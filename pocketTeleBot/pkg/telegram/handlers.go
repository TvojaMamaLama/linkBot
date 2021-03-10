package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	commandStart = "start"
	greeting     = "Привет перейди по ссылке, чтобы дать доступ к сохранению ссылок!\n%s"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Fuck!!!")
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func (b *Bot) handleCommand(command *tgbotapi.Message) error {
	switch command.Command() {
	case commandStart:
		return b.handleStart(command)
	default:
		return b.handleUnknown(command)
	}
}

func (b *Bot) handleStart(command *tgbotapi.Message) error {
	authLink, err := b.generateAuthLink(command.Chat.ID)
	msg := tgbotapi.NewMessage(command.Chat.ID, fmt.Sprintf(greeting, authLink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknown(command *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(command.Chat.ID, "Неизвестная команда")
	_, err := b.bot.Send(msg)
	return err
}
