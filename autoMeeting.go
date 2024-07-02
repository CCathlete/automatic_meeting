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
	chromeDriverPath := configYml.ChromeConfig.ChromeDriverPath
	chromeProfileDir := configYml.ChromeConfig.ChromeProfileDir
	seleniumPort := configYml.ChromeConfig.Port
	endMeeting := configYml.MeetingInfo.EndMeeting
	htmlIdentifiers := []string{
		configYml.ChromeConfig.ClassId1,
		configYml.ChromeConfig.ClassId2,
		configYml.ChromeConfig.Xpath,
	}

	mtgH.StartMeeting(meetingUrl, chromeProfileDir, chromeDriverPath, seleniumPort, htmlIdentifiers, endMeeting)

	for !endMeeting {
		endMeeting = ymlH.ParseYaml("config.yaml").MeetingInfo.EndMeeting
		mtgH.MonitorMeeting(meetingUrl, chromeProfileDir, chromeDriverPath, seleniumPort, checkInterval, htmlIdentifiers, endMeeting)
	}

}
