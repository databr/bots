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

type SabespBot struct{}

func (_ SabespBot) Run(db database.MongoDB) {
	getData("d", db)
	getData("s", db)
	getData("q", db)
	getData("m", db)
	getData("x", db)
	getData("a", db)
}

func getData(g string, db database.MongoDB) {
	url := "http://www.apolo11.com/reservatorios.php?step=" + g
	doc, err := goquery.NewDocument(url)
	parser.CheckError(err)
	doc.Find("body > center:nth-child(1) > table > tbody > tr > td:nth-child(1) > b > table").Each(func(_ int, s *goquery.Selection) {
		trs := s.Find("tr")
		title := "Sistema " + parser.Titlelize(strings.Replace(trs.Eq(0).Text(), "SISTEMA", "", -1))
		uri := models.MakeUri(title)

		getInfo := func(i int, ss *goquery.Selection) string {
			return ss.Text()
		}

		percent := trs.Eq(1).Find("font").Map(getInfo)
		date := trs.Eq(2).Find("font").Map(getInfo)

		data := make([]bson.M, 0)
		for i, _ := range percent {
			if strings.TrimSpace(date[i]) != "" && strings.TrimSpace(date[i]) != "/" {
				data = append(data, bson.M{"percent": percent[i], "date": date[i]})
			}
		}

		query := bson.M{"uri": uri, "granularity_letter": g}

		source := models.Source{
			Url:  "http://www.apolo11.com",
			Note: "Apolo11",
		}

		_, err := db.Upsert(query, bson.M{
			"$setOnInsert": bson.M{
				"createdat": time.Now(),
			},
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"uri":                uri,
				"name":               title,
				"granularity_letter": g,
				"granularity":        getGranularity(g),
				"data":               data,
				"source":             []models.Source{source},
			},
		}, models.Reservoir{})
		parser.CheckError(err)

		log.Println(uri, title, data)
	})

}

func getGranularity(letter string) string {
	return map[string]string{
		"d": "daily",
		"s": "weekly",
		"q": "fortnightly",
		"m": "monthly",
		"x": "semiannual",
		"a": "annual",
	}[letter]
}
