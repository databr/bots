package parser

import "github.com/camarabook/camarabook-api/models"

type SaveDeputiesAbout struct {
}

func (p SaveDeputiesAbout) Run(DB models.Database) {
	var ds []models.Parliamentarian

	DB.Find(&ds)

	for _, d := range ds {
		dURL := "http://www2.camara.leg.br/deputados/pesquisa/layouts_deputados_biografia?pk=" + strconv.Itoa(int(d.RegisterId))

		var doc *goquery.Document
		var e error

		if doc, e = goquery.NewDocument(dURL); e != nil {
			log.Fatal(e)
		}

		doc.Find("#bioDeputado .bioOutros").Each(func(i int, s *goquery.Selection) {
			title := s.Find(".bioOutrosTitulo").Text()
			if title != "" {
				title = strings.TrimSpace(title)
				title = strings.Replace(title, ":", "", -1)

				body := s.Find(".bioOutrosTexto").Text()

				hasher := md5.New()
				hasher.Write([]byte(title))
				sectionKey := hex.EncodeToString(hasher.Sum(nil))

				DB.Where(models.ParliamentarianAbout{
					SectionKey:        string(sectionKey),
					ParliamentarianId: d.RegisterId,
				}).Assign(models.ParliamentarianAbout{
					Body:  body,
					Title: title,
				}).FirstOrCreate(&models.ParliamentarianAbout{})
			}
		})
	}
}
