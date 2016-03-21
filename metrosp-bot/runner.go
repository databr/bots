package main

import (
  "github.com/databr/api/database"
  "github.com/databr/metrosp-bot/bot"
)

func main() {
  mongo := database.NewMongoDB()

  bot.StationBot{}.Run(mongo)
  bot.LineBot{}.Run(mongo)
  bot.StatusBot{}.Run(mongo)
}
