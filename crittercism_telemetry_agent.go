package main

import (
	"github.com/telemetryapp/agent_crittercism/flows"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

func reader(flowChan chan gotelemetry.Flow) {
	interval := 1.0

	next := time.Now()
	for {

		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(time.Millisecond * 100)
			timeout <- true
		}()

		select {
		case flow := <-flowChan:
			log.Print("New Flow To Update", flow)
			// Add the flow to the batch

		case <-timeout:
		}

		// sleep for the rest of the interval
		since := time.Since(next).Seconds()
		if since >= interval {
			next = time.Now().Add(1)
			log.Print("tick")
		}

	}
}

func main() {

	// Our channels for communication
	flowChan := make(chan gotelemetry.Flow)
	exitChan := make(chan bool)

	// Launch each of our flow ETLs
	go flows.DailyAppLoads(flowChan)

	// Launch our reader
	go reader(flowChan)

	// Wait for an exit message
	<-exitChan
}
