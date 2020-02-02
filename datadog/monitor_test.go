package datadog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	mon := NewMonitor("service", "")

	r.Equal("service", mon.ServiceName)
	r.Equal("localhost", mon.Host)
	r.Equal(8126, mon.Port)
}
