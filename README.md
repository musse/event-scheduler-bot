# Event Scheduler Bot

A Telegram bot to schedule events with your friends.

## How it works

Group admins may create an event with `/schedule` and cancel it with
`/cancel`. Regular users may add themselves to the event with `/go`
and remove themselves from it with `/quit`.

A max number of participants may be set. If so, there is an wait list.

The bot's messages and commands may be customized and/or translated to
your language.

For now, the bot supports one event at a time and works for a single
Telegram group.

## How to use

1. [Create a Telegram bot][1] and note down its API Token
2. Create a Telegram group with yours friends and your bot
3. [Get your group's chat ID][2]
4. Edit the `config.yaml` config file:
  - Set the `apiToken` and the `chatId` (mandatory)
  - Set the max number of participants
  - Change the messages and commands to your language
5. Add the commands you set in the config file to the [BotFather](3) with `/setcommands`
6. Install this bot with `go get github.com/musse/event-scheduler-bot`
7. Run the bot from the folder with your `config.yaml` file and leave it running

## Commands

- `/schedule {event-description}`: schedule a new event, replacing the current one if it
exists. Group admins only.
- `/go`: add yourself to the current event.
- `/go {participant-name}`: add another person to the event.
- `/quit`: remove yourself from the current event.
- `/quit {participant-name}`: remove another person from the event.
- `/cancel`: cancel the current event. Group admins only.

## Details

- Info about the current event in persisted in the `event.yaml` file. You may change
the file to which the event is persisted by setting eventFilePath in the config file.
- The default config file is `config.yaml` from the folder you run the bot. You may
change this by passing the your config file's path as an argument to the bot executable.

[1]: https://core.telegram.org/bots#3-how-do-i-create-a-bot
[2]: https://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id
[3]: https://core.telegram.org/bots#6-botfather