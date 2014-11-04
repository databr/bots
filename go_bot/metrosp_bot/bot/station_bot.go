package bot

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

type StationBot struct{}

func (_ StationBot) Run(db database.MongoDB) {
	doc, err := goquery.NewDocument("http://www.metro.sp.gov.br/app/trajeto/xt/estacoesTipoXML.asp")
	parser.CheckError(err)
	doc.Find("estacao").Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("estacaoid")
		name, _ := s.Attr("nome")

		lineId, _ := s.Attr("linhaid")
		_lineName, _ := s.Attr("linha")

		typeId, _ := s.Attr("tipoid")
		typeName, _ := s.Attr("tipo")

		if typeId == "3" {
			return
		}

		lineName := "Linha " + strings.Split(_lineName, " ")[0]
		uri := models.MakeUri(lineName)
		cannonicaluri := uri
		names := strings.Split(lineName, "-")
		if len(names) == 3 {
			cannonicaluri = models.MakeUri(strings.Replace(lineName, names[2], "", -1))
		}

		lineQ := bson.M{"id": uri}
		_, err := db.Upsert(lineQ, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"name":          lineName,
				"cannonicaluri": cannonicaluri,
				"metroid":       lineId,
				"type": bson.M{
					"id":   typeId,
					"name": typeName,
				},
			},
		}, models.Line{})
		parser.CheckError(err)

		q := bson.M{"id": models.MakeUri(name)}
		_, err = db.Upsert(q, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"metroid": id,
				"name":    name,
			},
		}, models.Station{})
		parser.CheckError(err)

		log.Println(id, name, lineId, lineName, typeName)
	})
}
