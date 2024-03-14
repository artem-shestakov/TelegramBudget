package main

import (
	"os"

	"github.com/artem-shestakov/telegram-budget/internal/bot"
	"github.com/artem-shestakov/telegram-budget/internal/repository"
	"github.com/artem-shestakov/telegram-budget/internal/service"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("'BOT_TOKEN' environment variable is empty")
	}

	db, err := repository.NewPgDb(logger)
	if err != nil {
		panic(err)
	}
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	tgBot := bot.NewTgBot(token, service, logger)

	tgBot.Run()
}
