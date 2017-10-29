package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (bot *EventBot) isAdmin(user *tgbotapi.User) bool {
	chatConfig := tgbotapi.ChatConfig{
		ChatID: bot.ChatID,
	}

	admins, err := bot.botAPI.GetChatAdministrators(chatConfig)
	if err != nil {
		panic(err)
	}

	for _, admin := range admins {
		if admin.User.ID == user.ID {
			return true
		}
	}
	return false
}

func (bot *EventBot) sendReplyToUpdate(replyText string, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := bot.botAPI.Send(msg)
	if err != nil {
		logrus.Errorf("Could not send message %s: %s", msg.Text, err)
	}
}
