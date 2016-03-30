package bot

import (
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

const (
	CITIES_BASE_URL = "http://cidades.ibge.gov.br/download/mapa_e_municipios.php?lang=&uf="
)

type BasicCityBot struct{}

func (self BasicCityBot) Run(db database.MongoDB) {
	for uf, _ := range STATES_NAME {
		self.getCitiesData(db, CITIES_BASE_URL+uf, uf)
	}
}

func (self BasicCityBot) getCitiesData(db database.MongoDB, url string, stateID string) {
	doc, err := goquery.NewDocument(url)
	parser.CheckError(err)
	source := models.Source{
		Url:  url,
		Note: "ibge",
	}

	doc.Find("#municipios tbody tr").Each(func(_ int, s *goquery.Selection) {
		data := s.Find("td")

		name := data.Eq(0).Text()
		parser.Log.Debug("Salving: " + name + " (" + stateID + ")")
		id := models.MakeUri(name)

		ibgecode, _ := strconv.Atoi(data.Eq(1).Text())

		q := bson.M{"id": id, "stateid": stateID}
		_, err = db.Upsert(q, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"name":       name,
				"ibgecode":   ibgecode,
				"gentile":    data.Eq(2).Text(),
				"population": toFloat(data.Eq(3).Text()),
				"area":       toFloat(data.Eq(4).Text()),
				"density":    toFloat(data.Eq(5).Text()),
				"pib":        toFloat(data.Eq(6).Text()),
			},
			"$addToSet": bson.M{
				"sources": source,
			},
		}, models.City{})
		parser.CheckError(err)
	})

	doc = nil
}
