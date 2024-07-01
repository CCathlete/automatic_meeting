//go:build windows
// +build windows

package main

import (
	mtgH "autoMeeting/meetingHandler"
	ymlH "autoMeeting/yamlHandler"

	"github.com/go-ole/go-ole"
)

func main() {
	// Initialising COM
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)

	configYml := ymlH.ParseYaml("config.yaml")
	meetingUrl := configYml.MeetingInfo.MeetingUrl
	checkInterval := configYml.MeetingInfo.CheckInterval

	mtgH.StartMeeting(meetingUrl)

	mtgH.MonitorMeeting(meetingUrl, checkInterval)

}
