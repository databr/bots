package parser

import (
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/models"
	"github.com/databr/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromSearch struct {
}

func (p SaveDeputiesFromSearch) Run(DB models.Database) {
	searchURL := "http://www2.camara.leg.br/deputados/pesquisa"

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(searchURL); e != nil {
		log.Critical(e.Error())
	}

	source := popolo.Source{
		Url:  toPtr(searchURL),
		Note: toPtr("Pesquisa Câmara"),
	}

	doc.Find("#deputado option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		if value != "" {
			info := regexp.MustCompile("=|%23|!|\\||\\?").Split(value, -1)

			name := titlelize(info[0])
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
					"identifiers": popolo.Identifier{
						Identifier: toPtr(info[2]), Scheme: toPtr("nMatricula"),
					},
				},
			}, models.Parliamentarian{})
			checkError(err)
		}
	})
}
