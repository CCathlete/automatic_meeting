//go:build windows
// +build windows

package meetinghandler

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/StackExchange/wmi"
)

type process struct {
	Name        string
	CommandLine string
}

func StartMeeting(meetUrl string) {
	err := exec.Command("cmd", "/C", "start", meetUrl).Run()
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
