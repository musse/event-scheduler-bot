package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (bot *EventBot) handleSchedule(update tgbotapi.Update) {
	user := update.Message.From
	if !bot.isAdmin(user) {
		logrus.Warningf("Can't schedule event: user is not admin")
		bot.sendReplyToUpdate(bot.NotAdminMessage, update)
		return
	}

	eventDescription := update.Message.CommandArguments()
	if eventDescription == "" {
		logrus.Warningf("Can't schedule event: no description provided")
		bot.sendReplyToUpdate(bot.NoEventDescriptionMessage, update)
		return
	}

	bot.currentEvent = &Event{
		Description:     eventDescription,
		Participants:    make([]Participant, 0),
		MaxParticipants: bot.MaxParticipants,
	}

	replyText := fmt.Sprintf(bot.ScheduledMessage + "\n" + eventDescription)
	bot.sendReplyToUpdate(replyText, update)

	logrus.Printf("New event scheduled: %s", eventDescription)
}

func (bot *EventBot) handleEdit(update tgbotapi.Update) {
	user := update.Message.From
	if !bot.isAdmin(user) {
		logrus.Warningf("Can't edit event: user is not admin")
		bot.sendReplyToUpdate(bot.NotAdminMessage, update)
		return
	}

	eventDescription := update.Message.CommandArguments()
	if eventDescription == "" {
		logrus.Warningf("Can't edit event: no description provided")
		bot.sendReplyToUpdate(bot.NoEventDescriptionMessage, update)
		return
	}

	bot.currentEvent.Description = eventDescription

	replyText := fmt.Sprintf(bot.EditedMessage + "\n" + eventDescription)
	bot.sendReplyToUpdate(replyText, update)

	logrus.Printf("Event description edited: %s", eventDescription)
}

func (bot *EventBot) handleGo(update tgbotapi.Update) {
	if bot.currentEvent == nil {
		logrus.Warningf("Can't add participant: no current event")
		bot.sendReplyToUpdate(bot.NoEventScheduledMessage, update)
		return
	}

	user := update.Message.From
	externalParticipant := update.Message.CommandArguments()

	var err error
	var participant Participant
	if externalParticipant == "" {
		participant, err = bot.currentEvent.addInGroupParticipant(*user)
	} else {
		participant, err = bot.currentEvent.addExternalParticipant(externalParticipant)
	}

	if err != nil {
		logrus.Warningf("Can't add participant %s: %s", participant, err)
		bot.sendReplyToUpdate(fmt.Sprintf("%s %s", participant, bot.AlreadyInEventMessage), update)
		return
	}

	bot.sendReplyToUpdate(bot.currentEvent.getEventMessage(bot.WaitlistedHeader), update)
	logrus.Printf("Participant %s added", participant)
}

func (bot *EventBot) handleQuit(update tgbotapi.Update) {
	if bot.currentEvent == nil {
		logrus.Warningf("Can't remove participant: no current event")
		bot.sendReplyToUpdate(bot.NoEventScheduledMessage, update)
		return
	}

	user := update.Message.From
	externalParticipant := update.Message.CommandArguments()

	var err error
	var participant Participant
	if externalParticipant == "" {
		participant, err = bot.currentEvent.removeInGroupParticipant(*user)
	} else {
		participant, err = bot.currentEvent.removeExternalParticipant(externalParticipant)
	}

	if err != nil {
		logrus.Warningf("Can't remove participant %s: %s", participant, err)
		bot.sendReplyToUpdate(fmt.Sprintf("%s %s", participant, bot.NotInEventMessage), update)
		return
	}

	bot.sendReplyToUpdate(bot.currentEvent.getEventMessage(bot.WaitlistedHeader), update)
	logrus.Printf("Participant %s removed", participant)
}

func (bot *EventBot) handleCancel(update tgbotapi.Update) {
	user := update.Message.From
	if !bot.isAdmin(user) {
		bot.sendReplyToUpdate(bot.NotAdminMessage, update)
		logrus.Warningf("Can't cancel event: user is not admin")
		return
	}

	if bot.currentEvent == nil {
		bot.sendReplyToUpdate(bot.NoEventScheduledMessage, update)
		logrus.Warningf("Can't cancel event: no current event")
		return
	}

	replyText := fmt.Sprintf(bot.CancelledMessage + "\n" + bot.currentEvent.Description)
	bot.sendReplyToUpdate(replyText, update)
	logrus.Printf("Event cancelled: %s", bot.currentEvent.Description)

	bot.currentEvent = nil
	bot.removeEventFile()
}
