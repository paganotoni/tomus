package tomus

import (
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/paganotoni/tomus/datadog"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/newrelic"
	"github.com/paganotoni/tomus/request"
)

var (
	//Logger ...
	Logger = logentries.NewLogger(
		envy.Get("GO_ENV", "development") == "development",
	)
)

// Setup receives the Config and from that
// it adds NewRelic and Logentries elements into the buffalo app.
func Setup(config Config) error {
	app := config.App

	if app == nil {
		return errors.New("app cannot be nil")
	}

	switch config.APMKind {
	case APMKindNewrelic:
		newrelic.MountTo(app, app.Logger)
	case APMKindDatadog:
		dd := datadog.Monitor{
			ServiceName: config.ServiceName,
			Enabled:     config.EnableAPM,
			Environment: config.Environment,
		}

		dd.Monitor()
	}

	app.Logger = Logger
	request.MountTo(app, config.RenderEngine)

	return nil
}

//TrackError allows to track errors that are not exactly inside a New Relic Tx.
func TrackError(c buffalo.Context, err error) error {
	return newrelic.NewTracker().TrackError(c, err)
}

// TrackBackgroundTransaction allows to track operations that happen in background
// by wrapping those operations within a New Relic tx.
func TrackBackgroundTransaction(name string, fn func() error) {
	newrelic.NewTracker().TrackBackgroundTransaction(name, fn)
}
