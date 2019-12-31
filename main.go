package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jessemillar/health"
	"github.com/jessemillar/sabacc/internal/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))
	// TODO Change to PUT
	e.GET("/sabacc", controllers.Play)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
