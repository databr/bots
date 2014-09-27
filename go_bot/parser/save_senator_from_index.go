package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/models"
	"github.com/databr/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveSenatorsFromIndex struct {
}

func (self SaveSenatorsFromIndex) Run(DB models.Database) {
	indexURL := "http://www.senado.gov.br"

	source := models.Source{
		Url:  indexURL,
		Note: "senado.gov.br website",
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(indexURL + "/senadores/"); e != nil {
		log.Critical(e.Error())
	}

	doc.Find("#senadores tbody tr").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		name := data.Eq(0).Text()
		link, okLink := data.Eq(0).Find("a").Attr("href")
		if !okLink {
			checkError(errors.New("link not found"))
		} else {
			link = indexURL + link
		}

		email, okEmail := data.Eq(6).Find("a").Attr("href")
		if !okEmail {
			email = ""
		} else {
			email = strings.Replace(email, "mailto:", "", -1)
		}

		partyId := models.MakeUri(data.Eq(1).Text())
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

		re := regexp.MustCompile("paginst/senador(.+)a.asp")
		senatorId := re.FindStringSubmatch(link)[1]

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
					"$each": []models.ContactDetail{
						{
							Label:   "Telefone",
							Type:    "phone",
							Value:   data.Eq(4).Text(),
							Sources: []models.Source{source},
						},
						{
							Label:   "Fax",
							Type:    "fax",
							Value:   data.Eq(5).Text(),
							Sources: []models.Source{source},
						},
					},
				},
				"identifiers": bson.M{
					"$each": []models.Identifier{
						{Identifier: senatorId, Scheme: "CodSenador"},
					},
				},
			},
			"$set": bson.M{
				"name":      name,
				"email":     email,
				"link":      link,
				"shortname": models.MakeUri(name),
			},
		}, models.Parliamentarian{})
		checkError(err)

		docDetails, e := goquery.NewDocument(link)
		if e != nil {
			log.Critical(e.Error())
		}
		info := docDetails.Find(".dadosSenador b")
		birthdateA := strings.Split(info.Eq(1).Text(), "/")
		year, _ := strconv.Atoi(birthdateA[2])
		month, _ := strconv.Atoi(birthdateA[1])
		day, _ := strconv.Atoi(birthdateA[0])
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		birthDate := popolo.Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)}

		_, err = DB.Upsert(q, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"birthdate": birthDate,
			},
			"$addToSet": bson.M{
				"sources": source,
				"othernames": models.OtherNames{
					Name: info.Eq(0).Text(),
					Note: "Nome de nascimento",
				},
				"contactdetails": models.ContactDetail{
					Label:   "Gabinete",
					Type:    "address",
					Value:   info.Eq(4).Text(),
					Sources: []models.Source{source},
				},
			},
		}, models.Parliamentarian{})

		createMembermeship(DB, models.Rel{
			Id:   parliamenrianId,
			Link: LinkTo("parliamenrians", parliamenrianId),
		}, models.Rel{
			Id:   partyId,
			Link: LinkTo("parties", partyId),
		}, source, "Filiado", "Partido")

		checkError(err)
	})
}
