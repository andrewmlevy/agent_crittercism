package crittercism

import (
	"github.com/telemetryapp/agent_critercism/flows"
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"math/rand"
	"time"
)

func init() {
	job.RegisterPlugin("Crittercism", CrittercismPluginFactory)
}

func CrittercismPluginFactory() job.PluginInstance {
	return &CrittercismPlugin{
		job.NewPluginHelper(),
	}
}

type CrittercismPlugin struct {
	*job.PluginHelper
}

func (r *CrittercismPlugin) Init(job *job.Job) error {

	config := job.Config()

	board := config["board"].(map[interface{}]interface{})
	boardName := board["name"].(string)
	boardPrefix := board["prefix"].(string)
	template := `{"name":"Crittercism","theme":"dark","aspect_ratio":"HDTV","font_family":"normal","font_size":"normal","widget_background":"","widget_margins":3,"widget_padding":8,"widgets":[{"variant":"image","tag":"image_1","column":0,"row":0,"width":9,"height":6,"in_board_index":0,"background":"none"},{"variant":"image","tag":"image_2","column":25,"row":1,"width":5,"height":4,"in_board_index":1,"background":"none"},{"variant":"box","tag":"box_1","column":0,"row":6,"width":32,"height":8,"in_board_index":2,"background":"default"},{"variant":"multivalue","tag":"multivalue_1","column":0,"row":14,"width":11,"height":6,"in_board_index":3,"background":"none"},{"variant":"graph","tag":"graph_1","column":11,"row":14,"width":11,"height":6,"in_board_index":4,"background":"none"},{"variant":"piechart","tag":"piechart_1","column":22,"row":14,"width":10,"height":6,"in_board_index":5,"background":"none"},{"variant":"value","tag":"value_2","column":0,"row":8,"width":8,"height":5,"in_board_index":106,"background":"none"},{"variant":"value","tag":"value_3","column":8,"row":8,"width":8,"height":5,"in_board_index":107,"background":"none"},{"variant":"value","tag":"value_4","column":24,"row":8,"width":8,"height":5,"in_board_index":108,"background":"none"},{"variant":"value","tag":"value_1","column":11,"row":0,"width":10,"height":6,"in_board_index":109,"background":"none"},{"variant":"text","tag":"text_1","column":8,"row":6,"width":17,"height":3,"in_board_index":110,"background":"none"},{"variant":"value","tag":"value_5","column":16,"row":8,"width":8,"height":5,"in_board_index":111,"background":"none"}]}`
	b, err := job.GetOrCreateBoard(boardName, boardPrefix, template)
	if err != nil {
		return err
	}

	if err = r.PluginHelper.AddTaskWithClosureFromBoardForFlowWithTag(flows.DailyActiveUsers, time.Second*3600, b, "daily_active_users"); err != nil {
		return err
	}

	return nil
}
