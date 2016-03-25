package bot

import (
	"time"

	"github.com/databr/api/database"
	"github.com/databr/api/models"
	"github.com/databr/bots/go_bot/parser"
	"github.com/lucasb-eyer/go-colorful"
	"gopkg.in/mgo.v2/bson"
)

type LineBot struct{}

func (_ LineBot) Run(db database.MongoDB) {
	LineColor("linha1azul", "#1a5ba3", db)
	LineColor("linha2verde", "#008569", db)
	LineColor("linha3vermelha", "#f04d43", db)
	LineColor("linha4amarela", "#ffd527", db)
	LineColor("linha5lilas", "#a84f9c", db)
	LineColor("linha7rubi", "#B81D64", db)
	LineColor("linha8diamante", "#8A8988", db)
	LineColor("linha9esmeralda", "#009496", db)
	LineColor("linha10turquesa", "#0088B0", db)
	LineColor("linha11coral", "#E87A65", db)
	LineColor("linha12safira", "#1C2D72", db)
}

func LineColor(uri, hex string, db database.MongoDB) {
	q := bson.M{"id": uri}

	c, _ := colorful.Hex(hex)
	r, g, b := c.RGB255()
	color := bson.M{
		"hex": hex,
		"rgb": []int{int(r), int(g), int(b)},
	}

	parser.Log.Debug("Save", uri, "with color", color)

	_, err := db.Upsert(q, bson.M{
		"$setOnInsert": bson.M{
			"createdat": time.Now(),
		},
		"$currentDate": bson.M{
			"updatedat": true,
		},
		"$set": bson.M{
			"color": color,
		},
	}, models.Line{})
	parser.CheckError(err)
}
