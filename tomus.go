package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/newrelic"
	"github.com/paganotoni/tomus/request"
)

// Setup receives the app it will add the logger and other tools and from that
// It adds
func Setup(app *buffalo.App) {
	logColorsEnabled := envy.Get("GO_ENV", "development") == "development"
	app.Logger = logentries.NewLogger(logColorsEnabled)

	newrelic.MountTo(app, app.Logger)
	request.MountTo(app)
}
