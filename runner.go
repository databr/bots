package main

import (
	"github.com/databr/api/database"
	"github.com/databr/bots/go_bot/sabesp_bot/bot"
)

func main() {
	mongo := database.NewMongoDB()

	bot.SabespBot{}.Run(mongo)
}
