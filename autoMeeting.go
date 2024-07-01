//go:build windows
// +build windows

package main

import (
	ymlH "autoMeeting/yamlHandler"
)

func main() {
	configYml := ymlH.ParseYaml("config.yaml")
	meetingUrl := configYml.MeetingInfo.MeetingUrl
}
