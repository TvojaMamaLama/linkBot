package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"pocketTeleBot/pkg/database"
	"pocketTeleBot/pkg/database/boltDB"
	"pocketTeleBot/pkg/pocketAPI"
	"pocketTeleBot/pkg/server"
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

	db, err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	tokenDB := boltDB.NewTokenDB(db)

	telegramBot := telegram.NewBot(bot, pocketClient, "http://localhost:8000/", tokenDB)
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Starting Server...")
	authServer := server.NewAuthorizationServer(pocketClient, tokenDB, "https://t.me/pocket_save_bot")
	if err := authServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDb() (*bolt.DB, error) {
	db, err := bolt.Open("bot.database", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(database.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(database.RequestToken))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, err
}
