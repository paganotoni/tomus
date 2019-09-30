package logentries

import (
	"github.com/bsphere/le_go"
	"github.com/markbates/safe"
	"github.com/sirupsen/logrus"
)

type LogentriesHook struct {
	token string
	conn  *le_go.Logger
}

// NewLogentriesHook creates a hook to be added to an instance of logger.
func NewLogentriesHook(token string) *LogentriesHook {
	le, err := le_go.Connect(token)
	if err != nil {
		logrus.Error(err)
	}

	return &LogentriesHook{
		token: token,
		conn:  le,
	}
}

// Fire is called when a log event is fired.
func (hook *LogentriesHook) Fire(entry *logrus.Entry) error {
	go func() {
		err := safe.Run(func() {
			if hook.conn == nil {
				return
			}

			msg, err := entry.String()
			if err != nil {
				return
			}

			hook.conn.Println(msg)
		})

		if err != nil {
			logrus.Error(err)
		}
	}()

	return nil
}

// Levels returns the available logging levels.
func (hook *LogentriesHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
