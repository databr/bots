package main

import (
	"github.com/databr/api/database"
	"github.com/databr/ibge-bot/bot"
)

func main() {
	mongo := database.NewMongoDB()

	bot.BasicStateBot{}.Run(mongo)
	bot.BasicCityBot{}.Run(mongo)
}
