package parser

import "github.com/camarabook/camarabook-api/models"

type Parser interface {
	Run(DB models.Database)
}
