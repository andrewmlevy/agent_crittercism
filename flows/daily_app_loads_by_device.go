package flows

import (
	"fmt"
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

// DailyAppLoadsByDevice will query the Crittercism API and get the daily app loads by device
// It will then emit a Flow object to the flowChan for sending up to Telemetry
func DailyAppLoadsByDevice(flowChan chan gotelemetry.Flow, config map[string]string) {

	const interval = 3600
	const label = "Daily App Loads By Device"
	const tag = "daily_app_loads_by_device"

	// Loop indefinitely, querying every interval seconds
	for {
		loopStartTime := time.Now()

		// Build the Request of Crittercism
		params := fmt.Sprintf(`{"params":{"groupBy": "service", "graph": "errors", "duration": 60, "appId": "%s"}}`, config["appId"])
		path := "performanceManagement/pie"

		// Get the data from Crittercism
		if jq, err := crittercism.Request("POST", path, params, config); err == nil {

			// Parse the result
			value, _ := jq.Float("data", "series", "0", "points", "0")
			// TODO

			// Construct the flow then emit to the channel
			data := gotelemetry.Value{Label: label, Value: value}
			flow := gotelemetry.NewFlow(tag, &data)
			flowChan <- *flow

		} else {
			log.Print("Error ", tag, err) // Request Error
		}

		// Sleep for the next loop
		if sleepTime := interval - time.Since(loopStartTime).Seconds(); sleepTime > 0 {
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}
}