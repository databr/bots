package bot

import "github.com/databr/api/database"

type LineBot struct{}

func (_ LineBot) Run(db database.MongoDB) {
}
