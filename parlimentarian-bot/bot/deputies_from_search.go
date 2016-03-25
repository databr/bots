package bot

import (
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromSearch struct {
}

func (p SaveDeputiesFromSearch) Run(DB database.MongoDB) {
	searchURL := "http://www2.camara.leg.br/deputados/pesquisa"

	if parser.IsCached(searchURL) {
		parser.Log.Info("SaveDeputiesFromSearch Cached")
		return
	}
	defer parser.DeferedCache(searchURL)

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(searchURL); e != nil {
		parser.Log.Critical(e.Error())
	}

	source := models.Source{
		Url:  searchURL,
		Note: "Pesquisa Câmara",
	}

	doc.Find("#deputado option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		if value != "" {
			info := regexp.MustCompile("=|%23|!|\\||\\?").Split(value, -1)

			name := parser.Titlelize(info[0])
			q := bson.M{
				"id": models.MakeUri(name),
			}

			_, err := DB.Upsert(q, bson.M{
				"$setOnInsert": bson.M{
					"createdat": time.Now(),
				},
				"$currentDate": bson.M{
					"updatedat": true,
				},
				"$addToSet": bson.M{
					"sources": source,
					"identifiers": models.Identifier{
						Identifier: info[2], Scheme: "nMatricula",
					},
				},
			}, models.Parliamentarian{})
			parser.CheckError(err)
		}
	})
}
