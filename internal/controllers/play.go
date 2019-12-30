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
	// TODO If there are players but no `turn` then start a new game
	// TODO If the `Draw` card has been taken (is empty now), populate it with a new card
	// TODO If we're on round 3, finish the game after the last player's turn

	query, err := parseQuery(c.QueryString())
	if err != nil {
		return err
	}

	deck := prepDeck(query)
	deck.Debug()

	// TODO Send email to everyone when the game is over
	/*
		err := email.Send()
		if err != nil {
			log.Println(err)
		}
	*/

	return c.JSON(200, query)
}

func parseQuery(queryString string) (models.Query, error) {
	log.Println(queryString)

	stringifiedJson, err := url.QueryUnescape(queryString)
	if err != nil {
		return models.Query{}, err
	}

	var query models.Query
	json.Unmarshal([]byte(stringifiedJson), &query)

	return query, nil
}

func prepDeck(query models.Query) deck.Deck {
	log.Println("Making deck")
	deck := deck.New()
	// hand := deck.Deal(2)

	// Remove cards in the discard pile from the deck
	for _, card := range query.AllDiscards {
		deck.Remove(card)
	}

	// Remove cards in player hands from the deck
	for _, player := range query.AllPlayers {
		for _, card := range player.Hand {
			deck.Remove(card)
		}
	}

	deck.Shuffle()

	return deck
}
