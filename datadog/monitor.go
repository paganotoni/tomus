package datadog

import (
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/events"
	"github.com/paganotoni/tomus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Ensuring monitor is tomus.APMMonitor
var _ tomus.APMMonitor = (*monitor)(nil)

type monitor struct {
	//ServiceName is the identifier of the service used when reporting to DD APM.
	ServiceName string

	// Host is the agent host where monitor reports traces, defaults to localhost
	Host string

	// Port is the port used to report traces on the agent host, it defaults to 8126
	Port int
}

func (dd *monitor) Listen() error {
	_, err := events.Listen(func(e events.Event) {
		switch e.Kind {
		case buffalo.EvtAppStart:
			dd.appStart(e)
		case buffalo.EvtAppStop:
			dd.appStop(e)
		case buffalo.EvtRouteStarted:
			dd.routeStarted(e)
		case buffalo.EvtRouteErr:
			dd.routeError(e)
		case buffalo.EvtRouteFinished:
			dd.routeFinished(e)
		}
	})

	return err
}

func (dd *monitor) appStart(e events.Event) {
	addr := net.JoinHostPort(
		dd.Host,
		strconv.Itoa(dd.Port),
	)

	tracer.Start(tracer.WithAgentAddr(addr))
}

func (dd *monitor) appStop(e events.Event) {
	defer tracer.Stop()
}

func (dd *monitor) routeStarted(e events.Event) {

	ro, err := e.Payload.Pluck("route")
	if err != nil {
		fmt.Printf("datadog monitor: error getting the route: %v\n", err.Error())
		return
	}

	currentRoute := ro.(buffalo.RouteInfo)

	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		fmt.Printf("datadog monitor: error getting the context: %v\n", err.Error())
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
		c.Logger().Errorf("datadog monitor: extracting headers: %v\n", err.Error())
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span := tracer.StartSpan("http.request", opts...)
	c.Set("span", span)
}

func (dd *monitor) routeError(e events.Event) {
	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		fmt.Printf("datadog monitor: error getting the context: %v\n", err.Error())
		return
	}

	c := ctx.(buffalo.Context)

	var span tracer.Span
	var ok bool
	if span, ok = c.Value("span").(tracer.Span); !ok {
		c.Logger().Errorf("datadog monitor: error getting the span: %v\n", err.Error())
		return
	}

	span.SetTag(ext.Error, e.Error.(error))
}

func (dd *monitor) routeFinished(e events.Event) {
	ctx, err := e.Payload.Pluck("context")
	if err != nil {
		fmt.Printf("datadog monitor: error getting the context: %v\n", err.Error())
		return
	}

	c := ctx.(buffalo.Context)

	var response *buffalo.Response
	var span tracer.Span
	var ok bool

	if span, ok = c.Value("span").(tracer.Span); !ok {
		c.Logger().Errorf("datadog monitor: error getting the span: %v\n", err.Error())
		return
	}

	if response, ok = c.Response().(*buffalo.Response); !ok {
		c.Logger().Errorf("datadog monitor: error getting the response: %v\n", err.Error())
		return
	}

	status := response.Status
	span.SetTag(ext.HTTPCode, strconv.Itoa(status))

	span.Finish()
}

// Track allows to track custom functions with DDog
func (dd *monitor) Track(name string, fn func() error) error {
	opts := []ddtrace.StartSpanOption{
		tracer.SpanType(ext.SpanType),
		tracer.ServiceName(dd.ServiceName),
		tracer.ResourceName(name),

		tracer.Tag(ext.EventSampleRate, math.NaN()),
		tracer.Tag("mux.host", dd.Host),
	}

	span := tracer.StartSpan(name, opts...)
	err := fn()
	if err != nil {
		span.SetTag(ext.Error, err)
	}

	span.Finish()
	return nil
}

// NewMonitor creates a new monitor for DataDog with the passed serviceName
func NewMonitor(serviceName, host string) tomus.APMMonitor {
	if host == "" {
		host = "localhost"
	}

	return &monitor{
		ServiceName: serviceName,
		Host:        host,
		Port:        8126,
	}
}
