package parser

import (
	"github.com/camarabook/camarabook-api/models"
	. "github.com/camarabook/go-popolo"
	"github.com/dukex/go-transparencia"
	"gopkg.in/mgo.v2/bson"
)

type SaveDeputiesFromTransparenciaBrasil struct {
}

func (p SaveDeputiesFromTransparenciaBrasil) Run(DB models.Database) {
	source := Source{
		Url:  toPtr("http://dev.transparencia.org.br/"),
		Note: toPtr("Transparencia Brasil"),
	}

	c := transparencia.New("kqOfbdNKSlpf")
	query := map[string]string{
		"casa": "1",
	}
	parliamenrians, err := c.Excelencias(query)
	checkError(err)

	for _, parliamenrian := range parliamenrians {
		uri := models.MakeUri(parliamenrian.Apelido)

		_, err := DB.Upsert(bson.M{"id": uri}, bson.M{
			"$currentDate": bson.M{
				"updatedat": true,
			},
			"$set": bson.M{
				"summary":          parliamenrian.MiniBio,
				"nationalidentity": parliamenrian.CPF,
			},
			"$addToSet": bson.M{
				"sources": source,
				"identifiers": bson.M{
					"identifier": toPtr(parliamenrian.Id),
					"scheme":     toPtr("TransparenciaBrasilID"),
				},
			},
		}, models.Parliamentarian{})
		checkError(err)
	}
}
