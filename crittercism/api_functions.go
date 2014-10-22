package crittercism

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
)

func (p *CrittercismPlugin) PostGraph(job *job.Job, path, name string, interval int, f *gotelemetry.Flow) {
	err := p.api.FetchGraphIntoFlow(path, name, interval, f)

	if err != nil {
		job.ReportError(err)
		return
	}

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
}

func (p *CrittercismPlugin) PostLastValueOfGraph(job *job.Job, path, name string, interval int, f *gotelemetry.Flow) {
	err := p.api.FetchLastValueOfGraphIntoFlow(path, name, interval, f)

	if err != nil {
		job.ReportError(err)
		return
	}

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
}

// Daily active users

func (p *CrittercismPlugin) DailyActiveUsers(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraph(job, "errorMonitoring/graph", "dau", 86400, f)
}

func (p *CrittercismPlugin) DailyUsers(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "dau", 1440, f)
}

// Monthly active users

func (p *CrittercismPlugin) MonthlyActiveUsers(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraph(job, "errorMonitoring/graph", "mau", 30*86400, f)
}

func (p *CrittercismPlugin) MonthlyUsers(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "mau", 86400, f)
}

// Daily app loads

func (p *CrittercismPlugin) DailyAppLoads(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "appLoads", 1440, f)
}

// Crashes

func (p *CrittercismPlugin) DailyAppCrashes(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "crashes", 1440, f)
}

func (p *CrittercismPlugin) DailyCrashRate(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "crashPercent", 1440, f)
}
