package remote

import (
	"TBack/conf"

	"github.com/labstack/echo"
)

// ValidateRemote ...
func ValidateRemote(username, password string, c echo.Context) (bool, error) {
	if username == conf.TBack.GetString("remote_user") && password == conf.TBack.GetString("remote_pass") {
		return true, nil
	}
	return false, nil
}
