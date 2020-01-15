package tomus

import (
	"github.com/gobuffalo/buffalo"
)

const (
	// APMKindNewrelic defines newrelic usage
	APMKindNewrelic = iota + 1

	// APMKindDatadog defines datadog usage
	APMKindDatadog
)

//Config is used to start Tomus monitoring
type Config struct {
	//App is the application where we would mount routes for monitoring
	App *buffalo.App

	//The APM monitor we would use for
	APMMonitor APMMonitor
}
