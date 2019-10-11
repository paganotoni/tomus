package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/newrelic"
	"github.com/paganotoni/tomus/request"
)

var logColorsEnabled = envy.Get("GO_ENV", "development") == "development"
var Logger = logentries.NewLogger(logColorsEnabled)

// Setup receives the app it will add the logger and other tools and from that
// it adds NewRelic and Logentries elements into the buffalo app.
func Setup(app *buffalo.App, r *render.Engine) {
	app.Logger = Logger

	newrelic.MountTo(app, app.Logger)
	request.MountTo(app, r)
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
