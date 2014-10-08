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

type SavePartiesFromTSE struct{}

func (_ SavePartiesFromTSE) Run(DB database.MongoDB) {
	url := "http://www.tse.jus.br/partidos/partidos-politicos/registrados-no-tse"
	source := models.Source{
		Url:  url,
		Note: "Tribunal Superior Eleitoral",
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(url); e != nil {
		parser.Log.Critical(e.Error())
	}

	const (
		IDX = iota
		SIGLA_IDX
		NAME_IDX
		DEFERIMENTO_IDX
		PRESIDENT_IDX
		N_IDX
	)

	doc.Find("#textoConteudo table tr").Each(func(i int, s *goquery.Selection) {
		if s.Find(".titulo_tabela").Length() < 6 && s.Find("td").Length() > 1 {
			info := s.Find("td")
			parser.Log.Info("%s - %s - %s - %s - %s - %s",
				info.Eq(IDX).Text(),
				info.Eq(SIGLA_IDX).Text(),
				info.Eq(NAME_IDX).Text(),
				info.Eq(DEFERIMENTO_IDX).Text(),
				info.Eq(PRESIDENT_IDX).Text(),
				info.Eq(N_IDX).Text(),
			)

			partyId := models.MakeUri(info.Eq(SIGLA_IDX).Text())
			DB.Upsert(bson.M{"id": partyId}, bson.M{
				"$setOnInsert": bson.M{
					"createdat": time.Now(),
				},
				"$currentDate": bson.M{
					"updatedat": true,
				},
				"$set": bson.M{
					"id":   partyId,
					"name": parser.Titlelize(info.Eq(NAME_IDX).Text()),
					"othernames": []bson.M{{
						"name": info.Eq(SIGLA_IDX).Text(),
					}},
					"classification": "party",
				},
				"$addToSet": bson.M{
					"sources": []models.Source{source},
				},
			}, &models.Party{})

			urlDetails, b := info.Eq(SIGLA_IDX).Find("a").Attr("href")
			if b {
				docDetails, err := goquery.NewDocument(urlDetails)
				if err != nil {
					parser.Log.Critical(err.Error())
				}

				sourceDetails := models.Source{
					Url:  urlDetails,
					Note: "Tribunal Superior Eleitoral",
				}

				contactdetails := make([]bson.M, 0)

				details := docDetails.Find("#ancora-text-um p")

				address := strings.Split(details.Eq(3).Text(), ":")[1]
				contactdetails = append(contactdetails, bson.M{
					"label":   "Endereço",
					"type":    "address",
					"value":   address,
					"sources": []models.Source{sourceDetails},
				})

				zipcode := strings.Split(details.Eq(4).Text(), ":")[1]
				contactdetails = append(contactdetails, bson.M{
					"label":   "CEP",
					"type":    "zipcode",
					"value":   zipcode,
					"sources": []models.Source{sourceDetails},
				})

				phoneString := strings.Split(details.Eq(5).Text(), ":")[1]
				phone := strings.Split(phoneString, "/")[0]
				contactdetails = append(contactdetails, bson.M{
					"label":   "Telefone",
					"type":    "phone",
					"value":   phone,
					"sources": []models.Source{sourceDetails},
				})

				faxString := strings.Split(details.Eq(6).Text(), ":")[1]
				fax := strings.Split(faxString, "/")[0]
				contactdetails = append(contactdetails, bson.M{
					"label":   "Fax",
					"type":    "fax",
					"value":   fax,
					"sources": []models.Source{sourceDetails},
				})

				website, ok := details.Eq(7).Find("a").Attr("href")
				if ok {
					contactdetails = append(contactdetails, bson.M{
						"label":   "Site",
						"type":    "website",
						"value":   website,
						"sources": []models.Source{sourceDetails},
					})
				}

				details.Eq(8).Find("a").Each(func(i int, ss *goquery.Selection) {
					email, ok := ss.Attr("href")
					if !ok {
						return
					}
					contactdetails = append(contactdetails, bson.M{
						"label":   "Email",
						"type":    "email",
						"value":   email,
						"sources": []models.Source{sourceDetails},
					})
				})

				//Label:   "Endereço",
				//Type:    "address",
				//Value:   street,
				//Sources: []models.Source{sourceDetails},
				//}, {
				//Label:   "CEP",
				//Type:    "zipcode",
				//Value:   zipcode,
				//Sources: []models.Source{sourceDetails},
				//}
				//

				data := bson.M{
					"$setOnInsert": bson.M{
						"createdat": time.Now(),
					},
					"$currentDate": bson.M{
						"updatedat": true,
					},
					"$set": bson.M{
						"contactdetails": contactdetails,
					},
				}

				DB.Upsert(bson.M{"id": partyId}, data, models.Party{})
			}
		}
	})
}
