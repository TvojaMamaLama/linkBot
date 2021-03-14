package main

import (
	"github.com/boltdb/bolt"
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

	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		log.Fatal()
	}

	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost/")
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
