package parser

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/camarabook/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesAbout struct {
}

func (p SaveDeputiesAbout) Run(DB models.Database) {
	var ds []models.Parliamentarian

	DB.FindAll(&ds)

	for _, d := range ds {
		id := getIdDeputado(d.Identifiers)
		bioURL := "http://www2.camara.leg.br/deputados/pesquisa/layouts_deputados_biografia?pk=" + id
		source := popolo.Source{
			Url:  toPtr(bioURL),
			Note: toPtr("Pesquisa Câmara"),
		}

		var doc *goquery.Document
		var e error

		if doc, e = goquery.NewDocument(bioURL); e != nil {
			log.Fatal(e)
		}

		bio := doc.Find("#bioDeputado .bioOutros")

		biographyItems := make([]string, 0)
		bio.Each(func(i int, s *goquery.Selection) {
			title := s.Find(".bioOutrosTitulo").Text()
			if title != "" {
				title = strings.TrimSpace(title)
				title = strings.Replace(title, ":", "", -1)

				body := s.Find(".bioOutrosTexto").Text()

				biographyItems = append(biographyItems, title)
				biographyItems = append(biographyItems, body)
				biographyItems = append(biographyItems, "")
			}
		})

		_, err := DB.Upsert(bson.M{"id": d.Id}, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"summary":   bio.Eq(0).Find(".bioOutrosTexto").Text(),
				"biography": strings.Join(biographyItems, "\n"),
			},
			"$addToSet": bson.M{
				"sources": source,
			},
		}, models.Parliamentarian{})
		checkError(err)
	}
}

func getIdDeputado(ids []popolo.Identifier) string {
	for _, id := range ids {
		if *id.Scheme == "ideCadastro" {
			return *id.Identifier
		}
	}
	panic("not found id")
}
