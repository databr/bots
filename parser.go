package main

import (
	"fmt"
	"os"

	"github.com/camarabook/camarabook-api/models"
	. "github.com/camarabook/camarabook-data/parser"
)

var usage = `Usage: camarabook-data <parsers>...

Available parsers:
    --all
    --save-deputies-from-search  Save deputies from official site search
    --save-deputies-from-xml     Save deputies from official site xml
    --save-deputies-about        Save deputies about information from official site
    --save-deputies-quotas
    --save-deputies-info-from-transparencia-brasil
    --save-senators-from-index
`

var mapp = map[string]Parser{
	"--save-deputies-from-search":                    SaveDeputiesFromSearch{},
	"--save-deputies-from-xml":                       SaveDeputiesFromXML{},
	"--save-deputies-about":                          SaveDeputiesAbout{},
	"--save-deputies-quotas":                         SaveDeputiesQuotas{},
	"--save-deputies-info-from-transparencia-brasil": SaveDeputiesFromTransparenciaBrasil{},
	"--save-senators-from-index":                     SaveSenatorsFromIndex{},
}

var DB models.Database

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println(usage)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	DB = models.New()

	parsers := os.Args[1:]

	if os.Args[1] == "--all" {
		parsers = []string{}
		for k, _ := range mapp {
			parsers = append(parsers, k)
		}
	}
	for i, _ := range parsers {
		Run(mapp[parsers[i]])
	}
}

func Run(p Parser) {
	p.Run(DB)
}
