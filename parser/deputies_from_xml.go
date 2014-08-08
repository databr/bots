package parser

import (
	"log"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/jinzhu/gorm"
)

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
