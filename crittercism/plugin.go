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
		map[string]interface{}{},
	}
}

type CrittercismPlugin struct {
	*job.PluginHelper,
  config map[string]interface{}
}

func (r *CrittercismPlugin) Init(job *job.Job, config map[string]interface{}) error {

  r.config = config

	board := config["board"].(map[interface{}]interface{})
	boardName := board["name"].(string)
	boardPrefix := board["prefix"].(string)
	template := "{\"name\":\"CrittercismTest\",\"theme\":\"dark\",\"aspect_ratio\":\"HDTV\",\"font_family\":\"normal\",\"font_size\":\"normal\",\"widget_background\":\"\",\"widget_margins\":3,\"widget_padding\":8,\"widgets\":[{\"variant\":\"value\",\"tag\":\"value_98\",\"column\":7,\"row\":7,\"width\":8,\"height\":5,\"in_board_index\":0,\"background\":\"default\"},{\"variant\":\"value\",\"tag\":\"value_99\",\"column\":17,\"row\":7,\"width\":8,\"height\":5,\"in_board_index\":1,\"background\":\"default\"}]}"

	b, err := job.GetOrCreateBoard(boardName, boardPrefix, template)
	if err != nil {
		return err
	}

	if err = r.PluginHelper.AddTaskWithClosureFromBoardForFlowWithTag(flows.DailyActiveUsers, time.Second*3600, b, "daily_active_users"); err != nil {
		return err
	}

	return nil
}
