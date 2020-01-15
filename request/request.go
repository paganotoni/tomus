package request

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gofrs/uuid"
)

//MountTo adds
func MountTo(app *buffalo.App) {
	app.Use(middleware)

	admin := app.Group("/admin")
	admin.GET("/info", healthCheck)
	admin.Middleware.Clear()
}

// healthCheck allows to check if the app is ready to respond.
func healthCheck(c buffalo.Context) error {
	r := &render.Engine{}
	return c.Render(200, r.String("OK"))
}

//middleware takes care of extracting the request id from the
//request. and putting it into the context.
func middleware(next buffalo.Handler) buffalo.Handler {
	requestIdentifier := "X-Unique-Id"

	return func(c buffalo.Context) error {
		id := c.Request().Header.Get(requestIdentifier)
		if id == "" {
			id = uuid.Must(uuid.NewV4()).String()
		}

		c.Set(requestIdentifier, id)
		c.LogField(requestIdentifier, id)

		return next(c)
	}
}
