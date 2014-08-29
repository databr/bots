package parser

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("camarabook")

var format = "%{color}%{time:15:04:05.000000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}"

func init() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	syslogBackend, err := logging.NewSyslogBackend("")

	logging.SetFormatter(logging.MustStringFormatter(format))

	if err != nil {
		log.Fatal(err)
	}
	logging.SetBackend(logBackend, syslogBackend)

	logging.SetLevel(logging.DEBUG, "example")
}
