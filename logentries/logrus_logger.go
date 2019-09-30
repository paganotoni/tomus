package logentries

import (
	"io"

	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

// Logrus is a Logger implementation backed by sirupsen/logrus
type LogrusLogger struct {
	logrus.FieldLogger
}

// SetOutput will try and set the output of the underlying
// logrus.FieldLogger if it can
func (l LogrusLogger) SetOutput(w io.Writer) {
	if lg, ok := l.FieldLogger.(logger.Outable); ok {
		lg.SetOutput(w)
	}
}

// WithField returns a new Logger with the field added
func (l LogrusLogger) WithField(s string, i interface{}) logger.FieldLogger {
	return LogrusLogger{l.FieldLogger.WithField(s, i)}
}

// WithFields returns a new Logger with the fields added
func (l LogrusLogger) WithFields(m map[string]interface{}) logger.FieldLogger {
	return LogrusLogger{l.FieldLogger.WithFields(m)}
}
