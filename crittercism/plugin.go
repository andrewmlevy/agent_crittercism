package crittercism

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"time"
)

func init() {
	job.RegisterPlugin("crittercism", CrittercismPluginFactory)
}

func CrittercismPluginFactory() job.PluginInstance {
	return &CrittercismPlugin{
		job.NewPluginHelper(),
		nil,
	}
}

type CrittercismPlugin struct {
	*job.PluginHelper
	api *CrittercismAPIClient
}

type crittercismPluginClosure struct {
	closure  job.PluginHelperClosureWithFlow
	interval int
	tag      string
}

func (p *CrittercismPlugin) registerClosures(b *gotelemetry.Board, closures []crittercismPluginClosure) error {
	for _, closure := range closures {
		if err := p.PluginHelper.AddTaskWithClosureFromBoardForFlowWithTag(closure.closure, time.Second*time.Duration(closure.interval), b, closure.tag); err != nil {
			return err
		}
	}

	return nil
}

func (p *CrittercismPlugin) Init(job *job.Job) error {
	var err error

	config := job.Config()

	p.api, err = NewCrittercismAPIClient(config["username"].(string), config["password"].(string), config["appId"].(string))

	if err != nil {
		return err
	}

	board := config["board"].(map[interface{}]interface{})
	boardName := board["name"].(string)
	boardPrefix := board["prefix"].(string)
	template := `{"name":"Crittercism","theme":"dark","aspect_ratio":"HDTV","font_family":"normal","font_size":"normal","widget_background":"","widget_margins":3,"widget_padding":8,"widgets":[{"flow":{"tag":"app_error_rates","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Service Error Rates"}},"variant":"barchart","column":24,"row":12,"width":8,"height":8,"in_board_index":0,"background":"default"},{"flow":{"tag":"crittercism_logo","data":{"mode":"fit","url":"https://s3.amazonaws.com/telemetrydemos/images/Logos/crittercism_crop.png"}},"variant":"image","column":0,"row":0,"width":32,"height":3,"in_board_index":1,"background":"none"},{"flow":{"tag":"crashes","data":{}},"variant":"box","column":0,"row":3,"width":8,"height":17,"in_board_index":2,"background":"default"},{"flow":{"tag":"users","data":{}},"variant":"box","column":8,"row":3,"width":16,"height":11,"in_board_index":3,"background":"default"},{"flow":{"tag":"app_store_ratings","data":{"icons":[{"color":"rgb(212, 212, 212)","label":"4.2","type":"fa-android"},{"color":"rgb(212, 212, 212)","label":"4.5","type":"fa-apple"}],"title":"App Store Ratings"}},"variant":"icon","column":8,"row":14,"width":16,"height":6,"in_board_index":4,"background":"default"},{"flow":{"tag":"crash_count","data":{"icon":"fa-bomb","label":"Total Crashes","value":"124"}},"variant":"value","column":0,"row":17,"width":8,"height":3,"in_board_index":105,"background":"none"},{"flow":{"tag":"crash_by_os","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}]}},"variant":"barchart","column":0,"row":7,"width":8,"height":10,"in_board_index":106,"background":"none"},{"flow":{"tag":"user_counts","data":{"values":[{"icon":"fa-download","label":"Loads","value":100},{"icon":"fa-group","label":"Daily","value":100},{"icon":"fa-group","label":"Monthly","value":100}]}},"variant":"multivalue","column":8,"row":3,"width":16,"height":4,"in_board_index":107,"background":"none"},{"flow":{"tag":"daily_active_users","data":{"renderer":"line","series":[{"values":[22,16,79,26,7,6,40,50,84,46]}],"title":"Daily Active Users"}},"variant":"graph","column":8,"row":7,"width":16,"height":7,"in_board_index":108,"background":"none"},{"flow":{"tag":"app_loads_by_device","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Loads By Device"}},"variant":"barchart","column":24,"row":3,"width":8,"height":9,"in_board_index":109,"background":"default"},{"flow":{"tag":"crash_rate","data":{"color":"rgb(246, 120, 0)","icon":"fa-bomb","label":"Crash Rate","value":"1.2","value_type":"percent"}},"variant":"value","column":0,"row":4,"width":8,"height":3,"in_board_index":110,"background":"none"}]}`

	b, err := job.GetOrCreateBoard(boardName, boardPrefix, template)

	if err != nil {
		return err
	}

	// Daily users

	const refreshInterval = 60

	err = p.registerClosures(
		b,
		[]crittercismPluginClosure{
			// Users and loads

			crittercismPluginClosure{p.DailyActiveUsers, refreshInterval, "daily_active_users"},
			crittercismPluginClosure{p.DailyMonthlyLoadsUsers, refreshInterval, "user_counts"},

			// Crashes

			crittercismPluginClosure{p.DailyAppCrashes, 1, "crash_count"},
			crittercismPluginClosure{p.DailyCrashRate, 1, "crash_rate"},

			// App loads by device

			crittercismPluginClosure{p.AppLoadsByDevice, refreshInterval, "app_loads_by_device"},

			// App service error rates

			crittercismPluginClosure{p.AppServiceErrorRates, refreshInterval, "app_error_rates"},

			// Crashes by OS

			crittercismPluginClosure{p.CrashesByOS, refreshInterval, "crash_by_os"},

			// Ratings

			crittercismPluginClosure{p.AppStoreRatings, refreshInterval, "app_store_ratings"},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
