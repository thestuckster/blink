package internal

import (
	"encoding/json"
	"log"
	"os"
)

const AddOnsFolder = "\\_retail_\\Interface\\AddOns"

type Config struct {
	GamePath string
	AddOns   []AddOn
}

type AddOn struct {
	Url     string
	Repo    string
	Version string
}

func LoadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}

	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Panic(err)
	}

	return config
}

func (config *Config) HasGamePath() bool {
	if config.GamePath == "" {
		return false
	}

	return true
}

func (config *Config) AddAddOn(url, repo, version string) {
	addOn := AddOn{
		Url:     url,
		Repo:    repo,
		Version: version,
	}

	config.AddOns = append(config.AddOns, addOn)
}

func (config *Config) Save() error {
	err := os.Remove("config.json")
	if err != nil {
		return err
	}

	jsonString, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	if err = os.WriteFile("config.json", jsonString, 0644); err != nil {
		return err
	}

	return nil
}
