package parser

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
	"github.com/camarabook/go-popolo"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromXML struct{}

func (p SaveDeputiesFromXML) Run(DB models.Database) {
	xmlURL := "http://www.camara.gov.br/SitCamaraWS/Deputados.asmx/ObterDeputados"

	source := popolo.Source{
		Url:  toPtr(xmlURL),
		Note: toPtr("Câmara API"),
	}

	var doc *goquery.Document
	var e error

	if doc, e = goquery.NewDocument(xmlURL); e != nil {
		log.Fatal(e)
	}

	doc.Find("deputado").Each(func(i int, s *goquery.Selection) {
		name := strings.Title(strings.ToLower(s.Find("nomeparlamentar").First().Text()))

		p := models.Parliamentarian{}
		p.Name = &name
		p.SortName = &name
		p.Id = toPtr(models.MakeUri(name))
		p.Gender = toPtr(s.Find("sexo").First().Text())
		p.Image = toPtr(s.Find("urlFoto").First().Text())
		p.Email = toPtr(s.Find("email").First().Text())

		q := bson.M{
			"email": p.Email,
		}
		DB.Upsert(q, &p)

		fullName := strings.Split(titlelize(s.Find("nome").First().Text()), " ")

		DB.Update(q, bson.M{
			"$set": bson.M{
				"sources": []popolo.Source{source},
				"identifiers": []popolo.Identifier{
					{Identifier: toPtr(s.Find("idParlamentar").First().Text()), Scheme: toPtr("idParlamentar")},
					{Identifier: toPtr(s.Find("ideCadastro").First().Text()), Scheme: toPtr("ideCadastro")},
				},
				"othernames": []popolo.OtherNames{
					{
						Name:       toPtr(titlelize(s.Find("nome").First().Text())),
						FamilyName: toPtr(fullName[len(fullName)-1:][0]),
						GivenName:  &fullName[0],
						Note:       toPtr("Nome de nascimento"),
					},
				},
				"contactdetails": []popolo.ContactDetail{
					{
						Label:   toPtr("Telefone"),
						Type:    toPtr("phone"),
						Value:   toPtr(s.Find("fone").First().Text()),
						Sources: []popolo.Source{source},
					},
					{
						Label:   toPtr("Gabinete"),
						Type:    toPtr("address"),
						Value:   toPtr(s.Find("gabinete").First().Text() + ", Anexo " + s.Find("anexo").First().Text()),
						Sources: []popolo.Source{source},
					},
				},
			},
		}, &p)
	})
}

func titlelize(s string) string {
	return strings.Title(strings.ToLower(s))
}

func toPtr(s string) *string {
	return &s
}
