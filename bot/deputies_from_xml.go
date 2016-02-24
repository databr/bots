package bot

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromXML struct{}

func (p SaveDeputiesFromXML) Run(DB database.MongoDB) {
	xmlURL := "http://www.camara.gov.br/SitCamaraWS/Deputados.asmx/ObterDeputados"

	source := models.Source{
		Url:  xmlURL,
		Note: "Câmara API",
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(xmlURL); e != nil {
		parser.Log.Critical(e.Error())
	}

	doc.Find("deputado").Each(func(i int, s *goquery.Selection) {
		name := parser.Titlelize(s.Find("nomeparlamentar").First().Text())
		parser.Log.Info("Saving " + name)
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
		fullName := strings.Split(parser.Titlelize(s.Find("nome").First().Text()), " ")

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
					Name:       parser.Titlelize(s.Find("nome").First().Text()),
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

		parser.CreateMembermeship(DB, models.Rel{
			Id:   parliamenrianId,
			Link: parser.LinkTo("parliamenrians", parliamenrianId),
		}, models.Rel{
			Id:   partyId,
			Link: parser.LinkTo("parties", partyId),
		}, source, "Filiado", "Partido")
		parser.CheckError(err)
	})
}
