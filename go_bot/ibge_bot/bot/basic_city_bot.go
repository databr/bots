package bot

import (
	"sync"
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
	var wg sync.WaitGroup
	for uf, _ := range STATES_NAME {
		wg.Add(1)
		go func(uf string) {
			self.getCitiesData(db, CITIES_BASE_URL+uf, uf)
			wg.Done()
		}(uf)
	}
	wg.Wait()
}

func (self BasicCityBot) getCitiesData(db database.MongoDB, url string, stateID string) {
	doc, err := goquery.NewDocument(url)
	parser.CheckError(err)

	doc.Find("#municipios tbody tr").Each(func(_ int, s *goquery.Selection) {
		data := s.Find("td")

		name := data.Eq(0).Text()
		parser.Log.Debug("Salving: " + name + " (" + stateID + ")")
		id := models.MakeUri(name)
		q := bson.M{"id": id, "state_id": stateID}
		_, err = db.Upsert(q, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"name":       name,
				"ibge_code":  data.Eq(1).Text(),
				"gentile":    data.Eq(2).Text(),
				"population": data.Eq(3).Text(),
				"area":       data.Eq(4).Text(),
				"density":    data.Eq(5).Text(),
				"pib":        data.Eq(6).Text(),
			},
		}, models.City{})
		parser.CheckError(err)
	})
}
