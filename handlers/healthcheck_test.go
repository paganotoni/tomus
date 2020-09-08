package request

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_HealthCheck(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.GET("/admin/info", HealthCheck)

	ht := httptest.New(app)
	res := ht.HTML("/admin/info").Get()

	r.Equal(http.StatusOK, res.Code, "/admin should respond OK status")
	r.Equal("OK", res.Body.String(), "/admin should respond OK content")
}
