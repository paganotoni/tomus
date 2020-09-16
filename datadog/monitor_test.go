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
	mon.Track("backgroundThing", func() error {
		flag = 10
		return nil
	})

	spans := mt.FinishedSpans()

	r.Equal(len(spans), 1)
	r.Equal(10, flag)

	mon.Track("backgroundThing", func() error {
		return errors.New("sample error")
	})

	spans = mt.FinishedSpans()
	r.Equal(len(spans), 2)

	err := spans[1].Tag(ext.Error)
	r.NotNil(err)

	r.Error(err.(error))
	r.Equal(err.(error).Error(), "sample error")
	r.Equal(spans[1].Tag("mux.host").(string), "localhost")
}
