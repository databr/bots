package parser

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/camarabook/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveSenatorsFromIndex struct {
}

func (self SaveSenatorsFromIndex) Run(DB models.Database) {
	indexURL := "http://www.senado.gov.br/senadores"

	source := popolo.Source{
		Url:  toPtr(indexURL),
		Note: toPtr("senado.gov.br website"),
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(indexURL); e != nil {
		log.Critical(e.Error())
	}

	doc.Find("#senadores tbody tr").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		name := data.Eq(0).Text()
		link, okLink := data.Eq(0).Find("a").Attr("href")
		if !okLink {
			link = ""
		} else {
			link = indexURL + link
		}

		email, okEmail := data.Eq(6).Find("a").Attr("href")
		if !okEmail {
			email = ""
		} else {
			email = strings.Replace(email, "mailto:", "", -1)
		}

		partyId := toPtr(models.MakeUri(data.Eq(1).Text()))
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
				"contactdetails": bson.M{
					"$each": []popolo.ContactDetail{
						{
							Label:   toPtr("Telefone"),
							Type:    toPtr("phone"),
							Value:   toPtr(data.Eq(4).Text()),
							Sources: []popolo.Source{source},
						},
						{
							Label:   toPtr("Fax"),
							Type:    toPtr("fax"),
							Value:   toPtr(data.Eq(5).Text()),
							Sources: []popolo.Source{source},
						},
					},
				},
			},
			"$set": bson.M{
				"name":  name,
				"email": email,
				"link":  toPtr(link),
			},
		}, models.Parliamentarian{})
		checkError(err)
	})
}
