package bot

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"gopkg.in/mgo.v2/bson"
)

var re = regexp.MustCompile(`(\d{1,2})`)

type StatusBot struct{}

func (_ StatusBot) Run(db database.MongoDB) {
	// Metro SP
	doc, err := goquery.NewDocument("http://www.metro.sp.gov.br/Sistemas/direto-do-metro-via4/diretodoMetroHome.aspx")
	parser.CheckError(err)
	metroSource := models.Source{
		Url: "http://www.metro.sp.gov.br/Sistemas/direto-do-metro-via4/diretodoMetroHome.aspx",
	}

	doc.Find("ul li").Each(func(_ int, s *goquery.Selection) {
		lineName := strings.TrimSpace(s.Find(".nomeDaLinha").Text())
		status := strings.TrimSpace(s.Find(".statusDaLinha").Text())
		saveStatus(db, lineName, status, metroSource)
	})

	// CPTM

	doc, err = goquery.NewDocument("http://www.cptm.sp.gov.br/Pages/atendimento.aspx")
	parser.CheckError(err)
	cptmSource := models.Source{
		Url: "http://www.cptm.sp.gov.br/Central-Relacionamento/situacao-linhas.asp",
	}

	lines := map[string]string{
		"rubi":      "7",
		"diamante":  "8",
		"esmeralda": "9",
		"turquesa":  "10",
		"coral":     "11",
		"safira":    "12",
	}

	doc.Find("#atendimento_consumidor .situacao_linhas .col-md-10 div.col-md-2").Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		uri := strings.TrimSpace(strings.Replace(class, "col-md-2", "", -1))
		name := s.Find(".nome_linha")
		lineNumber := lines[uri]
		status := s.Find("[data-toggle='tooltip']").Text()

		lineName := "Linha " + lineNumber + " - " + parser.Titlelize(strings.TrimSpace(name.Text()))

		saveStatus(db, lineName, status, cptmSource)
	})
}

func saveStatus(db database.MongoDB, lineName, status string, source models.Source) {
	uri := models.MakeUri(lineName)
	result := re.FindStringSubmatch(lineName)
	lineNumber, _ := strconv.Atoi(result[0])

	q := bson.M{"id": uri}

	_, err := db.Upsert(q, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"name":       lineName,
			"linenumber": lineNumber,
		},
		"$addToSet": bson.M{
			"sources": source,
		},
	}, models.Line{})

	log.Println(uri)
	parser.CheckError(err)
	var statusOld models.Status
	err = db.FindOne(bson.M{"line_id": uri}, &statusOld)

	statusQ := bson.M{"line_id": uri, "_id": bson.NewObjectId()}
	if err == nil && statusOld.Status == status {
		statusQ = bson.M{"_id": statusOld.Id, "line_id": uri}
	}

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
		"$addToSet": bson.M{
			"sources": source,
		},
	}, models.Status{})
	parser.CheckError(err)

	parser.Log.Debug(lineName + " - " + status)
	parser.Log.Info("-- Created Status to " + lineName)
	parser.Log.Info("Status: " + status)
	parser.Log.Info("------")

	if uri == "linha11coral" {
		saveStatus(db, "Linha 11-Coral-Expresso", status, source)
	}
}
