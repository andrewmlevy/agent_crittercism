package flows

import (
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"log"
	"time"
)

func DailyAppLoads(flowChan chan gotelemetry.Flow) {

	for {
		log.Print("looping")

		jq, err := crittercism.Request(
			"POST",
			"errorMonitoring/graph",
			`{"params":{"graph": "appLoads","duration": 1440,"appId": "519d53101386202089000007"}}`)

		if err == nil {
			loads, _ := jq.Int("data", "series", "0", "points", "0")
			log.Print("Loads: ", loads)

      data := gotelemetry.MultivalueValue{Label: "WEB", Value: www.rps, ValueType: "rps"}

      flow := 
		}

		time.Sleep(60 * time.Second)
	}

}
