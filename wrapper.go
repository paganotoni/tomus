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
	logger := logentries.NewLogger(
		envy.Get("GO_ENV", "development") == "development",
	)

	return Wrapper{
		config: config,
		logger: logger,

		Monitor: config.APMMonitor,
	}
}

type Wrapper struct {
	config Config
	logger buffalo.Logger

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

//Logger returns the logger used
func (h Wrapper) Logger() buffalo.Logger {
	if h.logger == nil {
		h.logger = logentries.NewLogger(envy.Get("GO_ENV", "development") == "development")
	}

	return h.logger
}
