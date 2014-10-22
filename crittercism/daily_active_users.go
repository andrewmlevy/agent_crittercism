package crittercism

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

//DailyActiveUsers will query the Crittercism API and get the daily active users
//It will then emit a Flow object to the flowChan for sending up to Telemetry
func (p *CrittercismPlugin) DailyActiveUsers(job *job.Job, f *gotelemetry.Flow) {
	if data, found := f.GraphData(); found == true {

		series, err := p.api.FetchGraph("errorMonitoring/graph", "dau", 86400)

		if err != nil {
			job.ReportError(err)
			return
		}

		data.Series[0].Values = series

		job.PostFlowUpdate(f)

		job.Logf("Updated flow %s", f.Tag)
	} else {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract value data from flow"+f.Tag))
	}

}
