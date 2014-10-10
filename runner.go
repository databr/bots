package main

import (
	"github.com/databr/api/database"
	"github.com/databr/bots/go_bot/ibge_bot/bot"
)

func main() {
	mongo := database.NewMongoDB()

	bot.BasicStateBot{}.Run(mongo)
}
