package flows

import (
	"github.com/telemetryapp/gotelemetry"
)

func ExampleGauge(flowChan chan gotelemetry.Flow) {

	data := gotelemetry.Gauge{Title: "Test Gauge"}
	flow := gotelemetry.NewFlow("test_gauge", &data)

	flowChan <- *flow

}
