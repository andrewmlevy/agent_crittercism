package flows

import (
	"fmt"
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

// DailyCrashRate will query the Crittercism API and get the daily crash rate
// It will then emit a Flow object to the flowChan for sending up to Telemetry
func DailyCrashRate(flowChan chan gotelemetry.Flow, config map[string]string) {

	const interval = 3600
	const label = "Daily Crash Rate"
	const tag = "daily_crash_rate"

	// Loop indefinitely, querying every interval seconds
	for {
		loopStartTime := time.Now()

		// Build the Request of Crittercism
		params := fmt.Sprintf(`{"params":{"graph": "crashPercent", "value_type": "percent", "duration": 1440, "appId": "%s"}}`, config["appId"])
		path := "errorMonitoring/graph"

		// Get the data from Crittercism
		if jq, err := crittercism.Request("POST", path, params, config); err == nil {

			// Parse the result
			value, _ := jq.Float("data", "series", "0", "points", "0")

			// Construct the flow then emit to the channel
			data := gotelemetry.Value{Label: label, Value: value, ValueType: "percent"}
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
