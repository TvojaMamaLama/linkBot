package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pocketTeleBot/pkg/database"
	"pocketTeleBot/pkg/pocketAPI"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	pocketClient *pocketAPI.Client
	tokenDB      database.TokenDB
	serverUrl    string
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocketAPI.Client, serverUrl string, tokenDB database.TokenDB) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, serverUrl: serverUrl, tokenDB: tokenDB}
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
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
