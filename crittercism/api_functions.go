package crittercism

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"strings"
)

func (p *CrittercismPlugin) PostGraph(job *job.Job, path, name string, interval int, scale int, f *gotelemetry.Flow) {
	err := p.api.FetchGraphIntoFlow(path, name, interval, scale, f)

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

// DAU/MAU/Loads

func (p *CrittercismPlugin) DailyActiveUsers(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraph(job, "errorMonitoring/graph", "dau", 86400, 86400*30, f)
}

func (p *CrittercismPlugin) DailyMonthlyLoadsUsers(job *job.Job, f *gotelemetry.Flow) {
	dau, err := p.api.FetchLastValueOfGraph("errorMonitoring/graph", "dau", 1440)

	if err != nil {
		job.ReportError(err)
		return
	}

	mau, err := p.api.FetchLastValueOfGraph("errorMonitoring/graph", "mau", 86400)

	if err != nil {
		job.ReportError(err)
		return
	}

	loads, err := p.api.FetchLastValueOfGraph("errorMonitoring/graph", "appLoads", 1440)

	if err != nil {
		job.ReportError(err)
		return
	}

	data, success := f.MultivalueData()

	if !success {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract multivalue data from flow"+f.Tag))
		return
	}

	data.Values[0].Value = loads
	data.Values[1].Value = dau
	data.Values[2].Value = mau

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
}

// Crashes

func (p *CrittercismPlugin) DailyAppCrashes(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "crashes", 1440, f)
}

func (p *CrittercismPlugin) DailyCrashRate(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "crashPercent", 1440, f)
}

// App Service Error Rates

func (p *CrittercismPlugin) PostGraphToBarchart(job *job.Job, path, name, groupBy string, interval int, f *gotelemetry.Flow) {
	if data, found := f.BarchartData(); found == true {
		jq, err := p.api.FetchGraphRaw(path, name, groupBy, interval)

		if err != nil {
			job.ReportError(err)
			return
		}

		slices, err := jq.ArrayOfObjects("data", "slices")

		if err != nil {
			job.ReportError(err)
			return
		}

		bars := []gotelemetry.BarchartBar{}

		for _, slice := range slices {
			bar := gotelemetry.BarchartBar{}

			bar.Color = "#267288"
			bar.Label = slice["label"].(string)
			bar.Value = slice["value"].(float64)

			bars = append(bars, bar)
		}

		data.Bars = bars

		job.PostFlowUpdate(f)
		job.Logf("Updated flow %s", f.Tag)

		return
	}

	job.ReportError(gotelemetry.NewError(400, "Cannot extract barchart data from flow"+f.Tag))
}

func (p *CrittercismPlugin) AppServiceErrorRates(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraphToBarchart(job, "performanceManagement/pie", "errors", "service", 60, f)
}

// App Loads by device

func (p *CrittercismPlugin) AppLoadsByDevice(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraphToBarchart(job, "errorMonitoring/pie", "appLoads", "device", 1440, f)
}

// Crash by OS

func (p *CrittercismPlugin) CrashesByOS(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraphToBarchart(job, "errorMonitoring/pie", "crashes", "os", 1440, f)
}

// App Store Ratings

func (p *CrittercismPlugin) AppStoreRatings(job *job.Job, f *gotelemetry.Flow) {
	jq, err := p.api.Request("GET", "apps?attributes=appType,rating", nil)

	if err != nil {
		job.ReportError(err)
		return
	}

	source, err := jq.Object()

	if err != nil {
		job.ReportError(err)
		return
	}

	ratings := map[string]float64{}
	counts := map[string]int{}

	for _, appObj := range source {
		app := appObj.(map[string]interface{})

		t := app["appType"].(string)
		rating := app["rating"].(float64)

		counts[t] += 1
		ratings[t] += rating
	}

	icons := []gotelemetry.IconIcon{}

	for os, count := range counts {
		if count > 0 {
			rating := int(ratings[os])

			icon := gotelemetry.IconIcon{
				Label: strings.Repeat("★", rating) + strings.Repeat("☆", 5-rating),
				Color: "rgb(212, 212, 212)",
			}

			switch os {
			case "ios":
				icon.Type = "fa-apple"

			case "android":
				icon.Type = "fa-android"

			case "wp":
				icon.Type = "fa-windows"

			case "html5":
				icon.Type = "fa-html5"
			}

			icons = append(icons, icon)
		}
	}

	data, success := f.IconData()

	if !success {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract icon data from flow"+f.Tag))
		return
	}

	data.Icons = icons

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
}
