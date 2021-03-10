package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pocketTeleBot/pkg/telegram"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("1603671499:AAFnHg70CW0EghxfffsFU_yu3vKwuNgyyrE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	telegramBot := telegram.NewBot(bot)
	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
