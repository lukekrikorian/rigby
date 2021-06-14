package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Configuration structure for the project
type Configuration struct {
	DatabaseURL string `json:"database"`
	Server      struct {
		Port   int32
		Origin string
	}
	HTTPS struct {
		Certificate string
		Key         string
	} `json:"https"`
}

// Config is the actual config data
var Config Configuration

// Init fetches the file data and decodes it, setting up the config
func init() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config.json file. Does it exist?")
		configFile.Close()
		return
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)
	if err != nil {
		log.Fatal("Error parsing json in config file.")
	}

	fmt.Println("Parsed config file")
}
