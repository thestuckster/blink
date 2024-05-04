package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/thestuckster/blink/cmd"
	"github.com/thestuckster/blink/internal"
	"log"
)

func main() {
	initBlink()
	cmd.Execute()
}

func initBlink() {
	config := internal.Config{}
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
