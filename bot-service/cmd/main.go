package main

import (
	"log"

	"github.com/rf-krcn/telegram-removeBG/bot-service/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const token = "YOU_BOT_TOKEN"

func main() {
	// Initialize Telegram bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	handlers.HandleTelegramUpdates(bot)

}
