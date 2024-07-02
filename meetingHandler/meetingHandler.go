//go:build windows
// +build windows

package meetinghandler

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/tebeka/selenium"
)

type process struct {
	Name        string
	CommandLine string
}

func StartMeeting(meetUrl string, chromeDriverPath string, port int) {
	options := []selenium.ServiceOption{}
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, options...)
	if err != nil {
		log.Fatalf("Couldn't start the ChromeDriver service: %v", err)
	}
	defer service.Stop()

	capabilities := selenium.Capabilities{
		"browserName": "chrome",
	}
	seleniumServerUrl := fmt.Sprintf("http://localhost:%d/wd/hub", port)
	// Starting a new session.
	webDriver, err := selenium.NewRemote(capabilities, seleniumServerUrl)
	if err != nil {
		log.Fatalf("Error connecting to the webDriver: %v", err)
	}
	defer webDriver.Quit()

	if err := webDriver.Get(meetUrl); err != nil {
		log.Fatalf("Error nevigating to the meeting's url: %v", err)
	}
	// Waiting for the page to load.
	time.Sleep(5 * time.Second)

	// Finding the join meeting button and clicking it.
	joinButton, err := webDriver.FindElement(selenium.ByCSSSelector, "")

	err = exec.Command("cmd", "/C", "start", meetUrl).Run()
	if err != nil {
		fmt.Printf("Failed to start the meeting with url: %s\n%v", meetUrl, err)
	} else {
		fmt.Printf("Google meet started at: %s", time.Now())
	}
}

func IsMeetRunning(meetUrl string) bool {
	var processes []process
	query := wmi.CreateQuery(&processes, "WHERE Name = 'chrome.exe'")
	err := wmi.Query(query, &processes)
	if err != nil {
		fmt.Printf("Failed to query processes: %v\n", err)
		return false
	}

	for _, process := range processes {
		if strings.Contains(process.CommandLine, meetUrl) {
			return true
		}
	}

	return false // No running processes with the meeting url.
}

func MonitorMeeting(meetUrl string, chromeDriverPath string, port int, checkInterval int) {
	time.Sleep(time.Duration(checkInterval) * time.Second)
	if !IsMeetRunning(meetUrl) {
		fmt.Println("Google meet is not running. Restarting at: ", time.Now())
		StartMeeting(meetUrl, chromeDriverPath, port)
	} else {
		fmt.Println("Google meet is running at ", time.Now())
	}
}
