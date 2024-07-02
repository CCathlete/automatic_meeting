//go:build windows
// +build windows

package meetinghandler

import (
	"fmt"
	"log"
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
	class1Name := "VfPpkd-RLmnJb"
	class2Name := "VfPpkd-vQzf8d"
	joinButton, err := webDriver.FindElement(selenium.ByClassName, class1Name)
	if err != nil {
		fmt.Printf("Couldn't find the button with the class name %s, trying another class name. Err: %v\n", class1Name, err)
		joinButton, err = webDriver.FindElement(selenium.ByClassName, class2Name)
		if err != nil {
			log.Fatalf("Couldn't find the button with the second class name %s. Err: %v\n", class2Name, err)
		}
	}

	if err := joinButton.Click(); err != nil {
		log.Fatalf("Error clicking the join button: %v", err)
	}

	// err = exec.Command("cmd", "/C", "start", meetUrl).Run()
	// if err != nil {
	// 	fmt.Printf("Failed to start the meeting with url: %s\n%v", meetUrl, err)
	// } else {
	// 	fmt.Printf("Google meet started at: %s", time.Now())
	// }
	fmt.Printf("Google meet started at: %s", time.Now())
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
