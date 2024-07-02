package yamlhandler

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Meeting struct {
	MeetingUrl    string `yaml:"url"`
	CheckInterval int    `yaml:"check interval"`
	EndMeeting    bool   `yaml:"end meeting?"`
}

type BrowserConfig struct {
	ChromeDriverPath string `yaml:"chrome driver path"`
	Port             int    `yaml:"port"`
}

type Config struct {
	MeetingInfo  Meeting       `yaml:"Meeting info"`
	ChromeConfig BrowserConfig `yaml:"Browser config"`
}

func ParseYaml(yamlPath string) Config {
	file, err := os.OpenFile(yamlPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Problem with opening the file in path %s: %v", yamlPath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("Problem with closing the file in path %s: %v", yamlPath, err)
		}
	}()
	fileInfo, _ := file.Stat()
	dataBytes := make([]byte, fileInfo.Size())

	// Reading the file we opened.
	_, err = file.Read(dataBytes)
	if err != nil {
		log.Fatalf("Problem with reading the file in path %s: %v", yamlPath, err)
	}
	// Similar to make(map[string]interface{} but also initialises it)
	configMap := Config{}
	err = yaml.Unmarshal(dataBytes, &configMap)
	if err != nil {
		log.Fatalf("Couldn't Unmarshal the yaml file in %s: %v", yamlPath, err)
	}

	return configMap
}
