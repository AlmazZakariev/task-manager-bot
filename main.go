package main

import (
	"context"
	"flag"
	"log"
	tgClient "tasks-manager-bot/clients/telegram"
	event_consumer "tasks-manager-bot/consumer/event-consumer"
	"tasks-manager-bot/events/telegram"
	"tasks-manager-bot/storage/postgres"
)

const (
	tgBotHost                 = "api.telegram.org"
	batchSize                 = 100
	pqStorageConnectionString = "postgres://postgres:admin@localhost:5432/taskmanagerbottest"
	pqConStr                  = "user=postgres password=admin dbname=taskmanagerbottest sslmode=disable"
)

type Client interface {
	Create()
	MakeRemind()
	Show()
	Delete()
}

func main() {
	s, err := postgres.New(pqConStr)
	if err != nil {
		log.Fatal("can't connect to storage", err)
	}
	err = s.Init(context.TODO())
	if err != nil {
		log.Fatal("can't init storage", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("servise is stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
func token() string {
	return ""
}
