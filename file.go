package main

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func (bot *EventBot) saveEvent() {
	eventData, err := yaml.Marshal(bot.currentEvent)
	if err != nil {
		logrus.Errorf("Could not marshal current event: %s", err)
	}

	err = ioutil.WriteFile(bot.EventFilePath, eventData, 0644)
	if err != nil {
		logrus.Errorf("Could not save event file %s: %s", bot.EventFilePath, err)
	}

	logrus.Infof("Saved event file %s", bot.EventFilePath)
}

func (bot *EventBot) loadEventFromFile() {
	eventData, err := ioutil.ReadFile(bot.EventFilePath)
	if err != nil {
		logrus.Warningf("Could not read event file %s: %s", bot.EventFilePath, err)
		return
	}
	bot.currentEvent = &Event{}
	err = yaml.Unmarshal(eventData, bot.currentEvent)
	if err != nil {
		logrus.Errorf("Could not unmarshal saved event: %s", err)
		return
	}

	logrus.Infof("Read event file %s", bot.EventFilePath)
}

func (bot *EventBot) removeEventFile() {
	err := os.Remove(bot.EventFilePath)
	if err != nil {
		logrus.Errorf("Could not remove event file %s: %s", bot.EventFilePath, err)
	}
}
