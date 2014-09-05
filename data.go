package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/camarabook/camarabook-api/models"
	. "github.com/databr/bots/parser"
)

var mapp = []Parser{
	SaveDeputiesFromSearch{},
	SaveDeputiesFromXML{},
	SaveDeputiesAbout{},
	// SaveDeputiesQuotas{},
	SaveDeputiesFromTransparenciaBrasil{},
	SaveSenatorsFromIndex{},
}

var DB models.Database

func main() {
	StartDispatcher(6)
	DB = models.New()

	max := len(mapp)
	c := 0

	go func() {
		for {
			if c == max {
				time.Sleep(1 * time.Hour)
				c = 0
			}
			Collector(mapp[c], reflect.ValueOf(mapp[c]).Type().Name())
			c++
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Finishing....")
	close(WorkerQueue)
}

// ---

type WorkRequest struct {
	Name   string
	Parser Parser
	Delay  time.Duration
}

var WorkQueue = make(chan WorkRequest, 100)

func Collector(parser Parser, name string) {
	delay := time.Second * 8
	work := WorkRequest{Parser: parser, Delay: delay, Name: name}
	WorkQueue <- work
	fmt.Println("Work", name, "queued")
	return
}

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}

	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				fmt.Printf("worker%d: Received work request, delaying for %f seconds, to %s\n", w.ID, work.Delay.Seconds(), work.Name)

				time.Sleep(work.Delay)
				fmt.Printf("worker%d: Hello, %s!\n", w.ID, work.Name)
				work.Parser.Run(DB)
			case <-w.QuitChan:
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

var WorkerQueue chan chan WorkRequest

func StartDispatcher(nworkers int) {
	WorkerQueue = make(chan chan WorkRequest, nworkers)
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received", work.Name, "requeust", work.Name)
				go func() {
					worker := <-WorkerQueue

					fmt.Println("Dispatching", work.Name, "request", work.Name)
					worker <- work
				}()
			}
		}
	}()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
