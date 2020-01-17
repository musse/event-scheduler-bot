package main

import (
	"fmt"
	"io/ioutil"

	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const defaultEventFile = "event.yaml"

type EventBot struct {
	EventBotConfig

	botAPI       *tgbotapi.BotAPI
	currentEvent *Event
}

type EventBotConfig struct {
	APIToken string `yaml:"apiToken"`
	ChatID   int64  `yaml:"chatId"`

	EventFilePath string `yaml:"eventFilePath"`

	MaxParticipants int `yaml:"maxParticipants"`

	ScheduleCommand string `yaml:"scheduleCommand"`
	EditCommand     string `yaml:"editCommand"`
	GoCommand       string `yaml:"goCommand"`
	QuitCommand     string `yaml:"quitCommand"`
	CancelCommand   string `yaml:"cancelCommand"`

	WaitlistedHeader          string `yaml:"waitlistedHeader"`
	ScheduledMessage          string `yaml:"scheduledMessage"`
	CancelledMessage          string `yaml:"cancelledMessage"`
	EditedMessage             string `yaml:"editedMessage"`
	NotAdminMessage           string `yaml:"notAdminMessage"`
	NoEventScheduledMessage   string `yaml:"noEventScheduledMessage"`
	NoEventDescriptionMessage string `yaml:"noEventDescriptionMessage"`
	AlreadyInEventMessage     string `yaml:"alreadyInEventMessage"`
	NotInEventMessage         string `yaml:"notInEventMessage"`
}

func newEventBot(configFilePath string) (*EventBot, error) {
	logrus.Infof("Reading config file %s", configFilePath)

	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %s", err)
	}

	config := EventBotConfig{}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config data: %s", err)
	}

	bot := EventBot{
		EventBotConfig: config,
		currentEvent:   nil,
	}

	if bot.EventFilePath == "" {
		bot.EventFilePath = defaultEventFile
	}

	bot.botAPI, err = tgbotapi.NewBotAPI(bot.APIToken)
	if err != nil {
		return nil, fmt.Errorf("could not create bot API: %s", err)
	}

	return &bot, nil
}

func (bot *EventBot) handleUpdates() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.botAPI.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("error getting updates channel: %s", err)
	}

	logrus.Infoln("Handling updates...")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != bot.ChatID {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		logrus.Infof("Handling update: %s", update.Message.Text)

		cmd := update.Message.Command()
		switch cmd {
		case bot.ScheduleCommand:
			bot.handleSchedule(update)
		case bot.EditCommand:
			bot.handleEdit(update)
		case bot.GoCommand:
			bot.handleGo(update)
		case bot.QuitCommand:
			bot.handleQuit(update)
		case bot.CancelCommand:
			bot.handleCancel(update)
		default:
			logrus.Errorf("Unsupported command: %s", cmd)
			continue
		}

		if bot.currentEvent != nil {
			bot.saveEvent()
		}
	}

	return errors.New("end of updates channel")
}
