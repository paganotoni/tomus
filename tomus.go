package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/newrelic"
	"github.com/paganotoni/tomus/request"
)

// Setup receives the app it will add the logger and other tools and from that
// it adds NewRelic and Logentries elements into the buffalo app.
func Setup(app *buffalo.App) {
	logColorsEnabled := envy.Get("GO_ENV", "development") == "development"
	app.Logger = logentries.NewLogger(logColorsEnabled)

	newrelic.MountTo(app, app.Logger)
	request.MountTo(app)
}

func TrackError(c buffalo.Context, err error) {
	newrelic.NewTracker().TrackError(c, err)
}

func TrackBackgroundTransaction(name string, fn func() error) {
	newrelic.NewTracker().TrackBackgroundTransaction(name, fn)
}
