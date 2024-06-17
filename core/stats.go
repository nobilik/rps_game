package core

import (
	"net/http"
	"rps/models/db"

	"github.com/labstack/echo/v4"
)

// stats for admins overal users bets
// @tags stats
// @summary gives admin stats for all users
// @description just example
// @produce json
// @success 200 {object} object
// @router       /stats [get]
func AllStats(c echo.Context) error {
	byStatus := db.GetBetStatsByStatus()
	byResult := db.GetBetStatsByResult()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"byStatus": byStatus,
		"byResult": byResult,
	})

}
