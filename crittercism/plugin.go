package crittercism

import (
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
	template := `{"name":"Default","theme":"dark","aspect_ratio":"HDTV","font_family":"normal","font_size":"normal","widget_background":"","widget_margins":3,"widget_padding":8,"widgets":[{"flow":{"tag":"box_1","data":{}},"variant":"box","column":8,"row":4,"width":24,"height":4,"in_board_index":0,"background":"default"},{"flow":{"tag":"box_2","data":{}},"variant":"box","column":8,"row":8,"width":24,"height":4,"in_board_index":1,"background":"default"},{"flow":{"tag":"box_3","data":{}},"variant":"box","column":0,"row":4,"width":8,"height":16,"in_board_index":2,"background":"default"},{"flow":{"tag":"image_1","data":{"mode":"fit","url":"https://s3.amazonaws.com/telemetrydemos/images/Logos/crittercism_crop.png"}},"variant":"image","column":0,"row":0,"width":32,"height":4,"in_board_index":3,"background":"none"},{"flow":{"tag":"app_store_ratings","data":{"cells":[[{"icon":"fa-android"},{"alignment":"left","value":4.5}],[{"icon":"fa-apple"},{"alignment":"left","value":4.3}]],"title":"App Store Ratings"}},"variant":"table","column":16,"row":12,"width":8,"height":8,"in_board_index":4,"background":"default"},{"flow":{"tag":"service_monitoring_error_rate","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Service Error Rates"}},"variant":"barchart","column":8,"row":12,"width":8,"height":8,"in_board_index":5,"background":"default"},{"flow":{"tag":"daily_app_crashes","data":{"icon":"fa-bug","label":"Crashes","value":124}},"variant":"value","column":0,"row":4,"width":8,"height":4,"in_board_index":106,"background":"none"},{"flow":{"tag":"barchart_1","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}]}},"variant":"barchart","column":0,"row":12,"width":8,"height":8,"in_board_index":107,"background":"none"},{"flow":{"tag":"daily_crash_rate","data":{"color":"rgb(246, 120, 0)","label":"Crash Percent","value":1.2,"value_type":"percent"}},"variant":"value","column":0,"row":8,"width":8,"height":4,"in_board_index":108,"background":"none"},{"flow":{"tag":"monthly_active_users","data":{"renderer":"line","series":[{"values":[35,62,85,50,93,47,97,37,98,83]}],"title":"Monthly Active Users"}},"variant":"graph","column":16,"row":4,"width":16,"height":4,"in_board_index":109,"background":"none"},{"flow":{"tag":"daily_app_loads","data":{"label":"App Loads","value":23000}},"variant":"value","column":8,"row":8,"width":4,"height":4,"in_board_index":110,"background":"none"},{"flow":{"tag":"monthly_users","data":{"label":"Monthly Users","value":18900}},"variant":"value","column":8,"row":4,"width":8,"height":4,"in_board_index":111,"background":"none"},{"flow":{"tag":"daily_users","data":{"label":"Daily Users","value":4545}},"variant":"value","column":12,"row":8,"width":4,"height":4,"in_board_index":112,"background":"none"},{"flow":{"tag":"daily_active_users","data":{"renderer":"line","series":[{"values":[22,16,79,26,7,6,40,50,84,46]}],"title":"Daily Active Users"}},"variant":"graph","column":16,"row":8,"width":16,"height":4,"in_board_index":113,"background":"none"},{"flow":{"tag":"daily_app_loads_by_device","data":{"bars":[{"color":"#267288","label":"New Bar","value":100}],"title":"App Loads By Device"}},"variant":"barchart","column":24,"row":12,"width":8,"height":8,"in_board_index":114,"background":"default"}]}`

	b, err := job.GetOrCreateBoard(boardName, boardPrefix, template)

	if err != nil {
		return err
	}

	if err = p.PluginHelper.AddTaskWithClosureFromBoardForFlowWithTag(p.DailyActiveUsers, time.Second*1, b, "daily_active_users"); err != nil {
		return err
	}

	return nil
}
