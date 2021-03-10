package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pocketTeleBot/pkg/pocketAPI"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocketAPI.Client
	redirectUrl  string
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocketAPI.Client, redirectUrl string) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectUrl: redirectUrl}
}

func (b *Bot) Start() error {
	log.Printf("Authorize on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		log.Panic(err)
	}

	b.updatesHandler(updates)

	return nil
}

func (b *Bot) updatesHandler(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
