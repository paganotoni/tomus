package request

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
)

var (
	requestIdentifier = "X-Request-ID"
)

var renderEngine *render.Engine

func MountTo(app *buffalo.App, r *render.Engine) {
	renderEngine = r
	app.Use(middleware)
	app.GET("/admin/info", healthCheck)
	app.Middleware.Clear()
}

// healthCheck allows to check if the app is ready to respond.
func healthCheck(c buffalo.Context) error {
	return c.Render(200, renderEngine.String("OK"))
}

//middleware takes care of extracting the request id from the
//request. and putting it into the context.
func middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		id := c.Request().Header.Get(requestIdentifier)
		if id != "" {
			id = uuid.Must(uuid.NewV4()).String()
		}

		c.Set("X-Request-ID", id)
		return next(c)
	}
}
