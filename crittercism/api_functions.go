package crittercism

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"time"
)

const tableDataLength = 80

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

func (p *CrittercismPlugin) DailyActiveUsersGraph(job *job.Job, f *gotelemetry.Flow) {
	p.PostGraph(job, "errorMonitoring/graph", "dau", 86400, 86400*30, f)
}

func (p *CrittercismPlugin) DailyActiveUsersValue(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "dau", 1440, f)
}

func (p *CrittercismPlugin) MonthlyActiveUsersValue(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "mau", 86400, f)
}

func (p *CrittercismPlugin) DailyAppLoadsValue(job *job.Job, f *gotelemetry.Flow) {
	p.PostLastValueOfGraph(job, "errorMonitoring/graph", "appLoads", 1440, f)
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

func (p *CrittercismPlugin) DailyMostFrequentCrashes(job *job.Job, f *gotelemetry.Flow) {
	data, found := f.TableData()

	if !found {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract table data from flow"+f.Tag))
	}

	crashes, err := p.api.FetchCrashStatus()

	if err != nil {
		job.ReportError(err)
		return
	}

	crashes = crashes.Aggregate()

	cells := [][]gotelemetry.TableCell{}

	var count = 8

	for _, crash := range crashes {
		name := ""

		if crash.Reason != "" {
			name = crash.Reason
		} else if crash.DisplayReason != nil {
			name = *crash.DisplayReason
		} else if crash.Name != nil {
			name = *crash.Name
		} else {
			name = "N/A (" + crash.Reason + ")"
		}

		if len(name) > tableDataLength {
			name = name[:tableDataLength-1]
		}

		cells = append(
			cells,
			[]gotelemetry.TableCell{
				gotelemetry.TableCell{Value: name},
				gotelemetry.TableCell{Value: crash.SessionCount},
			},
		)

		count -= 1

		if count == 0 {
			break
		}
	}

	for count > 0 {
		cells = append(
			cells,
			[]gotelemetry.TableCell{
				gotelemetry.TableCell{Value: ""},
				gotelemetry.TableCell{Value: ""},
			},
		)

		count -= 1
	}

	data.Cells = cells

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
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

		count := 10

		for _, slice := range slices {
			bar := gotelemetry.BarchartBar{}

			bar.Color = "#267288"
			bar.Label = slice["label"].(string)
			bar.Value = slice["value"].(float64)

			bars = append(bars, bar)

			count -= 1

			if count == 0 {
				break
			}
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

	data, found := f.ValueData()

	if !found {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract value data from flow"+f.Tag))
		return
	}

	finalRating := ratings[p.ratingKey]
	finalCount := float64(counts[p.ratingKey])

	data.Value = finalRating / finalCount

	switch p.ratingKey {
	case "ios":
		data.Icon = "fa-apple"

	case "android":
		data.Icon = "fa-android"

	case "wp":
		data.Icon = "fa-windows"

	case "html5":
		data.Icon = "fa-html5"
	}

	job.PostFlowUpdate(f)
	job.Logf("Updated flow %s", f.Tag)
}

func (p *CrittercismPlugin) SetAppName(job *job.Job, f *gotelemetry.Flow) {
	data, found := f.TextData()

	if !found {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract text data from flow"+f.Tag))
		return
	}

	data.Text = p.appName

	job.PostFlowUpdate(f)
	job.Logf("Set app name to flow %s", f.Tag)
}

func (p *CrittercismPlugin) SetDate(job *job.Job, f *gotelemetry.Flow) {
	data, found := f.TextData()

	if !found {
		job.ReportError(gotelemetry.NewError(400, "Cannot extract text data from flow"+f.Tag))
		return
	}

	data.Text = time.Now().Format("Monday, January 2")

	job.PostFlowUpdate(f)
	job.Logf("Set app name to flow %s", f.Tag)
}
