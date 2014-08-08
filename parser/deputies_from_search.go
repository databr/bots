package parser

import (
	"log"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/jinzhu/gorm"
)

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
