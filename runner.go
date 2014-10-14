package main

import (
	"github.com/databr/api/database"
	"github.com/databr/bots/go_bot/metrosp_bot/bot"
)

func main() {
	mongo := database.NewMongoDB()

	bot.StationBot{}.Run(mongo)
	bot.LineBot{}.Run(mongo)
	bot.StatusBot{}.Run(mongo)
}
