package main

import (
	"github.com/telemetryapp/gotelemetry"
	"github.com/telemetryapp/gotelemetry_agent/agent"
	"github.com/telemetryapp/gotelemetry_agent/agent/aggregations"
	"github.com/telemetryapp/gotelemetry_agent/agent/config"
	"github.com/telemetryapp/gotelemetry_agent/agent/functions"
	"github.com/telemetryapp/gotelemetry_agent/agent/graphite"
	"github.com/telemetryapp/gotelemetry_agent/agent/job"
	"github.com/telemetryapp/gotelemetry_agent/agent/server"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/telemetryapp/agent_crittercism/crittercism"
	_ "github.com/telemetryapp/gotelemetry_agent/plugin"
)

var configFile *config.ConfigFile
var errorChannel chan error
var completionChannel chan bool

func main() {
	if config.CLIConfig.WantsFunctionHelp {
		functions.PrintHelp(config.CLIConfig.FunctionHelpName)
		return
	}

	var err error

	configFile, err = config.NewConfigFile()

	if err != nil {
		log.Fatalf("Initialization error: %s", err)
	}

	errorChannel = make(chan error, 0)
	completionChannel = make(chan bool, 0)

	go run()

	for {
		select {
		case err := <-errorChannel:
			if e, ok := err.(*gotelemetry.Error); ok {
				logLevel := e.GetLogLevel()

				if logLevel >= config.CLIConfig.LogLevel {
					prefix := "Error"

					switch logLevel {
					case gotelemetry.LogLevelLog:
						prefix = "Info "

					case gotelemetry.LogLevelDebug:
						prefix = "Debug"
					}

					log.Printf("%s: %s", prefix, err)
				}

				continue
			}

			log.Printf("Error: %s", err.Error())

		case <-completionChannel:
			goto Done
		}
	}

Done:

	log.Println("No more jobs to run; exiting.\n")
}

func run() {
	if err := aggregations.Init(configFile.DataConfig().DataLocation, configFile.DataConfig().TTL, errorChannel); err != nil {
		log.Fatalf("Initialization error: %s", err)
	}

	if err := graphite.Init(configFile, errorChannel); err != nil {
		log.Fatalf("Initialization error: %s", err)
	}

	if config.CLIConfig.IsPiping {
		payload, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			errorChannel <- err
		}

		agent.ProcessPipeRequest(configFile, errorChannel, completionChannel, payload)
	} else if config.CLIConfig.IsNotifying {
		agent.ProcessNotificationRequest(configFile, errorChannel, completionChannel, config.CLIConfig.NotificationChannel, config.CLIConfig.NotificationFlow, config.CLIConfig.Notification)
	} else {
		if configFile.ListenAddress() != "" {
			if err := server.Init(configFile, errorChannel); err != nil {
				log.Fatal("Web server initialization error: %s", err)
			}
		}

		_, err := job.NewJobManager(configFile, errorChannel, completionChannel)

		if err != nil {
			log.Fatalf("Initialization error: %s", err)
		}
	}
}
