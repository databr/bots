package parser

import "github.com/jinzhu/gorm"

type Parser interface {
	Run(DB gorm.DB)
}
