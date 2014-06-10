package main

import (
	"github.com/go-martini/martini"
	"log"
	"time"
)

// The one and only martini instance.
var m *martini.Martini

func init() {
	m = martini.New()
	//setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	m.Use(martini.Static("public"))
}

func schedule() {
	var (
		count   float64
		average float64
	)
	count = 0
	average = 0

	dao := NewDAO()
	defer dao.Close()

	for {
		t, err := ReadTemperature(5)
		if err != nil {
			log.Fatalln("Read temperature error", err)
		}
		average = average + t
		time.Sleep(time.Millisecond * 10)
		count = count + 1
		if count >= 100 { // every 1 second
			m := (average / count)
			count = 0
			average = 0
			log.Printf("temperature: %f\n", m)
		}
	}
}

func main() {
	// temperature schedule
	go schedule()

	// run martini on default port 3000
	m.Run()
}
