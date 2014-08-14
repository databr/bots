package parser

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/camarabook/camarabook-api/models"
)

const (
	CAMARABASEURL      = "http://www.camara.gov.br/"
	QUOATAANALITICOURL = "http://www.camara.gov.br/cota-parlamentar/cota-analitico?nuDeputadoId=ID&numMes=MONTH&numAno=YEAR&numSubCota="
)

type SaveDeputiesQuotas struct {
}

func (p SaveDeputiesQuotas) Run(DB models.Database) {
	doc, err := goquery.NewDocument("http://www.camara.gov.br/cota-parlamentar/pg-cota-lista-deputados.jsp")
	if err != nil {
		log.Println("Problems base")
		return
	}

	doc.Find("#content ul li a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		getPages(CAMARABASEURL + url)
	})
}

func getPages(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Problems", url)
		return
	}

	_id := strings.Split(url, "nuDeputadoId=")

	doc.Find("#mesAno option").Each(func(_ int, s *goquery.Selection) {
		monthYear, ok := s.Attr("value")
		if ok {
			dateData := strings.Split(monthYear, "-")
			fullQuotasUrl := QUOATAANALITICOURL
			fullQuotasUrl = strings.Replace(fullQuotasUrl, "MONTH", dateData[0], -1)
			fullQuotasUrl = strings.Replace(fullQuotasUrl, "YEAR", dateData[1], -1)
			fullQuotasUrl = strings.Replace(fullQuotasUrl, "ID", _id[1], -1)
			getQuotaPage(fullQuotasUrl)
		}
	})
}

func getQuotaPage(url string) {
	<-time.After(5 * time.Second)
	log.Println("aaaa, ", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("Problems", url)
		return
	}

	doc.Find(".espacoPadraoInferior2 tr:not(.celulasCentralizadas)").Each(func(i int, s *goquery.Selection) {
		data := s.Find("td")
		cnpj := data.Eq(0).Text()

		if cnpj == "TOTAL" {
			return
		}
		suplier := data.Eq(1).Text()
		orderN := strings.TrimSpace(data.Eq(2).Text())

		switch len(data.Nodes) {
		case 4:
			value := data.Eq(3).Text()
			log.Println("normal:", cnpj, "|", suplier, "|", orderN, value)
		case 7:
			date := data.Eq(3).Text()
			passenger := data.Eq(4).Text()
			path := data.Eq(5).Text()
			value := data.Eq(6).Text()
			log.Println("aereo:", cnpj, "|", suplier, "|", orderN, "|", date, "|", passenger, "|", path, "|", value)
		default:
			log.Println("outro:", data.Text())
		}
	})
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
