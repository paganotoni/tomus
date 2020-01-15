package tomus

import (
	"github.com/paganotoni/tomus/datadog"
	"github.com/paganotoni/tomus/newrelic"
)

var (
	// Ensuring datadog Monitor meets the interface
	_ APMMonitor = datadog.NewMonitor("aaa")

	// Ensuring newrelic Monitor meets the interface
	_ APMMonitor = newrelic.NewMonitor("aaa", "development", "aaaaa")
)
