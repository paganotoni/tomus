package newrelic

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"

	newrelic "github.com/newrelic/go-agent"
)

var (
	config newrelic.Config
	nrApp  newrelic.Application
	logger buffalo.Logger

	license = envy.Get("NEWRELIC_LICENSE_KEY", "")
	name    = envy.Get("NEWRELIC_APP_NAME", "app/byod")
	env     = envy.Get("NEWRELIC_ENV", "staging")
)

func MountTo(app *buffalo.App, logger buffalo.Logger) {
	app.Use(middleware)
}

func init() {
	suffix := envy.Get("NEWRELIC_SUFFIX", fmt.Sprintf("(%v)", env))

	config = newrelic.NewConfig(fmt.Sprintf("%v %v", name, suffix), license)
	config.Enabled, _ = strconv.ParseBool(os.Getenv("ENABLE_NEWRELIC"))
	config.DistributedTracer.Enabled = true

	config.Labels = map[string]string{
		"ENVIRONMENT": env,
	}

	var err error
	if nrApp, err = newrelic.NewApplication(config); err != nil {
		logger.Error(errors.Wrap(err, "tomus error creating newrelic app"))
	}
}

func middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if !config.Enabled {
			c.Logger().Info("Skipping newrelic middleware because its disabled")
			return next(c)
		}

		txn := nrApp.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())

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
func TrackBackgroundTransaction(name string, fn func() error) {
	if nrApp == nil {
		if err := fn(); err != nil {
			logger.Error(err)
		}

		return
	}

	txn := nrApp.StartTransaction(name, nil, nil)
	if err := fn(); err != nil {
		if nrErr := txn.NoticeError(err); nrErr != nil {
			logger.Error(nrErr)
		}
	}

	if err := txn.End(); err != nil {
		logger.Error(err)
	}
}

//TrackError tracks a newrelic error this is useful for error pages and crashes.
func TrackError(c buffalo.Context, err error) error {
	var oerr error

	if !config.Enabled {
		c.Logger().Info("Not notifying error to NR since is disabled")
		return nil
	}

	txn := nrApp.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())
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
