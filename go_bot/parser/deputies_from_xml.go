package parser

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/models"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromXML struct{}

func (p SaveDeputiesFromXML) Run(DB models.Database) {
	xmlURL := "http://www.camara.gov.br/SitCamaraWS/Deputados.asmx/ObterDeputados"

	source := models.Source{
		Url:  xmlURL,
		Note: "Câmara API",
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(xmlURL); e != nil {
		log.Critical(e.Error())
	}

	doc.Find("deputado").Each(func(i int, s *goquery.Selection) {
		name := titlelize(s.Find("nomeparlamentar").First().Text())
		log.Info("Saving " + name)
		partyId := models.MakeUri(s.Find("partido").First().Text())
		DB.Upsert(bson.M{"id": partyId}, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"id":             partyId,
				"classification": "party",
			},
		}, &models.Party{})

		parliamenrianId := models.MakeUri(name)
		q := bson.M{
			"id": parliamenrianId,
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
				"id":       models.MakeUri(name),
				"gender":   s.Find("sexo").First().Text(),
				"image":    s.Find("urlFoto").First().Text(),
				"email":    s.Find("email").First().Text(),
			},
			"$addToSet": bson.M{
				"sources": source,
				"identifiers": bson.M{
					"$each": []models.Identifier{
						{Identifier: s.Find("idParlamentar").First().Text(), Scheme: "idParlamentar"},
						{Identifier: s.Find("ideCadastro").First().Text(), Scheme: "ideCadastro"},
					},
				},
				"othernames": models.OtherNames{
					Name:       titlelize(s.Find("nome").First().Text()),
					FamilyName: fullName[len(fullName)-1:][0],
					GivenName:  fullName[0],
					Note:       "Nome de nascimento",
				},
				"contactdetails": bson.M{
					"$each": []models.ContactDetail{
						{
							Label:   "Telefone",
							Type:    "phone",
							Value:   s.Find("fone").First().Text(),
							Sources: []models.Source{source},
						},
						{
							Label:   "Gabinete",
							Type:    "address",
							Value:   s.Find("gabinete").First().Text() + ", Anexo " + s.Find("anexo").First().Text(),
							Sources: []models.Source{source},
						},
					},
				},
			},
		}, &models.Parliamentarian{})

		createMembermeship(DB, models.Rel{
			Id:   parliamenrianId,
			Link: LinkTo("parliamenrians", parliamenrianId),
		}, models.Rel{
			Id:   partyId,
			Link: LinkTo("parties", partyId),
		}, source)
		checkError(err)
	})
}
