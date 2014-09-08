package parser

import (
	"github.com/databr/api/models"
	"github.com/databr/go-popolo"
	"github.com/dukex/go-transparencia"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromTransparenciaBrasil struct {
}

func (p SaveDeputiesFromTransparenciaBrasil) Run(DB models.Database) {
	source := popolo.Source{
		Url:  toPtr("http://dev.transparencia.org.br/"),
		Note: toPtr("Transparencia Brasil"),
	}

	if isCached("http://dev.transparencia.org.br/") {
		return
	}
	defer deferedCache("http://dev.transparencia.org.br/")

	log.Info("Starting SaveDeputiesFromTransparenciaBrasil")

	c := transparencia.New("kqOfbdNKSlpf")
	query := map[string]string{
		"casa": "1",
	}
	parliamenrians, err := c.Excelencias(query)
	checkError(err)

	for _, parliamenrian := range parliamenrians {
		uri := models.MakeUri(parliamenrian.Apelido)
		log.Info("Saving %s", parliamenrian.Nome)

		_, err := DB.Upsert(bson.M{"id": uri}, bson.M{
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"summary":          parliamenrian.MiniBio,
				"nationalidentify": parliamenrian.CPF,
			},
			"$addToSet": bson.M{
				"sources": source,
				"identifiers": bson.M{
					"$each": []bson.M{
						{
							"identifier": parliamenrian.Id,
							"scheme":     "TransparenciaBrasilID",
						},
						{
							"identifier": parliamenrian.CPF,
							"scheme":     "CPF",
						},
					},
				},
			},
		}, models.Parliamentarian{})
		checkError(err)
	}
}
