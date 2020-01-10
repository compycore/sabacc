package controllers

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"

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
		gameDeck, database.Draw = deck.DealSingle(gameDeck)
	}

	if len(database.AllDiscards) == 0 {
		discard := deck.Card{}
		gameDeck, discard = deck.DealSingle(gameDeck)
		database.AllDiscards = append(database.AllDiscards, discard)
	}

	if database.Round > 0 {
		database.Turn = database.Turn + 1
	}

	// Start a new game if needed
	if len(database.AllPlayers[0].Hand) == 0 {
		for i, _ := range database.AllPlayers {
			gameDeck, database.AllPlayers[i].Hand = deck.Deal(gameDeck, 2)
		}

		// Only set the round to 1 if it's a new game (as opposed to deleting hands because the dice were doubles)
		if database.Round == 0 {
			database.Round = 1
		}
	}

	if database.Turn == len(database.AllPlayers) {
		database.Round = database.Round + 1
		database.Turn = 0
	}

	// Calculate player scores
	database = calculatePlayerScores(database)

	// Send an email confirmation to the player that just took their turn
	if database.Round > 0 && database.Turn > 0 {
		previousTurn := database.Turn - 1
		if previousTurn < 0 {
			previousTurn = len(database.AllPlayers) - 1
		}
		email.SendConfirmation(database.AllPlayers[previousTurn].Email, getHandString(database.AllPlayers[previousTurn].Hand), strconv.Itoa(database.AllPlayers[previousTurn].Score))
	}

	// If the game is still going
	if database.Round <= 3 && database.Turn < len(database.AllPlayers) && len(database.AllPlayers) > 1 {
		encodedDatabase, err := encodeDatabase(database)
		if err != nil {
			return models.Database{}, err
		}

		allEmailAddresses := ""
		for _, player := range database.AllPlayers {
			allEmailAddresses = allEmailAddresses + player.Email + ", "
		}

		err = email.SendLink(database.AllPlayers[database.Turn].Email, allEmailAddresses, os.Getenv("SABACC_UI_HREF")+"?"+encodedDatabase, database.Round)
		if err != nil {
			return models.Database{}, err
		}
	} else {
		// TODO Determine who won
		// TODO Break this into a function
		finalResultsMessage := ""
		for _, player := range database.AllPlayers {
			finalResultsMessage = finalResultsMessage + player.Email + " got a final score of " + strconv.Itoa(player.Score) + " with a hand of " + getHandString(player.Hand) + "\n\n"
		}

		rematchDatabase := models.Database{}
		rematchDatabase.Rematch = database.AllPlayers

		rematchDatabaseString, err := encodeDatabase(rematchDatabase)
		if err != nil {
			return models.Database{}, err
		}

		finalResultsMessage = finalResultsMessage + `\n\n<a href="` + os.Getenv("SABACC_UI_HREF") + "?" + rematchDatabaseString + `">Click here for a rematch!</a>`

		// Send an email to every player
		for _, player := range database.AllPlayers {
			// TODO Make the function smart enough to not need both HTML and plain if only plain is passed
			err = email.SendMessage(player.Email, finalResultsMessage, finalResultsMessage)
			if err != nil {
				return models.Database{}, err
			}
		}

		database.Result = finalResultsMessage
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
		preppedDeck = deck.Remove(preppedDeck, card)
	}

	// Remove cards in player hands from the deck
	for _, player := range database.AllPlayers {
		for _, card := range player.Hand {
			preppedDeck = deck.Remove(preppedDeck, card)
		}
	}

	// Remove the card that's available to draw from the deck
	preppedDeck = deck.Remove(preppedDeck, database.Draw)

	preppedDeck = deck.Shuffle(preppedDeck)

	return preppedDeck
}

func getHandString(hand deck.Deck) string {
	handString := ""

	for _, card := range hand {
		handString = handString + "\n" + card.Stave + " " + strconv.Itoa(card.Value)
	}

	return handString
}
