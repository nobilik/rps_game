package auth

import "github.com/labstack/echo/v4"

func AddRoutes(api *echo.Group) {
	api.POST("/register", Register)
	api.POST("/login", Login)
	api.DELETE("/logout", Logout)
}
