package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/thestuckster/blink/cmd"
	"github.com/thestuckster/blink/internal"
	"log"
	"os"
)

func main() {
	initBlink()
	cmd.Execute()
}

func initBlink() {
	createDefaultConfigIfMissing()
	config := internal.LoadConfig()
	if !config.HasGamePath() {
		path, err := internal.FindGameInstallationDirectory()
		if err != nil {
			log.Panic(err)
		}
		config.GamePath = path

		err = config.Save()
		if err != nil {
			log.Panic(err)
		}
	}
}

func createDefaultConfigIfMissing() {
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		config := internal.Config{}
		config.Save()
	}
}
