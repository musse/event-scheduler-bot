package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

const defaultConfigFile = "config.yaml"

func main() {
	var configFilePath string
	if len(os.Args) > 1 {
		configFilePath = os.Args[1]
	} else {
		configFilePath = defaultConfigFile
	}

	bot, err := newEventBot(configFilePath)
	if err != nil {
		logrus.Fatalf("Could not start Event Bot: %s", err)
	}

	bot.loadEventFromFile()

	logrus.Fatal(bot.handleUpdates())
}
