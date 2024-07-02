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
	"github.com/tebeka/selenium/chrome"
)

type process struct {
	Name        string
	CommandLine string
}

func StartMeeting(meetUrl, chromeProfileDir, chromeDriverPath string, port int, htmlIdentifiers []string) {
	options := []selenium.ServiceOption{}
	selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, options...)
	if err != nil {
		log.Fatalf("Couldn't start the ChromeDriver service: %v", err)
	}
	defer service.Stop()

	chromeOptions := chrome.Capabilities{
		Args: []string{
			"--disable-gpu",
			"--disable-extensions",
			"--disable-dev-shm-usage",
			"--no-sandbox",
			"--disable-infobars",
			"--start-maximized",
			"--disable-browser-side-navigation",
			"--disable-blink-features=AutomationControlled",
			fmt.Sprintf("--user-data-dir=%s", chromeProfileDir),
		},
	}

	capabilities := selenium.Capabilities{
		"browserName":   "chrome",
		"chromeOptions": chromeOptions,
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
	class1Name := htmlIdentifiers[0]
	class2Name := htmlIdentifiers[1]
	xpath := htmlIdentifiers[2]
	joinButton, err := webDriver.FindElement(selenium.ByClassName, class1Name)
	var err2 error
	if err == nil {
		// class 1 was found so we're trying to click it.
		err2 = joinButton.Click()
		if err2 != nil {
			fmt.Printf("Error clicking the join button with class 1 %s, trying class 2: %v", class1Name, err)
		}
	}
	if err != nil || err2 != nil {
		// class 1 wasn't found or wasn't clickable so we're trying class 2.
		if err != nil {
			fmt.Printf("\nCouldn't find the button with the class name %s, trying another class name. \nErr: %v\n", class1Name, err)
		}
		joinButton, err = webDriver.FindElement(selenium.ByClassName, class2Name)
		if err == nil {
			// class 2 was found so we're trying to click it.
			err2 = joinButton.Click()
			if err2 != nil {
				fmt.Printf("Error clicking the join button with class 2 %s, trying xpath: %v", class2Name, err)
			}
		}
		if err != nil || err2 != nil {
			// class 2 wasn't found or wasn't clickable so we're trying xpath.
			if err != nil {
				fmt.Printf("\nCouldn't find the button with the second class name %s, trying using xpath. \nErr: %v\n", class2Name, err)
			}
			joinButton, err = webDriver.FindElement(selenium.ByXPATH, xpath)
			if err == nil {
				// xpath was found so we're trying to click it.
				err2 = joinButton.Click()
				if err2 != nil {
					log.Fatalf("Error clicking the join button using xpath %s: %v", xpath, err)
				}
			}
		}
		if err != nil {
			// xpath wasn't found.
			log.Fatalf("\nCouldn't find the button with the xpath %s. \nErr: %v\n", xpath, err)
		}
	}

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

func MonitorMeeting(meetUrl, chromeProfileDir, chromeDriverPath string,
	port int, checkInterval int, htmlIdentifiers []string) {
	time.Sleep(time.Duration(checkInterval) * time.Second)
	if !IsMeetRunning(meetUrl) {
		fmt.Println("Google meet is not running. Restarting at: ", time.Now())
		StartMeeting(meetUrl, chromeProfileDir, chromeDriverPath, port, htmlIdentifiers)
	} else {
		fmt.Println("Google meet is running at ", time.Now())
	}
}
