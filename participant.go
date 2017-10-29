package main

import "github.com/go-telegram-bot-api/telegram-bot-api"

type Participant struct {
	// only set one of the below
	User *tgbotapi.User `yaml:"user,omitempty"`
	Name string         `yaml:"name,omitempty"`
}

func (p Participant) String() string {
	if p.User != nil {
		if p.User.UserName != "" {
			return "@" + p.User.UserName
		}
		return p.User.String()
	}
	return p.Name
}
