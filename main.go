package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jessemillar/health"
	"github.com/jessemillar/sabacc/internal/controllers"
	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/email"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("Making deck")
	deck := deck.New()
	deck.Debug()
	log.Println("Shuffling deck")
	deck.Shuffle()
	deck.Debug()
	log.Println("Dealing hand")
	hand := deck.Deal(2)
	hand.Debug()
	// TODO Write unit tests
	log.Println("Checking deck length")
	deck.Debug()

	err := email.Send()
	if err != nil {
		log.Println(err)
	}

	log.Println("Configuring server")

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))
	e.PUT("/play", controllers.Play)

	// TODO Fail if $SABACC_PORT is not set
	e.Logger.Fatal(e.Start(":" + os.Getenv("SABACC_PORT")))
}
