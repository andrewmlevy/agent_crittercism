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
		"",
		map[string]*gotelemetry.Flow{},
		"",
	}
}

type CrittercismPlugin struct {
	*job.PluginHelper
	api       *CrittercismAPIClient
	appName   string
	flows     map[string]*gotelemetry.Flow
	ratingKey string
}

type crittercismPluginClosure struct {
	closure  job.PluginHelperClosureWithFlow
	interval int
	tag      string
}

func (p *CrittercismPlugin) registerClosures(b *gotelemetry.Board, closures []crittercismPluginClosure) error {
	for _, closure := range closures {
		if err := p.PluginHelper.AddTaskWithClosureForFlowWithTag(closure.closure, time.Second*time.Duration(closure.interval), p.flows, closure.tag); err != nil {
			return err
		}
	}

	return nil
}

func (p *CrittercismPlugin) Init(job *job.Job) error {
	var err error

	config := job.Config()

	p.api, err = NewCrittercismAPIClient(config["apiKey"].(string), config["appId"].(string))

	if err != nil {
		return err
	}

	board := config["board"].(map[interface{}]interface{})
	boardName := board["name"].(string)
	boardPrefix := board["prefix"].(string)
	template := `{"name":"Crittercism","theme":"dark","aspect_ratio":"HDTV","font_family":"normal","font_size":"normal","widget_background":"","widget_margins":3,"widget_padding":8,"widgets":[{"flow":{"tag":"crittercism_logo","data":{"mode":"fit","url":"https://s3.amazonaws.com/telemetrydemos/images/Logos/crittercism_crop.png"}},"variant":"image","column":0,"row":0,"width":11,"height":3,"in_board_index":0,"background":"none"},{"flow":{"tag":"users","data":{}},"variant":"box","column":0,"row":3,"width":32,"height":5,"in_board_index":1,"background":"default"},{"flow":{"tag":"app_name","data":{"alignment":"right","text":"Application Name"}},"variant":"text","column":14,"row":0,"width":18,"height":3,"in_board_index":2,"background":"none"},{"flow":{"tag":"app_error_rates","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Service Error Rates"}},"variant":"barchart","column":24,"row":8,"width":8,"height":6,"in_board_index":3,"background":"default"},{"flow":{"tag":"most_frequent_crashes","data":{"cells":[[{"value":"Row1 Col1"},{"value":"4"}]],"title":"Most Frequent Crashes"}},"variant":"table","column":8,"row":13,"width":16,"height":7,"in_board_index":4,"background":"default"},{"flow":{"tag":"daily_users","data":{"icon":"fa-group","label":"Daily Users","value":100}},"variant":"value","column":16,"row":5,"width":5,"height":3,"in_board_index":105,"background":"none"},{"flow":{"tag":"crash_count","data":{"icon":"fa-bug","label":"Total Crashes","value":"124"}},"variant":"value","column":5,"row":5,"width":5,"height":3,"in_board_index":106,"background":"none"},{"flow":{"tag":"crash_by_os","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"Crash By OS"}},"variant":"barchart","column":0,"row":8,"width":8,"height":12,"in_board_index":107,"background":"default"},{"flow":{"tag":"monthly_users","data":{"icon":"fa-group","label":"Monthly Users","value":100}},"variant":"value","column":21,"row":5,"width":6,"height":3,"in_board_index":108,"background":"none"},{"flow":{"tag":"daily_active_users","data":{"renderer":"line","series":[{"values":[22,16,79,26,7,6,40,50,84,46]}],"title":"Daily Active Users"}},"variant":"graph","column":8,"row":8,"width":16,"height":5,"in_board_index":109,"background":"default"},{"flow":{"tag":"crash_rate","data":{"color":"rgb(246, 120, 0)","icon":"fa-bug","label":"Crash Rate","value":"1.2","value_type":"percent"}},"variant":"value","column":0,"row":5,"width":5,"height":3,"in_board_index":110,"background":"none"},{"flow":{"tag":"app_loads","data":{"icon":"fa-download","label":"App Loads","value":100}},"variant":"value","column":10,"row":5,"width":6,"height":3,"in_board_index":111,"background":"none"},{"flow":{"tag":"app_loads_by_device","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Loads By Device"}},"variant":"barchart","column":24,"row":14,"width":8,"height":6,"in_board_index":112,"background":"default"},{"flow":{"tag":"app_rating","data":{"icon":"fa-android","label":"Rating","rounding":1,"value":4.3}},"variant":"value","column":27,"row":5,"width":5,"height":3,"in_board_index":113,"background":"none"},{"flow":{"tag":"date","data":{"alignment":"center","text":"Friday, October 31st"}},"variant":"text","column":10,"row":3,"width":12,"height":3,"in_board_index":114,"background":"none"}]}`

	p.appName = config["appName"].(string)
	p.ratingKey = config["ratingKey"].(string)

	b, err := job.GetOrCreateBoard(boardName, boardPrefix, template)

	if err != nil {
		return err
	}

	job.Logf("Created board %s", b.Id)

	// Daily users

	refreshInterval, ok := config["refresh"].(int)

	if !ok {
		refreshInterval = 60
	}

	p.flows, err = b.MapWidgetsToFlows()

	if err != nil {
		return err
	}

	err = p.registerClosures(
		b,
		[]crittercismPluginClosure{

			crittercismPluginClosure{p.SetAppName, 500000, "app_name"},
			crittercismPluginClosure{p.SetDate, 60, "date"},

			// Users and loads

			crittercismPluginClosure{p.DailyActiveUsersGraph, refreshInterval, "daily_active_users"},

			crittercismPluginClosure{p.DailyActiveUsersValue, refreshInterval, "daily_users"},
			crittercismPluginClosure{p.MonthlyActiveUsersValue, refreshInterval, "monthly_users"},
			crittercismPluginClosure{p.DailyAppLoadsValue, refreshInterval, "app_loads"},

			// Crashes

			crittercismPluginClosure{p.DailyAppCrashes, refreshInterval, "crash_count"},
			crittercismPluginClosure{p.DailyCrashRate, refreshInterval, "crash_rate"},
			crittercismPluginClosure{p.DailyMostFrequentCrashes, refreshInterval, "most_frequent_crashes"},

			// App loads by device

			crittercismPluginClosure{p.AppLoadsByDevice, refreshInterval, "app_loads_by_device"},

			// App service error rates

			crittercismPluginClosure{p.AppServiceErrorRates, refreshInterval, "app_error_rates"},

			// Crashes by OS

			crittercismPluginClosure{p.CrashesByOS, refreshInterval, "crash_by_os"},

			// Ratings

			crittercismPluginClosure{p.AppStoreRatings, refreshInterval, "app_rating"},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
