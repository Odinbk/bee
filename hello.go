package main

import (
	"encoding/json"
	"fmt"
	C "github.com/go_proj/proj/core"
	"os"
	"time"
)

type Settings struct {
	TestURL string
	Params  map[string]string
	Header  map[string]string
	SSL     bool
}

var settings Settings

func main() {
	get_settings(&settings)

	foo := func() (task C.Runner) {
		task = &C.HttpTask{
			ID:     1,
			URL:    settings.TestURL,
			Engine: C.Engine{},
			Params: settings.Params,
			Header: settings.Header,
			SSL:    settings.SSL,
		}
		return task
	}

	buffer := 10
	consumerCount := 10
	c := C.Bee{
		QPS:           buffer,
		Duration:      5 * time.Second,
		Timeout:       time.Second,
		Tasks:         make(chan C.Runner, buffer),
		Medium:        make(chan C.Runner),
		ProducerCount: 1,
		ConsumerCount: consumerCount,
		TaskHandler:   C.GetTask(foo),
	}
	c.Launcher()
}

func get_settings(settings *Settings) {
	file, _ := os.Open("settings.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(settings)
	if err != nil {
		fmt.Println("Decode settings file error:", err)
	}
}
