package parser

import (
	"fmt"

	"github.com/wsxiaoys/terminal"
)

type llog struct {
}

func (l llog) Log(v ...interface{}) {
	message := fmt.Sprintln(v...)
	l.log("y", "LOG", message)
}

func (l llog) Debug(v ...interface{}) {
	message := fmt.Sprintln(v...)
	l.log("b", "DEBUG", message)
}

func (l llog) Error(v ...interface{}) {
	message := fmt.Sprintln(v...)
	l.log("r", "ERROR", message)
}

func (l llog) log(color, prefix, message string) {
	terminal.Stdout.Color(color).
		Print(prefix + ": " + message).Reset()
}

var LLog *llog

func init() {
	LLog = new(llog)
}

//func main() {
//.Color("y").
//Print("Hello world").Nl().
//Reset().
//Colorf("@{kW}Hello world\n")

//color.Print("@rHello world")
//}
