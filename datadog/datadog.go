package datadog

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/events"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// DatadogMonitor ...
type DatadogMonitor struct {
	ServiceName string
	Host        string
	TracingPort string
	Enabled     bool
}

// Monitor ...
func (dd DatadogMonitor) Monitor() {
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
}

func (dd DatadogMonitor) appStart(e events.Event) {
	addr := net.JoinHostPort(
		dd.Host,
		dd.TracingPort,
	)

	tracer.Start(tracer.WithAgentAddr(addr))
}

func (dd DatadogMonitor) appStop(e events.Event) {
	defer tracer.Stop()
}

func (dd DatadogMonitor) routeStarted(e events.Event) {

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
		tracer.Tag("mux.host", dd.Host),
		tracer.Tag(ext.HTTPMethod, c.Request().Method),
		tracer.Tag(ext.HTTPURL, c.Request().URL.Path),
	}

	if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(c.Request().Header)); err == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span := tracer.StartSpan("http.request", opts...)
	c.Set("span", span)
}

func (dd DatadogMonitor) routeFinished(e events.Event) {
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

	c.Logger().Info(status)
	span.Finish()
}

// DatadogTracer ...
func DatadogTracer(next buffalo.Handler) buffalo.Handler {

	return func(c buffalo.Context) error {
		return next(c)
	}
}
