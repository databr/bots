package parser

import (
	"log"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/jinzhu/gorm"
)

type Parser interface {
	Run(DB gorm.DB)
}

type SaveDeputiesFromSearch struct {
}

func (p SaveDeputiesFromSearch) Run(DB gorm.DB) {
	searchURL := "http://www2.camara.leg.br/deputados/pesquisa"

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(searchURL); e != nil {
		log.Fatal(e)
	}

	doc.Find("#deputado option").Each(func(i int, s *goquery.Selection) {
		value, _ := s.Attr("value")
		if value != "" {
			info := regexp.MustCompile("=|%23|!|\\||\\?").Split(value, -1)

			var party models.Party
			DB.Where(models.Party{Title: info[4]}).FirstOrCreate(&party)

			RegisterID, _ := strconv.Atoi(info[5])

			DB.Where(models.Parliamentarian{
				RegisterId: int64(RegisterID),
			}).Assign(models.Parliamentarian{
				Name:       info[0],
				PartyId:    party.Id,
				State:      info[3],
				RegisterId: int64(RegisterID),
			}).FirstOrCreate(&models.Parliamentarian{})
		}
	})
}

// ---

type SaveDeputiesFromXML struct{}

func (p SaveDeputiesFromXML) Run(DB gorm.DB) {
	xmlURL := "http://www.camara.gov.br/SitCamaraWS/Deputados.asmx/ObterDeputados"

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(xmlURL); e != nil {
		log.Fatal(e)
	}

	doc.Find("deputado").Each(func(i int, s *goquery.Selection) {
		govID, _ := strconv.Atoi(s.Find("idParlamentar").First().Text())
		registerID, _ := strconv.Atoi(s.Find("ideCadastro").First().Text())

		var party models.Party
		DB.Where(models.Party{Title: s.Find("partido").First().Text()}).FirstOrCreate(&party)

		DB.Where(models.Parliamentarian{
			GovId: int64(govID),
		}).Assign(models.Parliamentarian{
			Name:       s.Find("nomeparlamentar").First().Text(),
			PartyId:    party.Id,
			GovId:      int64(govID),
			Gender:     s.Find("sexo").First().Text(),
			Phone:      s.Find("fone").First().Text(),
			Email:      s.Find("email").First().Text(),
			RegisterId: int64(registerID),
			RealName:   s.Find("nome").First().Text(),
			State:      s.Find("uf").First().Text(),
		}).FirstOrCreate(&models.Parliamentarian{})
	})
}

// ---
