package parser

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/models"
	"github.com/databr/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromXML struct{}

func (p SaveDeputiesFromXML) Run(DB models.Database) {
	xmlURL := "http://www.camara.gov.br/SitCamaraWS/Deputados.asmx/ObterDeputados"

	source := popolo.Source{
		Url:  toPtr(xmlURL),
		Note: toPtr("Câmara API"),
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(xmlURL); e != nil {
		log.Critical(e.Error())
	}

	doc.Find("deputado").Each(func(i int, s *goquery.Selection) {
		partyId := toPtr(models.MakeUri(s.Find("partido").First().Text()))
		DB.Upsert(bson.M{"id": partyId}, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"id":             partyId,
				"classification": toPtr("party"),
			},
		}, &models.Party{})

		//PartyId:    party.Id,
		//State:      s.Find("uf").First().Text(),

		name := titlelize(s.Find("nomeparlamentar").First().Text())
		q := bson.M{
			"id": models.MakeUri(name),
		}
		fullName := strings.Split(titlelize(s.Find("nome").First().Text()), " ")

		_, err := DB.Upsert(q, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"name":     &name,
				"sortname": &name,
				"id":       toPtr(models.MakeUri(name)),
				"gender":   toPtr(s.Find("sexo").First().Text()),
				"image":    toPtr(s.Find("urlFoto").First().Text()),
				"email":    toPtr(s.Find("email").First().Text()),
			},
			"$addToSet": bson.M{
				"sources": source,
				"identifiers": bson.M{
					"$each": []popolo.Identifier{
						{Identifier: toPtr(s.Find("idParlamentar").First().Text()), Scheme: toPtr("idParlamentar")},
						{Identifier: toPtr(s.Find("ideCadastro").First().Text()), Scheme: toPtr("ideCadastro")},
					},
				},
				"othernames": popolo.OtherNames{
					Name:       toPtr(titlelize(s.Find("nome").First().Text())),
					FamilyName: toPtr(fullName[len(fullName)-1:][0]),
					GivenName:  &fullName[0],
					Note:       toPtr("Nome de nascimento"),
				},
				"contactdetails": bson.M{
					"$each": []popolo.ContactDetail{
						{
							Label:   toPtr("Telefone"),
							Type:    toPtr("phone"),
							Value:   toPtr(s.Find("fone").First().Text()),
							Sources: []popolo.Source{source},
						},
						{
							Label:   toPtr("Gabinete"),
							Type:    toPtr("address"),
							Value:   toPtr(s.Find("gabinete").First().Text() + ", Anexo " + s.Find("anexo").First().Text()),
							Sources: []popolo.Source{source},
						},
					},
				},
			},
		}, &models.Parliamentarian{})
		checkError(err)
	})
}

func titlelize(s string) string {
	return strings.Title(strings.ToLower(s))
}

func toPtr(s string) *string {
	return &s
}
