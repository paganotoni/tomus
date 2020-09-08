package request

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// HealthCheck allows to check if the app is ready to respond.
func HealthCheck(c buffalo.Context) error {
	r := render.New(render.Options{})
	return c.Render(200, r.String("OK"))
}
