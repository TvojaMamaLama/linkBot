package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"pocketTeleBot/pkg/database"
)

func (b *Bot) makeAuth(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(greeting, authLink))

	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	redirectUrl := b.generateRedirectLink(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	if err != nil {
		return "", err
	}

	if err := b.tokenDB.Save(chatID, requestToken, database.RequestToken); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationUrl(requestToken, redirectUrl)
}

func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.serverUrl, chatID)
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenDB.Get(chatID, database.AccessToken)
}
