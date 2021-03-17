package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user not authorized")
	errUnableToSave = errors.New("unable save url")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "")
	switch err {
	case errInvalidURL:
		msg.Text = ""
	case errUnableToSave:
		msg.Text = ""
	case errUnauthorized:
		msg.Text = ""
	default:
		msg.Text = ""
	}
	b.bot.Send(msg)
}
