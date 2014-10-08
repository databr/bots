package main

import (
	"github.com/databr/api/database"
	"github.com/databr/bots/go_bot/parliamentarian_bot/bot"
)

func main() {
	mongo := database.NewMongoDB()

	bot.SaveDeputiesAbout{}.Run(mongo)
	bot.SaveDeputiesFromSearch{}.Run(mongo)
	bot.SaveDeputiesFromTransparenciaBrasil{}.Run(mongo)
	bot.SaveDeputiesFromXML{}.Run(mongo)
	//	bot.SaveDeputiesQuotas{}.Run(mongo)
	bot.SavePartiesFromTSE{}.Run(mongo)
	bot.SaveSenatorsFromIndex{}.Run(mongo)
}
