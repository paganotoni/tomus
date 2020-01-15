package datadog

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/events"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type monitor struct {
	ServiceName string
}

func (dd *monitor) Listen() error {
	events.Listen(func(e events.Event) {
		switch e.Kind {
		case buffalo.EvtAppStart:
			dd.appStart(e)
		case buffalo.EvtAppStop:
			dd.appStop(e)
		case buffalo.EvtRouteStarted:
			dd.routeStarted(e)
		case buffalo.EvtRouteFinished:
			dd.routeFinished(e)
		}
	})

	return nil
}

func (dd *monitor) appStart(e events.Event) {
	tracer.Start(tracer.WithServiceName(dd.ServiceName))
}

func (dd *monitor) appStop(e events.Event) {
	defer tracer.Stop()
}

func (dd *monitor) routeStarted(e events.Event) {

	ro, err := e.Payload.Pluck("route")
	if err != nil {
		//Do something
		return
	}

	currentRoute := ro.(buffalo.RouteInfo)

	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		//Do something
		return
	}

	c := ctx.(buffalo.Context)

	resourceName := c.Request().Method + " " + currentRoute.Path

	opts := []ddtrace.StartSpanOption{
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.ServiceName(dd.ServiceName),
		tracer.ResourceName(resourceName),

		tracer.Tag(ext.EventSampleRate, math.NaN()),
		tracer.Tag("mux.host", ""),
		tracer.Tag(ext.HTTPMethod, c.Request().Method),
		tracer.Tag(ext.HTTPURL, c.Request().URL.Path),
	}

	if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(c.Request().Header)); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span := tracer.StartSpan("http.request", opts...)
	c.Set("span", span)
}

func (dd *monitor) routeFinished(e events.Event) {
	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		return
	}

	c := ctx.(buffalo.Context)

	var response *buffalo.Response
	var span tracer.Span
	var ok bool

	if span, ok = c.Value("span").(tracer.Span); !ok {
		return
	}

	if response, ok = c.Response().(*buffalo.Response); !ok {
		return
	}

	status := response.Status
	span.SetTag(ext.HTTPCode, strconv.Itoa(status))
	if status >= 500 && status < 600 {
		span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(status)))
	}

	c.Logger().Debug("Finishing Span")
	span.Finish()
}

// Track ...
func (dd *monitor) Track(name string, fn func() error) error {
	return errors.New("Needs to be implemented")
}

// NewMonitor creates a new monitor for DataDog with the passed serviceName
func NewMonitor(serviceName string) *monitor {
	return &monitor{
		ServiceName: serviceName,
	}
}
