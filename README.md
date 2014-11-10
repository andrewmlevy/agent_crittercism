# Telemetry Crittercism Agent

This app simplifies the creation of standardized [Telemetry](https://telemetryapp.com) boards that express data from [Crittercism](http://www.crittercism.com).

The app is written in Go and depends on several external packages, including:

- [Gotelemetry](https://github.com/telemetryapp/gotelemetry) as the underlying interface to Telemetry
- [Gotelemetry_agent](https://github.com/telemetryapp/gotelemetry_agent) as the base agent app

## Capabilities

The Agent is designed to create and update an arbitrary number of standardized boards. The data is pulled directly from the Crittercism API and used to populate the appropriate Telemetry flows.

There are no limits to the number of boards that can be created or maintained, although API limits may come into play based on the number of concurrent operations performed for a given Crittercism or Telemetry account.

## Installation

Binaries for Windows, OS X, and Linux are included in the repository under the `bin` directory. They are completely self-contained and have no external dependencies.

The Agent takes one parameter, `-config`, which points to the location of a configuration file that describes the operations that should be performed. For example:

```
./bin/linux-amd64/crittercism_telemetry_agent.go -config /path/to/config.yaml
```

## Configuration

The configuration file is written in [YAML](http://www.yaml.org) and is organized as an array of Telemetry accounts, each containing one or more jobs, each of which identifies a board to be updated:

```yaml
accounts: 
  - api_key: <your telemetry api key here>
    submission_interval: 1
    jobs:
      - id: <job id>
        plugin: crittercism
        config:
        	refresh: 60
          apiKey: <crittercism api key>
          appId: <crittercism app id>
          appName: <app name>
          ratingKey: <rating key>
          board:
            name: <unique board name>
            prefix: <unique board prefix>
```

The top-level `accounts` array contains one entry for each Telemetry account which must receive data. In turn, each account contains a `jobs` array that describes an arbitrary number of boards that should be created and updated for that particular account. Each account also contains a `submission_interval` entry, which indicates how often data updates are coealesced and sent to the Telemetry API; it's safe to leave this value to a default of `1` second.

Jobs possess the following properties:

- `id` (string) — must be unique across the entire configuration file. This is used primarily for logging purposes, and can be any string that is convenient to the user.
- `plugin` (string) — must be set to `crittercism` to indicate that the agent should talk the Crittercism plugin with running the job.
- `config.refresh` (int, defaults to `60`) — indicates the data refresh interval in seconds. It should be set in a way so as to strike the correct balance between keeping the data fresh and avoiding to overwhelm either the Telemetry or Crittercism APIs.
- `config.apiKey` (string) — the Crittercism oAuth token to be used when making requests for this board.
- `config.appId` (string) — the Crittercism ID of the app whose data is to be polled for this job.
- `config.appName` (string) — the name of the app whose data is to be polled. Used to populat the title in the resulting board.
- `config.ratingKey` (string, one of `ios`, `android`, `wp8`, `html5`) — identifies the platform whose rating should be shown on the board.
- `config.board.name` (string) — a unique name for the board that the job will output. This name must be unique across the entire Telemetry account, and will only be displayed in the Telemetry editor.
- `config.board.prefix` (string, must contain only letters, numbers, dashes, and underscores) — A prefix that is associated with every flow in the board created by the job. This prefix is used to ensure that the contents of a board remain unique across the entire Telemetry account—a good strategy is to use (for example) the App ID, or some other similarly unique value.

All the parameters listed above, with the exception of `config.refresh` are required.



