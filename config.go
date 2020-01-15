package tomus

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

const (
	// APMKindNewrelic defines newrelic usage
	APMKindNewrelic = iota + 1

	// APMKindDatadog defines datadog usage
	APMKindDatadog = iota + 1
)

//Config is used to start tomus monitoring
type Config struct {
	App          *buffalo.App
	RenderEngine *render.Engine

	APMKind     int
	ServiceName string
	Environment string
	EnableAPM   bool
}
