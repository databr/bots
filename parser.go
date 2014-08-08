package main

import (
	"fmt"
	"os"

	. "github.com/camarabook/camarabook-data/parser"
)

var usage = `Usage: camarabook-data <parsers>...

Available parsers:

    --save-deputies-from-search  Save deputies from oficial site search
    --save-deputies-from-xml     Save deputies from oficial site xml
`

var mapp = map[string]Parser{
	"--save-deputies-from-search": SaveDeputiesFromSearch{},
	"--save-deputies-from-xml":    SaveDeputiesFromXML{},
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println(usage)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println(usage)
	}

	parsers := os.Args[1:]

	for i, _ := range parsers {
		Run(mapp[parsers[i]])
	}
}

func Run(p Parser) {
	p.Run()
}
