package newrelic

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	newrelic "github.com/newrelic/go-agent"
)

var (
	license = envy.Get("NEWRELIC_LICENSE_KEY", "")
	env     = envy.Get("NEWRELIC_ENV", "staging")
)

func MountTo(app *buffalo.App, logger buffalo.Logger) {
	trk := NewTracker()
	trk.logger = logger

	app.Use(trk.Middleware)
}

type tracker struct {
	config newrelic.Config
	app    newrelic.Application
	logger buffalo.Logger
}

func NewTracker() tracker {
	result := tracker{}

	name := envy.Get("NEWRELIC_APP_NAME", "app/byod")
	suffix := envy.Get("NEWRELIC_SUFFIX", fmt.Sprintf("(%v)", env))

	result.config = newrelic.NewConfig(fmt.Sprintf("%v %v", name, suffix), license)
	result.config.Enabled, _ = strconv.ParseBool(os.Getenv("ENABLE_NEWRELIC"))
	result.config.DistributedTracer.Enabled = true
	result.config.Labels = map[string]string{
		"ENVIRONMENT": env,
	}

	var err error
	if result.app, err = newrelic.NewApplication(result.config); err != nil {
		logrus.Error(errors.Wrap(err, "tomus error creating newrelic app"))
	}

	return result
}

func (t tracker) Middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if !t.config.Enabled {
			c.Logger().Info("Skipping newrelic middleware because its disabled")
			return next(c)
		}

		txn := t.app.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())

		ri := c.Value("current_route").(buffalo.RouteInfo)
		if err := txn.AddAttribute("PathName", ri.PathName); err != nil {
			c.Logger().Error(err)
		}

		if err := txn.AddAttribute("RequestID", c.Value("request_id")); err != nil {
			c.Logger().Error(err)
		}

		defer func() {
			if err := txn.End(); err != nil {
				c.Logger().Error(err)
			}
		}()

		err := next(c)
		if err != nil {
			c.Logger().Error(err)

			if terr := txn.NoticeError(err); terr != nil {
				c.Logger().Error(err)
			}
		}

		return err
	}
}

//TrackBackgroundTransaction allows to track non web transaction functions.
func (t tracker) TrackBackgroundTransaction(name string, fn func() error) {
	if t.app == nil {
		if err := fn(); err != nil {
			t.logger.Error(err)
		}

		return
	}

	txn := t.app.StartTransaction(name, nil, nil)
	if err := fn(); err != nil {
		if nrErr := txn.NoticeError(err); nrErr != nil {
			t.logger.Error(nrErr)
		}
	}

	if err := txn.End(); err != nil {
		t.logger.Error(err)
	}
}

//TrackError tracks a newrelic error this is useful for error pages and crashes.
func (t tracker) TrackError(c buffalo.Context, err error) error {
	var oerr error

	if !t.config.Enabled {
		c.Logger().Info("Not notifying error to NR since is disabled")
		return nil
	}

	txn := t.app.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())
	defer func() {
		if oerr = txn.End(); oerr != nil {
			c.Logger().Error(oerr)
		}
	}()

	if oerr = txn.AddAttribute("RequestID", c.Value("request_id")); oerr != nil {
		c.Logger().Error(oerr)
	}

	ri := c.Value("current_route").(buffalo.RouteInfo)
	if oerr = txn.AddAttribute("PathName", ri.PathName); oerr != nil {
		c.Logger().Error(oerr)
	}

	if oerr = txn.NoticeError(err); oerr != nil {
		return errors.Wrap(oerr, "MonitoringMW error notifying failed")
	}

	return nil
}
