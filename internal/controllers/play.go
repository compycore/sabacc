package controllers

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/email"
	"github.com/jessemillar/sabacc/internal/models"
	"github.com/labstack/echo"
)

func Play(c echo.Context) error {
	result, err := gameLoop(c.QueryString())
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, result)
}

// Broken into a function for easier unit testing (don't have to mock an Echo context this way)
func gameLoop(queryString string) (models.Database, error) {
	database, err := parseDatabase(queryString)
	if err != nil {
		return models.Database{}, err
	}

	gameDeck := prepDeck(database)

	if database.Draw == (deck.Card{}) {
		database.Draw = gameDeck.Deal(1)[0]
	}

	if len(database.AllDiscards) == 0 {
		database.AllDiscards = append(database.AllDiscards, gameDeck.Deal(1)[0])
	}

	if database.Round > 0 {
		database.Turn = database.Turn + 1
		if database.Turn >= len(database.AllPlayers) {
			database.Turn = 0
		}
	}

	// Start a new game if needed
	if len(database.AllPlayers[0].Hand) == 0 {
		for i, _ := range database.AllPlayers {
			database.AllPlayers[i].Hand = gameDeck.Deal(2)
		}

		// Only set the round to 1 if it's a new game (as opposed to deleting hands because the dice were doubles)
		if database.Round == 0 {
			database.Round = 1
		}
	}

	if database.Turn == len(database.AllPlayers)-1 {
		database.Round = database.Round + 1
	}

	// Calculate player scores
	database = calculatePlayerScores(database)

	// If the game is still going
	if database.Round <= 3 && database.Turn < len(database.AllPlayers) {
		encodedDatabase, err := encodeDatabase(database)
		if err != nil {
			return models.Database{}, err
		}

		err = email.SendLink(database.AllPlayers[database.Turn].Email, os.Getenv("SABACC_UI_HREF")+"?"+encodedDatabase, database.Round)
		if err != nil {
			return models.Database{}, err
		}
	} else {
		// TODO Determine who won
		for _, player := range database.AllPlayers {
			// TODO Make the function smart enough to not need both HTML and plain if only plain is passed
			err = email.SendMessage(player.Email, "Game over", "Game over")
			if err != nil {
				return models.Database{}, err
			}
		}
	}

	return database, nil
}

func parseDatabase(databaseString string) (models.Database, error) {
	stringifiedJson, err := url.QueryUnescape(databaseString)
	if err != nil {
		return models.Database{}, err
	}

	var database models.Database
	json.Unmarshal([]byte(stringifiedJson), &database)

	return database, nil
}

func encodeDatabase(database models.Database) (string, error) {
	encodedDatabase, err := json.Marshal(database)
	if err != nil {
		return "", err
	}

	return url.QueryEscape(string(encodedDatabase)), nil
}

func calculatePlayerScores(database models.Database) models.Database {
	for i, player := range database.AllPlayers {
		database.AllPlayers[i].Score = 0

		for _, card := range player.Hand {
			database.AllPlayers[i].Score = database.AllPlayers[i].Score + card.Value
		}
	}

	return database
}

func prepDeck(database models.Database) deck.Deck {
	preppedDeck := deck.New()

	// Remove cards in the discard pile from the deck
	for _, card := range database.AllDiscards {
		preppedDeck.Remove(card)
	}

	// Remove cards in player hands from the deck
	for _, player := range database.AllPlayers {
		for _, card := range player.Hand {
			preppedDeck.Remove(card)
		}
	}

	// Remove the card that's available to draw from the deck
	preppedDeck.Remove(database.Draw)

	preppedDeck.Shuffle()

	return preppedDeck
}
