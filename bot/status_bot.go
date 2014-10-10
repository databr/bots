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

type StatusBot struct{}

func (_ StatusBot) Run(db database.MongoDB) {
	// Metro SP
	doc, err := goquery.NewDocument("http://www.metro.sp.gov.br/sistemas/direto-do-metro-via4/index.aspx")

	parser.CheckError(err)

	doc.Find("#diretoMetro ul li").Each(func(_ int, s *goquery.Selection) {
		lineName := s.Find(".linha").Text()
		status := strings.TrimSpace(s.Find(".status").Text())
		saveStatus(db, lineName, status)
	})

	doc, err = goquery.NewDocument("http://www.cptm.sp.gov.br/Central-Relacionamento/situacao-linhas.asp")
	parser.CheckError(err)

	doc.Find(".linhaStatus").Each(func(_ int, s *goquery.Selection) {
		data := s.Find("td")
		nameTD := data.Eq(0)
		status := data.Eq(2).Text()
		nameImage, _ := nameTD.Find("img").Attr("src")
		lineNumber := strings.Split(strings.Split(nameImage, "-")[1], ".")[0]

		lineName := "Linha " + lineNumber + "-" + parser.ToUtf8(parser.Titlelize(strings.TrimSpace(strings.Split(nameTD.Text(), "-")[1])))

		saveStatus(db, lineName, parser.ToUtf8(status))
	})
}

func saveStatus(db database.MongoDB, lineName, status string) {
	uri := models.MakeUri(lineName)

	q := bson.M{"id": uri}

	_, err := db.Upsert(q, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"name": lineName,
		},
	}, models.Line{})

	parser.CheckError(err)

	statusQ := bson.M{"line_id": uri}
	_, err = db.Upsert(statusQ, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"status":  status,
			"line_id": uri,
		},
	}, models.Status{})
	parser.CheckError(err)

	parser.Log.Debug(lineName + " - " + status)
	parser.Log.Info("-- Created Status to " + lineName)
	parser.Log.Info("Status: " + status)
	parser.Log.Info("------")
}
