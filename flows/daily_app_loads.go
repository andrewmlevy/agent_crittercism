package flows

import (
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

// Our interval to update this Flow on
const interval = 60 * time.Second

// DailyAppLoads will query the Crittercism API and get the daily app loads
// It will then emit a Flow object to the flowChan for sending up to Telemetry
func DailyAppLoads(flowChan chan gotelemetry.Flow) {
	// Loop indefinitely, querying every interval seconds
	for {
		// Get the data from the API
		if jq, err := crittercism.Request(
			"POST",
			"errorMonitoring/graph",
			`{"params":{"graph": "appLoads","duration": 1440,"appId": "519d53101386202089000007"}}`); err == nil {

			if loads, err := jq.Float("data", "series", "0", "points", "0"); err == nil {
				// Parse the result and construct the flow
				data := gotelemetry.Value{Label: "Daily App Loads", Value: loads}
				flow := gotelemetry.NewFlow("daily_app_loads", &data)

				// Emit to the channel
				flowChan <- *flow

			} else {
				log.Print("Error ", err)
			}
		} else {
			log.Print("Error ", err)
		}
		time.Sleep(interval)
	}
}
