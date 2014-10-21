package flows

import (
	"fmt"
	"github.com/telemetryapp/agent_crittercism/crittercism"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

// DailyActiveUsers will query the Crittercism API and get the daily active users
// It will then emit a Flow object to the flowChan for sending up to Telemetry
func DailyActiveUsers(job *job.Job, f *gotelemetry.Flow) {

	if data, err := f.ValueData(); err == nil {
		// Build the Request of Crittercism
		params := fmt.Sprintf(`{"params":{"groupBy": "service", "graph": "errors", "duration": 60, "appId": "%s"}}`, job.Config["appId"])
		path := "performanceManagement/pie"

		// Get the data from Crittercism
		if jq, err := crittercism.Request("POST", path, params, job.Config); err == nil {

			if value, err := jq.Float("data", "series", "0", "points", "0"); err == nil {
				data.Value = value
			}
			//TODO

			job.PostFlowUpdate(f)

		} else {
			job.ReportError(err)
		}

	} else {
		job.ReportError(err)
	}

}
