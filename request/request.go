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
