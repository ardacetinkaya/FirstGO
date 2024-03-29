package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ConnectionString      string `json:"SQLConnection"`
	Port                  string `json:"Port"`
	Token                 string `json:"Token"`
	AzureQueueAccountName string `json:"AZQAccountName"`
	AzureQueueAccountKey  string `json:"AZQAccountKey"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
