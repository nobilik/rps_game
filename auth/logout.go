package auth

import (
	"net/http"
	"rps/models/db"

	"github.com/labstack/echo/v4"
)

// @tags auth
// @summary Logout
// @description Deletes current auth session
// @param X-Auth-Token header string true "token"
// @success 200 {string} string "logged out"
// @failure 400 {string} string "Some err"
// @router       /logout [delete]
func Logout(c echo.Context) error {
	token, err := Token(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	db.RedisClient.Del("SESSION:" + token)

	return c.JSON(http.StatusOK, "logged out")
}
