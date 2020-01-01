package controllers

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/email"
	"github.com/jessemillar/sabacc/internal/models"
	"github.com/labstack/echo"
)

func Play(c echo.Context) error {
	// TODO Limit the game to 8 players

	database, err := ParseDatabase(c.QueryString())
	if err != nil {
		return err
	}

	gameDeck := prepDeck(database)

	if database.Draw == (deck.Card{}) {
		database.Draw = gameDeck.Deal(1)[0]
	}

	// Start a new game if needed
	if len(database.AllPlayers[0].Hand) == 0 {
		for i, _ := range database.AllPlayers {
			database.AllPlayers[i].Hand = gameDeck.Deal(2)
		}
	}

	database.Turn = database.Turn + 1
	if database.Turn > len(database.AllPlayers) {
		database.Turn = 0
	}

	database.Round = database.Round + 1

	if database.Round < 3 {
		encodedDatabase, err := EncodeDatabase(database)
		if err != nil {
			return err
		}

		err = email.Send(database.AllPlayers[database.Turn].Email, os.Getenv("SABACC_UI_HREF")+"?"+encodedDatabase)
		if err != nil {
			log.Println(err)
		}
	} else {
		// TODO Send email to everyone when the game is over
		// TODO Determine who won
	}

	return c.JSON(200, database)
}

func ParseDatabase(databaseString string) (models.Database, error) {
	log.Println(databaseString)

	stringifiedJson, err := url.QueryUnescape(databaseString)
	if err != nil {
		return models.Database{}, err
	}

	var database models.Database
	json.Unmarshal([]byte(stringifiedJson), &database)

	return database, nil
}

func EncodeDatabase(database models.Database) (string, error) {
	encodedDatabase, err := json.Marshal(database)
	if err != nil {
		return "", err
	}

	return url.QueryEscape(string(encodedDatabase)), nil
}

func prepDeck(database models.Database) deck.Deck {
	log.Println("Making deck")
	preppedDeck := deck.New()

	log.Println("Removing discarded cards from deck")
	// Remove cards in the discard pile from the deck
	for _, card := range database.AllDiscards {
		preppedDeck.Remove(card)
	}

	log.Println("Removing player hand cards from deck")
	// Remove cards in player hands from the deck
	for _, player := range database.AllPlayers {
		for _, card := range player.Hand {
			preppedDeck.Remove(card)
		}
	}

	log.Println("Shuffling")
	preppedDeck.Shuffle()

	return preppedDeck
}
