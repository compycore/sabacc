package controllers

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/email"
	"github.com/jessemillar/sabacc/internal/helpers"
	"github.com/jessemillar/sabacc/internal/models"
	"github.com/labstack/echo"
)

// Play is the main function that is called by Echo from main.go
func Play(c echo.Context) error {
	result, err := gameLoop(c.QueryString())
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, result)
}

// Broken into a function for easier unit testing (don't have to mock an Echo context this way)
func gameLoop(queryString string) (models.Database, error) {
	// Parse the encoded-as-a-query-string database into a models.Database{} struct
	database, err := parseDatabase(queryString)
	if err != nil {
		return models.Database{}, err
	}

	// We use codenames to prevent all Sabacc emails from being threaded into the same thread in Gmail
	database.Codename = getCodename(database)

	// Build a new deck of cards and make it match the contents of the database (remove cards from the players' hands, populate the discard pile, etc.)
	gameDeck := prepDeck(database)

	// Put a card on top of the draw pile if necessary
	gameDeck, database.Draw = populateDraw(database, gameDeck)

	// Populate the discard pile
	database, gameDeck = populateDiscard(database, gameDeck)

	// Deal hands if we're starting a new game or if hands were discarded because of a dice roll
	database, gameDeck = dealHands(database, gameDeck)

	// Calculate player scores
	database = calculatePlayerScores(database)

	// Send an email confirmation to the player that just took their turn
	sendTurnConfirmationEmail(database)

	// We use round 0 as an indicator that a game is just starting; we need to increase that to round 1 now that everything is set up
	database.Round = endRoundZero(database)

	// Change whose turn it is
	database.Turn = increaseTurn(database)

	// Send emails
	if !isGameOver(database) {
		err = sendNextTurnEmail(database)
		if err != nil {
			return models.Database{}, err
		}
	} else {
		// TODO Determine who won
		finalResultsMessage := ""
		for _, player := range database.AllPlayers {
			finalResultsMessage = finalResultsMessage + player.Email + " got a final score of " + strconv.Itoa(player.Score) + " with a hand of " + getHandString(player.Hand)
		}

		rematchDatabase := models.Database{}
		rematchDatabase.Rematch = database.AllPlayers

		rematchDatabaseString, err := encodeDatabase(rematchDatabase)
		if err != nil {
			return models.Database{}, err
		}

		finalResultsMessage = finalResultsMessage + `<br><br><a href="` + os.Getenv("SABACC_UI_HREF") + "?" + rematchDatabaseString + `">Click here for a rematch!</a>`

		// Send an email to every player
		for _, player := range database.AllPlayers {
			// TODO Make the function smart enough to not need both HTML and plain if only plain is passed
			err = email.SendMessage(player.Email, database.Codename, finalResultsMessage, finalResultsMessage)
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
		handString = handString + card.Stave + " " + strconv.Itoa(card.Value) + "<br>"
	}

	return handString
}

func hasGameStarted(database models.Database) bool {
	return database.Round > 0
}

func isGameOver(database models.Database) bool {
	if database.Round <= 3 && database.Turn < len(database.AllPlayers) && len(database.AllPlayers) > 1 {
		return false
	}

	return true
}

func getCodename(database models.Database) string {
	if len(database.Codename) == 0 {
		database.Codename = helpers.GetCodename()
	}

	// Remove "+" characters that get put there by the Golang marshal process (Go encodes spaces as "+" instead of "%20")
	return strings.ReplaceAll(database.Codename, "+", " ")
}

func populateDraw(database models.Database, gameDeck deck.Deck) (deck.Deck, deck.Card) {
	if database.Draw == (deck.Card{}) {
		gameDeck, database.Draw = deck.DealSingle(gameDeck)
		return deck.DealSingle(gameDeck)
	}

	return gameDeck, database.Draw
}

func increaseTurn(database models.Database) int {
	if database.Round > 0 {
		database.Turn = database.Turn + 1
	}

	if database.Turn == len(database.AllPlayers) {
		database.Round = database.Round + 1
		database.Turn = 0
	}

	return database.Turn
}

func dealHands(database models.Database, gameDeck deck.Deck) (models.Database, deck.Deck) {
	if len(database.AllPlayers[0].Hand) == 0 {
		for i, _ := range database.AllPlayers {
			gameDeck, database.AllPlayers[i].Hand = deck.Deal(gameDeck, 2)
		}
	}

	return database, gameDeck
}

func populateDiscard(database models.Database, gameDeck deck.Deck) (models.Database, deck.Deck) {
	if len(database.AllDiscards) == 0 {
		discard := deck.Card{}
		gameDeck, discard = deck.DealSingle(gameDeck)
		database.AllDiscards = append(database.AllDiscards, discard)
	}

	return database, gameDeck
}

func endRoundZero(database models.Database) int {
	if database.Round == 0 {
		return 1
	}

	return database.Round
}

func sendTurnConfirmationEmail(database models.Database) {
	if hasGameStarted(database) && !isGameOver(database) {
		email.SendConfirmation(database.AllPlayers[database.Turn].Email, database.Codename, getHandString(database.AllPlayers[database.Turn].Hand), strconv.Itoa(database.AllPlayers[database.Turn].Score))
	}
}

func sendNextTurnEmail(database models.Database) error {
	encodedDatabase, err := encodeDatabase(database)
	if err != nil {
		return err
	}

	allEmailAddresses := ""
	for _, player := range database.AllPlayers {
		allEmailAddresses = allEmailAddresses + player.Email + ", "
	}

	err = email.SendLink(database.AllPlayers[database.Turn].Email, allEmailAddresses, database.Codename, os.Getenv("SABACC_UI_HREF")+"?"+encodedDatabase, database.Round)
	if err != nil {
		return err
	}

	return nil
}
