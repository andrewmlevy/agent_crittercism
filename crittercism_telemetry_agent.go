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
	var config map[string]string

	//519d53101386202089000007

	// Launch our reader
	go reader(flowChan)

	// Launch each of our flow ETLs
	go flows.DailyActiveUsers(flowChan, config)
	go flows.DailyAppCrashes(flowChan, config)
	go flows.DailyAppLoads(flowChan, config)
	go flows.DailyAppLoadsByDevice(flowChan, config)
	go flows.DailyCrashRate(flowChan, config)
	go flows.DailyCrashesByOs(flowChan, config)
	go flows.MonthlyActiveUsers(flowChan, config)
	go flows.ServiceMonitoringErrorRate(flowChan, config)

	// Wait for an exit message
	<-exitChan
}
