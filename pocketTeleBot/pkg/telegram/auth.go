package telegram

import (
	"context"
	"fmt"
)

func (b *Bot) generateAuthLink(chatID int64) (string, error) {
	redirectUrl := b.generateRedirectLink(chatID)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	if err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationUrl(requestToken, redirectUrl)
}

func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatID)
}
