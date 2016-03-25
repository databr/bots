package bot

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"github.com/databr/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesAbout struct {
}

func (p SaveDeputiesAbout) Run(DB database.MongoDB) {
	var ds []models.Parliamentarian

	DB.FindAll(&ds)

	var wg sync.WaitGroup

	for _, d := range ds {
		id, ok := getIdDeputado(d.Identifiers)
		if !ok {
			continue
		}
		wg.Add(1)
		go func(_id string) {
			defer wg.Done()
			saveDeputies(_id, d, DB)
		}(id)
	}

	wg.Wait()
}

func saveDeputies(id string, d models.Parliamentarian, DB database.MongoDB) {

	bioURL := "http://www2.camara.leg.br/deputados/pesquisa/layouts_deputados_biografia?pk=" + id

	if parser.IsCached(bioURL) {
		parser.Log.Info("SaveDeputiesAbout(%s) Cached", id)
		return
	}

	source := models.Source{
		Url:  bioURL,
		Note: "Pesquisa Câmara",
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(bioURL); e != nil {
		parser.Log.Critical(e.Error())
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
	case "74447":
		year = 1936
	case "74474":
		year = 1940
	default:
		parser.Log.Debug("(%s) %s", id, birthdateA)
		if len(birthdateA) != 3 {
			parser.Log.Debug("Error, deputies without year %s", bioURL)
			year = 0
		} else {
			year, _ = strconv.Atoi(birthdateA[2])
		}
	}

	birthDate := popolo.Date{}

	if len(birthdateA) > 1 {
		month, _ := strconv.Atoi(birthdateA[1])
		day, _ := strconv.Atoi(birthdateA[0])
		loc, _ := time.LoadLocation("America/Sao_Paulo")
		birthDate = popolo.Date{time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)}
	}

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

	parser.CheckError(err)
	parser.CacheURL(bioURL)
}

func getIdDeputado(ids []models.Identifier) (string, bool) {
	for _, id := range ids {
		if id.Scheme == "ideCadastro" {
			return id.Identifier, true
		}
	}
	return "", false
}
