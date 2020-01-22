package newrelic

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/events"
	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type monitor struct {
	newrelicApplication newrelic.Application
}

func (nr *monitor) Listen() error {
	events.Listen(func(e events.Event) {
		switch e.Kind {
		case buffalo.EvtRouteStarted:
			nr.routeStarted(e)
		case buffalo.EvtRouteErr:
			nr.routeError(e)
		case buffalo.EvtRouteFinished:
			nr.routeFinished(e)
		}
	})

	return nil
}

func (nr *monitor) routeStarted(e events.Event) {
	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		return
	}

	c := ctx.(buffalo.Context)

	txn := nr.newrelicApplication.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())

	ri := c.Value("current_route").(buffalo.RouteInfo)
	if err := txn.AddAttribute("PathName", ri.PathName); err != nil {
		c.Logger().Error(err)
	}

	if err := txn.AddAttribute("RequestID", c.Value("request_id")); err != nil {
		c.Logger().Error(err)
	}

	c.Set("newrelicTransaction", txn)
}

func (nr *monitor) routeError(e events.Event) {
	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		return
	}

	c := ctx.(buffalo.Context)
	nrTxn := c.Value("newrelicTransaction")
	txn, ok := nrTxn.(newrelic.Transaction)
	if !ok {
		return
	}

	txn.NoticeError(err)
}

func (nr *monitor) routeFinished(e events.Event) {

	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		return
	}

	c := ctx.(buffalo.Context)

	nrTxn := c.Value("newrelicTransaction")
	txn, ok := nrTxn.(newrelic.Transaction)
	if !ok {
		return
	}

	if err := txn.End(); err != nil {
		c.Logger().Error(err)
	}
}

// Track ...
func (nr *monitor) Track(name string, fn func() error) error {
	return errors.New("Needs to be implemented")
}

// NewMonitor creates a new monitor for DataDog with the passed serviceName
func NewMonitor(serviceName, env, license string) *monitor, error {
	nrName := strings.Replace(serviceName, "/", "", 1)
	nrName = fmt.Sprintf("%v (%v)", nrName, env)

	config := newrelic.NewConfig(nrName, license)
	config.Enabled = true
	config.DistributedTracer.Enabled = true
	config.Labels = map[string]string{
		"ENVIRONMENT": env,
	}

	var app newrelic.Application
	var err error
	if app, err = newrelic.NewApplication(config); err != nil {
		logrus.Error(errors.Wrap(err, "tomus error creating newrelic app"))
	}

	return &monitor{
		newrelicApplication: app,
	}
}
