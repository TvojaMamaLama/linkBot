package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"pocketTeleBot/pkg/pocketAPI"
)

const (
	commandStart = "start"
	greeting     = "Привет перейди по ссылке, чтобы дать доступ к сохранению ссылок!\n%s"
	alreadyAuth  = "Вы уже авторизированы"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Успешно сохранена!")

	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		msg.Text = "Некорректная ссылка."
		_, err := b.bot.Send(msg)
		return err
	}

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		msg.Text = "Ты не авторизирован, отправь /start."
		_, err := b.bot.Send(msg)
		return err
	}

	if err := b.pocketClient.Add(context.Background(), pocketAPI.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		msg.Text = "Не удалось сохранить ссылку.Попробуйте позже."
		_, err := b.bot.Send(msg)
		return err
	}

	_, err = b.bot.Send(msg)
	return err
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
	_, err := b.getAccessToken(command.Chat.ID)
	if err != nil {
		return b.makeAuth(command)
	}
	msg := tgbotapi.NewMessage(command.Chat.ID, alreadyAuth)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknown(command *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(command.Chat.ID, "Неизвестная команда")
	_, err := b.bot.Send(msg)
	return err
}
