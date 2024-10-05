package config

import (
	"encoding/json"
	"os"
)

var Config = TtmConfig{
	TasksPath: home + "/ttm",
}

var home = os.Getenv("HOME")
var configPath = home + "/.config/ttm"

type TtmConfig struct {
	TasksPath string `json:"tasksPath"`
}

func Init() error {
	file, err := os.Open(configPath + "/ttm.json")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(configPath, os.ModePerm)
			if err != nil {
				return err
			}
			file, err = os.Create(configPath + "/ttm.json")
			if err != nil {
				return err
			}
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(&Config)
			if err != nil {
				return err
			}
		}
	}

	decoder := json.NewDecoder(file)
	decoder.Decode(&Config)

	return nil
}
