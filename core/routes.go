package core

import "github.com/labstack/echo/v4"

func AddRoutes(api *echo.Group) {
	api.GET("/users", GetUsers)
	api.GET("/bets", GetBets)
	api.POST("/bets", PlaceBet)
	api.PATCH("/bets/:action", RespondToBet)

	api.GET("/transactions", GetPersonalTransactions)
	api.POST("/transactions/:action", ManageBalance)

	api.GET("/stats", AllStats)
}
