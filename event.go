package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Event struct {
	Description     string        `yaml:"description"`
	MaxParticipants int           `yaml:"maxParticipants"`
	Participants    []Participant `yaml:"participants"`
}

func (e *Event) addInGroupParticipant(user tgbotapi.User) (Participant, error) {
	newParticipant := Participant{
		User: &user,
	}

	return newParticipant, e.addParticipant(newParticipant)
}

func (e *Event) addExternalParticipant(name string) (Participant, error) {
	newParticipant := Participant{
		Name: name,
	}

	return newParticipant, e.addParticipant(newParticipant)
}

func (e *Event) getEventMessage(waitlistedHeader string) string {
	var msg string
	msg += e.Description + "\n"

	var participantsCount int
	if len(e.Participants) <= e.MaxParticipants {
		participantsCount = len(e.Participants)
	} else {
		participantsCount = e.MaxParticipants
	}

	for i, participant := range e.Participants[:participantsCount] {
		msg += fmt.Sprintf("%d. %s\n", i+1, participant)
	}

	if len(e.Participants) > e.MaxParticipants {
		if e.MaxParticipants != 0 {
			msg += waitlistedHeader + "\n"
		}
		for i, participant := range e.Participants[e.MaxParticipants:] {
			msg += fmt.Sprintf("%d. %s\n", i+1, participant)
		}
	}

	return msg
}

func (e *Event) addParticipant(participant Participant) error {
	isInEvent, _ := e.isInEvent(participant)
	if isInEvent {
		return fmt.Errorf("%s is already in the event", participant)
	}

	e.Participants = append(e.Participants, participant)

	return nil
}

func (e *Event) removeInGroupParticipant(user tgbotapi.User) (Participant, error) {
	newParticipant := Participant{
		User: &user,
	}

	return newParticipant, e.removeParticipant(newParticipant)
}

func (e *Event) removeExternalParticipant(name string) (Participant, error) {
	newParticipant := Participant{
		Name: name,
	}

	return newParticipant, e.removeParticipant(newParticipant)
}

func (e *Event) removeParticipant(participant Participant) error {
	isInEvent, index := e.isInEvent(participant)
	if !isInEvent {
		return fmt.Errorf("%s is not in the event", participant)
	}

	e.Participants = append(e.Participants[:index], e.Participants[index+1:]...)

	return nil
}

func (e *Event) isInEvent(participant Participant) (bool, int) {
	participantString := participant.String()
	for i, p := range e.Participants {
		pString := p.String()
		if participantString == pString {
			return true, i
		}
	}
	return false, -1
}
