package logentries

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

var token = envy.Get("LOGENTRIES_TAG", "")

//NewLogger returns an instance of our logentries logger.
func NewLogger(enableColors bool) buffalo.Logger {
	if token == "" {
		return logger.NewLogger("debug")
	}

	return newLogentriesLogger(logrus.DebugLevel, token, enableColors)
}

// New based on the specified log level, defaults to "debug".
// This logger will log to the STDOUT in a human readable,
// but parseable form.
func newLogentriesLogger(lvl logger.Level, token string, enableColors bool) logger.FieldLogger {
	l := logrus.New()
	l.AddHook(NewLogentriesHook(token))
	l.SetOutput(os.Stdout)

	l.Level = lvl
	l.Formatter = &textFormatter{
		ForceColors: enableColors,
	}

	return LogrusLogger{l}
}
