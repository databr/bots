package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/databr/api/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	CAMARABASEURL      = "http://www.camara.gov.br/"
	QUOATAANALITICOURL = "http://www.camara.gov.br/cota-parlamentar/cota-analitico?nuDeputadoId=ID&numMes=MONTH&numAno=YEAR&numSubCota="
)

type SaveDeputiesQuotas struct {
}

func (p SaveDeputiesQuotas) Run(DB models.Database) {
	url := "http://www.camara.gov.br/cota-parlamentar/pg-cota-lista-deputados.jsp"

	if isCached(url) {
		return
	}
	defer deferedCache(url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
		return
	}

	doc.Find("#content ul li a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		name_party := strings.Split(s.Text(), "-")
		name := strings.TrimSpace(name_party[0])
		id := models.MakeUri(name)
		if !strings.Contains(id, "lideranca") {
			getPages(CAMARABASEURL+url, id, DB)
		}
	})
}

func getPages(url, id string, DB models.Database) {
	if isCached(url) {
		return
	}
	defer deferedCache(url)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Critical("Problems %s", url)
		return
	}

	_id := strings.Split(url, "nuDeputadoId=")

	doc.Find("#mesAno option").Each(func(_ int, s *goquery.Selection) {
		monthYear, ok := s.Attr("value")
		if ok {
			dateData := strings.Split(monthYear, "-")
			if dateData[1] == "2014" {
				fullQuotasUrl := QUOATAANALITICOURL
				fullQuotasUrl = strings.Replace(fullQuotasUrl, "MONTH", dateData[0], -1)
				fullQuotasUrl = strings.Replace(fullQuotasUrl, "YEAR", dateData[1], -1)
				fullQuotasUrl = strings.Replace(fullQuotasUrl, "ID", _id[1], -1)
				getQuotaPage(id, fullQuotasUrl, DB)
			}
		}
	})
}

func getQuotaPage(id, url string, DB models.Database) {
	if isCached(url) {
		return
	}
	defer deferedCache(url)

	<-time.After(2 * time.Second)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Error(err.Error(), url)
		return
	}

	var p models.Parliamentarian
	DB.FindOne(bson.M{
		"id": id,
	}, &p)

	doc.Find(".espacoPadraoInferior2 tr:not(.celulasCentralizadas)").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		cnpj := data.Eq(0).Text()

		if cnpj == "TOTAL" {
			return
		}
		suplier := data.Eq(1).Text()
		orderN := strings.TrimSpace(data.Eq(2).Text())
		companyUri := models.MakeUri(suplier)

		if cnpj == "" {
			cnpj = companyUri
		}

		_, err := DB.Upsert(bson.M{"id": cnpj}, bson.M{
			"$set": bson.M{
				"name": suplier,
				"uri":  companyUri,
			},
		}, models.Company{})
		checkError(err)

		switch len(data.Nodes) {
		case 4:
			// value :=  data.Eq(3).Text()
			// log.Println("normal:", cnpj, "|", suplier, "|", orderN, value)
			// log.Println("skip")
		case 7:
			sendedAt, _ := time.Parse("2006-01-02", strings.Split(data.Eq(3).Text(), " ")[0])
			value := strings.Replace(data.Eq(6).Text(), "R$", "", -1)
			value = strings.Replace(value, ".", "", -1)
			value = strings.Replace(value, "-", "", -1)
			value = strings.TrimSpace(strings.Replace(value, ",", ".", -1))
			valueF, _ := strconv.ParseFloat(value, 64)

			log.Debug(orderN)

			orderNS := strings.Split(orderN, ":")
			var ticket string
			if len(orderNS) == 1 {
				ticket = strings.TrimSpace(orderNS[0])
			} else {
				ticket = strings.TrimSpace(orderNS[1])
			}

			_, err = DB.Upsert(bson.M{"order": orderN, "parliamentarian": p.Id}, bson.M{
				"$set": bson.M{
					"company":        cnpj,
					"date":           sendedAt,
					"passenger_name": data.Eq(4).Text(),
					"route":          data.Eq(5).Text(),
					"value":          valueF,
					"ticket":         ticket,
				},
			}, models.Quota{})
			checkError(err)

		default:
			panic(data.Text())
		}
	})
}
