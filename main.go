package main

import (
	tgClient "example/tlgbot/clients/telegram"
	"example/tlgbot/consumer/event_consumer"
	"example/tlgbot/events/telegram"
	"example/tlgbot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	storagePath = "data"
	batchSize = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Println("server started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err:= consumer.Start(); err != nil {
		log.Fatal()
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access tot telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not valid")
	}

	return *token
}
