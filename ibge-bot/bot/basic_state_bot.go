package bot

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

const (
	STATE_BASE_URL = "http://www.ibge.gov.br/estadosat/"
)

var (
	STATES_NAME = map[string]string{
		"ac": "Acre",
		"al": "Alagoas",
		"am": "Amazonas",
		"ap": "Amapá",
		"ba": "Bahia",
		"ce": "Ceará",
		"es": "Espírito Santo",
		"df": "Distrito Federal",
		"go": "Goiás",
		"ma": "Maranhão",
		"mg": "Minas Gerais",
		"ms": "Mato Grosso do Sul",
		"mt": "Mato Grosso",
		"pa": "Pará",
		"pe": "Pernambuco",
		"pb": "Paraíba",
		"pi": "Piauí",
		"pr": "Paraná",
		"rj": "Rio de Janeiro",
		"rn": "Rio Grande do Norte",
		"ro": "Rondônia",
		"rr": "Roraima",
		"rs": "Rio Grande do Sul",
		"sc": "Santa Catarina",
		"se": "Sergipe",
		"sp": "São Paulo",
		"to": "Tocantins",
	}
)

type BasicStateBot struct{}

func (self BasicStateBot) Run(db database.MongoDB) {

	var wg sync.WaitGroup

	doc, err := goquery.NewDocument(STATE_BASE_URL)
	parser.CheckError(err)

	doc.Find("#menu a").Each(func(_ int, s *goquery.Selection) {
		wg.Add(1)
		partialUrl, _ := s.Attr("href")
		url := STATE_BASE_URL + partialUrl

		go func() {
			self.ParseState(db, url)
			wg.Done()
		}()
	})

	wg.Wait()
}

func (self BasicStateBot) ParseState(db database.MongoDB, stateUrl string) {
	doc, err := goquery.NewDocument(stateUrl)
	parser.CheckError(err)
	source := models.Source{
		Url:  stateUrl,
		Note: "ibge",
	}

	data := doc.Find("#sintese tr")

	pUrl, _ := url.Parse(stateUrl)

	id := pUrl.Query().Get("sigla")
	capital := parser.ToUtf8(data.Eq(0).Find(".total").Text())
	population2014 := data.Eq(1).Find(".total").Text()
	population2010 := data.Eq(2).Find(".total").Text()
	area := data.Eq(3).Find(".total").Text()
	populationDensity := data.Eq(4).Find(".total").Text()
	numberOfMunicipalities, _ := strconv.Atoi(data.Eq(5).Find(".total").Text())

	log.Println(id, capital, population2014, population2010, area, populationDensity, numberOfMunicipalities)

	if STATES_NAME[id] == "" {
		panic(id)
	}
	q := bson.M{"id": id}
	_, err = db.Upsert(q, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"name":                   STATES_NAME[id],
			"capitalid":              models.MakeUri(capital),
			"population":             toFloat(population2014),
			"area":                   toFloat(area),
			"populationdensity":      toFloat(populationDensity),
			"numberofmunicipalities": numberOfMunicipalities,
		},
		"$addToSet": bson.M{
			"sources": source,
		},
	}, models.State{})
	doc = nil
	parser.CheckError(err)
}

func toFloat(n string) float64 {
	n = strings.Replace(strings.Replace(n, ".", "", -1), ",", ".", -1)
	f, _ := strconv.ParseFloat(n, 64)
	return f
}
