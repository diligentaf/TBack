package middleware

import (
	"github.com/labstack/echo"
)

// ResponseDump print http response body
func ResponseDump(c echo.Context, _, resBody []byte) {
	trID := c.Response().Header().Get(echo.HeaderXRequestID)
	if trID == "" {
		return
	}

	mlog.Infow("[ResponseDump]", "method", c.Request().Method, "trID", trID, "url", c.Path(), "status", c.Response().Status)
}
