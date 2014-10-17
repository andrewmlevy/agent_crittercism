package flows

import (
	"fmt"
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

// Our interval to update this Flow on
const interval = 60

// DailyAppLoads will query the Crittercism API and get the daily app loads
// It will then emit a Flow object to the flowChan for sending up to Telemetry
func DailyAppLoads(flowChan chan gotelemetry.Flow) {

	// Loop indefinitely, querying every interval seconds
	for {
		loopStartTime := time.Now()

		// Build the Request of Crittercism
		params := fmt.Sprintf(`{"params":{"graph": "appLoads", "duration": 1440, "appId": "%s"}}`, "519d53101386202089000007")
		path := "errorMonitoring/graph"

		// Get the data from Crittercism
		if jq, err := crittercism.Request("POST", path, params); err == nil {
			if loads, err := jq.Float("data", "series", "0", "points", "0"); err == nil {

				// Parse the result and construct the flow then emit to the channel
				data := gotelemetry.Value{Label: "Daily App Loads", Value: loads}
				flow := gotelemetry.NewFlow("daily_app_loads", &data)
				flowChan <- *flow

			} else {
				log.Print("Error ", err) // Parsing Error
			}
		} else {
			log.Print("Error ", err) // Request Error
		}

		// Sleep for the next loop
		if sleepTime := interval - time.Since(loopStartTime).Seconds(); sleepTime > 0 {
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}
}
