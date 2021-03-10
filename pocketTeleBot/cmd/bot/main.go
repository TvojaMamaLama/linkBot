package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pocketTeleBot/pkg/pocketAPI"
	"pocketTeleBot/pkg/telegram"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocketAPI.NewCLient(consumerKey)
	if err != nil {
		log.Fatal(err)
	}

	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost/")
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
