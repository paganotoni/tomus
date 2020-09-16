package datadog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	mon := NewMonitor("service", "").(*monitor)

	r.Equal("service", mon.ServiceName)
	r.Equal("localhost", mon.Host)
	r.Equal(8126, mon.Port)
}

func Test_Track(t *testing.T) {
	r := require.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	var flag int
	mon := NewMonitor("service", "")
	err := mon.Track("backgroundThing", func() error {
		flag = 10
		return nil
	})

	spans := mt.FinishedSpans()

	r.NoError(err)
	r.Equal(len(spans), 1)
	r.Equal(10, flag)

	err = mon.Track("backgroundThing", func() error {
		return errors.New("sample error")
	})

	r.NoError(err)

	spans = mt.FinishedSpans()
	r.Equal(len(spans), 2)

	errs := spans[1].Tag(ext.Error)
	r.NotNil(errs)

	r.Error(errs.(error))
	r.Equal(errs.(error).Error(), "sample error")
	r.Equal(spans[1].Tag("mux.host").(string), "localhost")
}
