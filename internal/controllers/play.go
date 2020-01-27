package controllers

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/compycore/sabacc/internal/deck"
	"github.com/compycore/sabacc/internal/email"
	"github.com/compycore/sabacc/internal/helpers"
	"github.com/compycore/sabacc/internal/models"
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

	// Set the player names when they're not provided (for transitioning to names)
	// TODO Remove this
	database = tempSetPlayerNames(database)

	// Build a new deck of cards and make it match the contents of the database (remove cards from the players' hands, populate the discard pile, etc.)
	gameDeck := prepDeck(database)

	// Put a card on top of the draw pile if necessary
	gameDeck, database.Draw = populateDraw(database, gameDeck)

	// Pre-roll the dice to prevent browser refresh cheating
	database.Dice = rollDice()

	// Populate the discard pile if it's currently empty
	database, gameDeck = populateDiscard(database, gameDeck)

	// Deal hands if we're starting a new game or if hands were discarded because of a dice roll (notify players if hands were discarded)
	database, gameDeck, err = dealHands(database, gameDeck)
	if err != nil {
		return models.Database{}, err
	}

	// Calculate player scores
	database = calculatePlayerScores(database)

	// Send an email confirmation to the player that just took their turn and notify players of new games starting
	err = sendNotices(database)
	if err != nil {
		return models.Database{}, err
	}

	// We use round 0 as an indicator that a game is just starting; we need to increase that to round 1 now that everything is set up
	database = endRoundZero(database)

	// Change whose turn it is (also increase the round and change the dealer if necessary)
	database = endTurn(database)

	// Send emails
	if !isGameOver(database) {
		// Since the game is not over, notify the next player that it's their turn
		err = sendNextTurnEmail(database)
		if err != nil {
			return models.Database{}, err
		}
	} else {
		// TODO Determine who won
		// Set a value to database.Result so the tests can know that a match finished
		database.Result, err = generateResultString(database)
		if err != nil {
			return models.Database{}, err
		}

		// Send a game over email to every player
		err = sendGameOverEmails(database)
		if err != nil {
			return models.Database{}, err
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

func tempSetPlayerNames(database models.Database) models.Database {
	for i, player := range database.AllPlayers {
		if player.Email == "hellojessemillar@gmail.com" {
			database.AllPlayers[i].Name = "Jesse"
		} else if player.Email == "rileyjmillar@gmail.com" {
			database.AllPlayers[i].Name = "Riley"
		} else if player.Email == "jameston2001@gmail.com" {
			database.AllPlayers[i].Name = "James"
		}
	}

	return database
}

func populateDraw(database models.Database, gameDeck deck.Deck) (deck.Deck, deck.Card) {
	if database.Draw == (deck.Card{}) {
		gameDeck, database.Draw = deck.DealSingle(gameDeck)
		return deck.DealSingle(gameDeck)
	}

	return gameDeck, database.Draw
}

func endTurn(database models.Database) models.Database {
	if database.Rolled && database.Turn == database.Dealer {
		database.Round = database.Round + 1
		database.Dealer = changeDealer(database)
		database.Turn = database.Dealer
		database.Rolled = false
	}

	database.Turn = increaseTurn(database)

	return database
}

func increaseTurn(database models.Database) int {
	if database.Turn+1 < len(database.AllPlayers) {
		return database.Turn + 1
	}

	return 0
}

func changeDealer(database models.Database) int {
	if database.Dealer+1 < len(database.AllPlayers) {
		return database.Dealer + 1
	}

	return 0
}

func dealHands(database models.Database, gameDeck deck.Deck) (models.Database, deck.Deck, error) {
	if len(database.AllPlayers[0].Hand) == 0 {
		for i, _ := range database.AllPlayers {
			cardCount := database.AllPlayers[i].HandSize
			if cardCount == 0 {
				cardCount = 2
			}

			gameDeck, database.AllPlayers[i].Hand = deck.Deal(gameDeck, cardCount)
		}

		if database.Round > 0 {
			// Put an extra card on top of the discard pile
			database, gameDeck = dealIntoDiscard(database, gameDeck)

			err := sendHandDiscardEmails(database)
			if err != nil {
				return models.Database{}, deck.Deck{}, err
			}
		}
	}

	return database, gameDeck, nil
}

func populateDiscard(database models.Database, gameDeck deck.Deck) (models.Database, deck.Deck) {
	if len(database.AllDiscards) == 0 {
		database, gameDeck = dealIntoDiscard(database, gameDeck)
	}

	return database, gameDeck
}

func dealIntoDiscard(database models.Database, gameDeck deck.Deck) (models.Database, deck.Deck) {
	discard := deck.Card{}
	gameDeck, discard = deck.DealSingle(gameDeck)
	database.AllDiscards = append(database.AllDiscards, discard)

	return database, gameDeck
}

func endRoundZero(database models.Database) models.Database {
	if database.Round == 0 {
		database.Round = 1
	}

	return database
}

func sendNotices(database models.Database) error {
	if hasGameStarted(database) && !isGameOver(database) {
		return email.SendConfirmation(database.AllPlayers[database.Turn].Email, database.Codename, getHandString(database.AllPlayers[database.Turn].Hand), strconv.Itoa(database.AllPlayers[database.Turn].Score))
	} else if !hasGameStarted(database) {
		return email.SendGameStartNotice(database)
	}

	return nil
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

func sendHandDiscardEmails(database models.Database) error {
	// Send an email to every player
	for _, player := range database.AllPlayers {
		// TODO Make the function smart enough to not need both HTML and plain if only plain is passed
		// TODO Decide if I'm gonna handle text-only emails or if HTML is required
		err := email.SendHandDiscardNotice(player.Email, database.Codename, getHandString(player.Hand), strconv.Itoa(player.Score))
		if err != nil {
			return err
		}
	}

	return nil
}

func sendGameOverEmails(database models.Database) error {
	// Send an email to every player
	for _, player := range database.AllPlayers {
		// TODO Make the function smart enough to not need both HTML and plain if only plain is passed
		// TODO Decide if I'm gonna handle text-only emails or if HTML is required
		err := email.SendMessage(player.Email, database.Codename, database.Result, database.Result)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateResultString(database models.Database) (string, error) {
	rematchLink, err := generateRematchLink(database)
	if err != nil {
		return "", err
	}

	return generateHandSummaries(database) + "<br><br>" + rematchLink, nil
}

func generateHandSummaries(database models.Database) string {
	finalResultsMessage := ""
	for _, player := range database.AllPlayers {
		finalResultsMessage = finalResultsMessage + player.Email + " got a final score of " + strconv.Itoa(player.Score) + " with a hand of " + getHandString(player.Hand)
	}

	return finalResultsMessage
}

func generateRematchLink(database models.Database) (string, error) {
	rematchDatabase := models.Database{}
	rematchDatabase.Rematch = database.AllPlayers

	rematchDatabaseString, err := encodeDatabase(rematchDatabase)
	if err != nil {
		return "", err
	}

	return `<a href="` + os.Getenv("SABACC_UI_HREF") + "?" + rematchDatabaseString + `">Click here for a rematch!</a>`, nil
}

func rollDice() []int {
	dice := []int{}

	for i := 0; i < 2; i++ {
		dice = append(dice, rand.Intn(6)+1)
	}

	return dice
}
