package datadog

import (
	"math"

	"github.com/paganotoni/tomus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

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

func (dd *monitor) TrackChild(options tomus.TrackingOptions, fn func() error) error {
	opts := []ddtrace.StartSpanOption{
		tracer.SpanType(ext.SpanType),
		tracer.ServiceName(dd.ServiceName),
		tracer.ResourceName(options.Name),

		tracer.Tag(ext.EventSampleRate, math.NaN()),
		tracer.Tag("mux.host", dd.Host),
	}

	if parentSpan, ok := options.ParentSpan.(tracer.Span); ok {
		opts = append(opts, tracer.ChildOf(parentSpan.Context()))
	}

	span := tracer.StartSpan(options.Name, opts...)
	err := fn()
	if err != nil {
		span.SetTag(ext.Error, err)
	}

	span.Finish()

	return nil
}
