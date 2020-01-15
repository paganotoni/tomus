package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/paganotoni/tomus/logentries"
	"github.com/paganotoni/tomus/request"
)

// New receives the Config and from that
// creates a new tomusWrapper
func New(config Config) Wrapper {
	return Wrapper{
		config: config,

		Logger: logentries.NewLogger(
			envy.Get("GO_ENV", "development") == "development",
		),
		Monitor: config.APMMonitor,
	}
}

type Wrapper struct {
	config Config

	Logger  buffalo.Logger
	Monitor APMMonitor
}

// Start ...
func (h Wrapper) Start() error {
	app := h.config.App

	app.Logger = h.Logger
	request.MountTo(app)

	if h.config.APMMonitor == nil {
		return nil
	}

	return h.config.APMMonitor.Listen()
}
