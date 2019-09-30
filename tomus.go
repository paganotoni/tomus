package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/newrelic"
	"github.com/paganotoni/tomus/request"
)

func Setup(app *buffalo.App) {
	app.Logger = logentries.NewLogger()

	newrelic.MountTo(app)
	request.MountTo(app)
}
