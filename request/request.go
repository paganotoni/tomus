package request

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

var (
	requestIdentifier = "X-Request-ID"
)

func MountTo(app *buffalo.App) {
	app.Use(middleware)
}

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
