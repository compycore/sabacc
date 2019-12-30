package controllers

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/models"
	"github.com/labstack/echo"
)

func Play(c echo.Context) error {
	// Get all variable values from the URI params
	round := c.QueryParam("round")
	log.Println(round)
	turn := c.QueryParam("turn")
	log.Println(turn)
	allPlayers := c.QueryParam("player")
	log.Println(allPlayers)

	// TODO If there are players but no `turn` then start a new game

	// TODO If we're on round 3, finish the game after the last player's turn

	log.Println("Making deck")
	deck := deck.New()
	// deck.Debug()
	log.Println("Shuffling deck")
	deck.Shuffle()
	// deck.Debug()
	log.Println("Dealing hand")
	hand := deck.Deal(2)
	hand.Debug()
	// TODO Write unit tests
	log.Println("Checking deck length")
	// deck.Debug()

	/*
		err := email.Send()
		if err != nil {
			log.Println(err)
		}
	*/

	log.Println(c.QueryString())

	stringifiedJson, err := url.QueryUnescape(c.QueryString())
	if err != nil {
		return err
	}

	var query models.Query
	json.Unmarshal([]byte(stringifiedJson), &query)

	// TODO Do I need/want to send a response back?
	return c.JSON(200, query)
}
