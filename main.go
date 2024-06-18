// @title RPS game API - OpenAPI 2.0
// @description simple RPS game
// @license name:Apache 2.0 url:http://www.apache.org/licenses/LICENSE-2.0.html
// @version 1.0.0
// @Schemes http
// @host      localhost:3000
// @BasePath  /api/v1

// swag init -g main.go -o docs --parseDependency --parseInternal

package main

import (
	"net/http"
	"rps/auth"
	"rps/core"
	"rps/models/db"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	swg "github.com/swaggo/echo-swagger"

	_ "rps/docs"
)

func main() {
	time.Sleep(10 * time.Second)
	db.SetDB()
	db.SetRedis()
	defer db.RedisClient.Close()

	server()
}

func server() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // TODO: Allow strict origins
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	api := e.Group("/api/v1")

	auth.AddRoutes(api)
	core.AddRoutes(api)

	e.GET("/docs/*", swg.WrapHandler)

	// go core.PlaceBetByAI()

	e.Logger.Fatal(e.Start(":3000"))
}
