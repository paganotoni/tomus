package middleware

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

requestIDToken := "X-Unique-Id"

// RequestID looks for a X-Unique-Id header in the request
// if its not present it generates a UUID and adds a log field 
// with that token and value.
func RequestID(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		id := c.Request().Header.Get(requestIDToken)
		if id == "" {
			id = uuid.Must(uuid.NewV4()).String()
		}

		c.Set(requestIDToken, id)
		c.LogField(requestIDToken, id)

		return next(c)
	}
}
