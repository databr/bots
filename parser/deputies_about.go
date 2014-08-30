package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	. "github.com/camarabook/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesAbout struct {
}

func (p SaveDeputiesAbout) Run(DB models.Database) {
	var ds []models.Parliamentarian

	DB.FindAll(&ds)

	for _, d := range ds {
		id, ok := getIdDeputado(d.Identifiers)
		if !ok {
			continue
		}

		bioURL := "http://www2.camara.leg.br/deputados/pesquisa/layouts_deputados_biografia?pk=" + id
		source := Source{
			Url:  toPtr(bioURL),
			Note: toPtr("Pesquisa Câmara"),
		}

		var doc *goquery.Document
		var e error

		if doc, e = goquery.NewDocument(bioURL); e != nil {
			log.Critical(e.Error())
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

		bioDetails := doc.Find("#bioDeputado .bioDetalhes strong")

		birthdateA := strings.Split(bioDetails.Eq(1).Text(), "/")

		var year int
		switch id {
		case "123756", "160635":
			year = 1970
		case "74230", "129618":
			year = 1952
		case "74665", "141387":
			year = 1953
		case "73933":
			year = 1959
		case "73786":
			year = 1939
		case "74124":
			year = 1964
		default:
			log.Debug("(%s) %s", id, birthdateA)
			year, _ = strconv.Atoi(birthdateA[2])
		}

		month, _ := strconv.Atoi(birthdateA[1])
		day, _ := strconv.Atoi(birthdateA[0])
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		birthDate := Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)}

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
				"link":      "http://www.camara.gov.br/internet/Deputado/dep_Detalhe.asp?id=" + id,
				"birthdate": birthDate,
			},
			"$addToSet": bson.M{
				"sources": source,
			},
		}, models.Parliamentarian{})
		checkError(err)
	}
}

func getIdDeputado(ids []Identifier) (string, bool) {
	for _, id := range ids {
		if *id.Scheme == "ideCadastro" {
			return *id.Identifier, true
		}
	}
	return "", false
}
