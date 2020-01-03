package controllers

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/models"
)

var databaseStruct = models.Database{
	Round: 0,
	Turn:  0,
	AllPlayers: []models.Player{
		{
			Email: "hellojessemillar@gmail.com",
			Hand: deck.Deck{
				{
					1,
					"triangle",
				},
				{
					3,
					"circle",
				},
			},
		},
		{
			Email: "penguinshatestuff@gmail.com",
			Hand: deck.Deck{
				{
					-3,
					"square",
				},
				{
					5,
					"circle",
				},
			},
		},
	},
	Draw: deck.Card{
		5,
		"triangle",
	},
	AllDiscards: deck.Deck{
		{
			10,
			"triangle",
		},
		{
			10,
			"circle",
		},
	},
}

func TestParseDatabase(t *testing.T) {
	encodedDatabase, err := encodeDatabase(databaseStruct)
	if err != nil {
		t.Error(err)
	}

	resultDatabase, err := parseDatabase(encodedDatabase)
	if err != nil {
		t.Error(err)
	}

	input, err := json.MarshalIndent(databaseStruct, "", "\t")
	if err != nil {
		t.Error(err)
	}

	result, err := json.MarshalIndent(resultDatabase, "", "\t")
	if err != nil {
		t.Error(err)
	}

	if string(input) != string(result) {
		t.Errorf("Decoded database did not match source database: input: %s, result: %s", input, result)
	}
}

// This test is to recreate an out of range error
func TestProblematicDatabase(t *testing.T) {
	_, err := gameLoop("%7B%22round%22%3A1%2C%22turn%22%3A1%2C%22players%22%3A%5B%7B%22email%22%3A%22hellojessemillar%40gmail.com%22%2C%22hand%22%3A%5B%7B%22value%22%3A-6%2C%22stave%22%3A%22square%22%7D%2C%7B%22value%22%3A-3%2C%22stave%22%3A%22triangle%22%7D%5D%7D%2C%7B%22email%22%3A%22penguinshatestuff%40gmail.com%22%2C%22hand%22%3A%5B%7B%22value%22%3A9%2C%22stave%22%3A%22triangle%22%7D%2C%7B%22value%22%3A-5%2C%22stave%22%3A%22square%22%7D%5D%7D%5D%2C%22draw%22%3A%7B%22value%22%3A1%2C%22stave%22%3A%22triangle%22%7D%2C%22discards%22%3Anull%7D")
	if err != nil {
		t.Error(err)
	}
}

func TestPrepDeck(t *testing.T) {
	_ = prepDeck(databaseStruct)
}

func TestPrepDeckEmptyDatabase(t *testing.T) {
	_ = prepDeck(models.Database{})
}

func TestPrepDeckNewGame(t *testing.T) {
	_ = prepDeck(models.Database{
		AllPlayers: []models.Player{
			{
				Email: "test@test.com",
			},
			{
				Email: "test2@test.com",
			},
		},
	})
}

func databaseToURI(database models.Database) string {
	encodedDatabase, err := encodeDatabase(database)
	if err != nil {
		log.Fatal(err)
	}

	return encodedDatabase
}

func TestGameFlow(t *testing.T) {
	// ----------
	// Game init
	// ----------

	// An empty struct with only emails starts the game
	startingDatabase := models.Database{
		AllPlayers: []models.Player{
			{
				Email: "hellojessemillar@gmail.com",
			},
			{
				Email: "penguinshatestuff@gmail.com",
			},
		},
	}

	// Pass the bare database to the game loop to start the game
	resultDatabase, err := gameLoop(databaseToURI(startingDatabase))
	if err != nil {
		t.Error(err)
	}

	// Check that the round increased
	if resultDatabase.Round != 1 {
		t.Errorf("Round number incorrect; want: %d, got: %d", 1, resultDatabase.Round)
	}

	// Check that it's player 1's turn
	if resultDatabase.Turn != 0 {
		t.Errorf("Turn number incorrect; want: %d, got: %d", 0, resultDatabase.Turn)
	}

	// Verify that player 1 got a hand dealt to them
	if len(resultDatabase.AllPlayers[0].Hand) != 2 {
		t.Errorf("Player 1 hand size incorrect; want: %d, got: %d", 2, len(resultDatabase.AllPlayers[0].Hand))
	}

	// Verify that player 2 got a hand dealt to them
	if len(resultDatabase.AllPlayers[1].Hand) != 2 {
		t.Errorf("Player 2 hand size incorrect; want: %d, got: %d", 2, len(resultDatabase.AllPlayers[1].Hand))
	}

	// Make sure there's something in the discard pile
	if len(resultDatabase.AllDiscards) != 1 {
		t.Errorf("Discard pile size incorrect; want: %d, got: %d", 1, len(resultDatabase.AllDiscards))
	}

	// Make sure there's a card available to be drawn
	if resultDatabase.Draw == (deck.Card{}) {
		t.Error("There is no card available to be drawn")
	}

	// ----------
	// Round 1 - Player 1
	// ----------

	// Player 1 draws a card
	resultDatabase.AllPlayers[0].Hand = append(resultDatabase.AllPlayers[0].Hand, resultDatabase.Draw)
	resultDatabase.Draw = deck.Card{}

	// Send the new database to the game loop
	resultDatabase, err = gameLoop(databaseToURI(resultDatabase))
	if err != nil {
		t.Error(err)
	}

	// Verify that player 1 has the new card
	if len(resultDatabase.AllPlayers[0].Hand) != 3 {
		t.Errorf("Player 1 hand size incorrect; want: %d, got: %d", 3, len(resultDatabase.AllPlayers[0].Hand))
	}

	// Check that there's a new card in the draw pile
	if resultDatabase.Draw == (deck.Card{}) {
		t.Error("There's no card in the draw pile")
	}

	// Check that the round is the same
	if resultDatabase.Round != 1 {
		t.Errorf("Round number incorrect; want: %d, got: %d", 1, resultDatabase.Round)
	}

	// ----------
	// Round 1 - Player 2
	// ----------

	// Check that it's player 2's turn
	if resultDatabase.Turn != 1 {
		t.Errorf("Turn number incorrect; want: %d, got: %d", 1, resultDatabase.Turn)
	}

	// Player 2 swaps their card with the one in the discard pile
	// Save the two cards to different variables
	handSwap := resultDatabase.AllPlayers[1].Hand[len(resultDatabase.AllPlayers[1].Hand)-1]
	discardSwap := resultDatabase.AllDiscards[len(resultDatabase.AllDiscards)-1]
	// Remove the cards from the player's hand and discard pile
	resultDatabase.AllPlayers[1].Hand = deck.Remove(resultDatabase.AllPlayers[1].Hand, handSwap)
	resultDatabase.AllDiscards = deck.Remove(resultDatabase.AllDiscards, discardSwap)
	// Swap the cards
	resultDatabase.AllDiscards = append(resultDatabase.AllDiscards, handSwap)
	resultDatabase.AllPlayers[1].Hand = append(resultDatabase.AllPlayers[1].Hand, discardSwap)

	// Send the new database to the game loop
	resultDatabase, err = gameLoop(databaseToURI(resultDatabase))
	if err != nil {
		t.Error(err)
	}

	// Verify that the swap happened
	if deck.Contains(resultDatabase.AllDiscards, discardSwap) {
		t.Error("Discard pile still contains swapped card")
	}
	if deck.Contains(resultDatabase.AllPlayers[1].Hand, handSwap) {
		t.Error("Player hand still contains swapped card")
	}
	if !deck.Contains(resultDatabase.AllPlayers[1].Hand, discardSwap) {
		t.Error("Player hand does not contain newly swapped card")
	}
	if !deck.Contains(resultDatabase.AllDiscards, handSwap) {
		t.Error("Discard pile does not contain newly swapped card")
	}

	// Verify that player 2 has the correct hand size
	if len(resultDatabase.AllPlayers[1].Hand) != 2 {
		t.Errorf("Player 1 hand size incorrect; want: %d, got: %d", 2, len(resultDatabase.AllPlayers[0].Hand))
	}

	// Make sure the discard pile is the correct size
	if len(resultDatabase.AllDiscards) != 1 {
		t.Errorf("Discard pile size incorrect; want: %d, got: %d", 1, len(resultDatabase.AllDiscards))
	}

	// Check that the round increased
	if resultDatabase.Round != 2 {
		t.Errorf("Round number incorrect; want: %d, got: %d", 2, resultDatabase.Round)
	}

	// ----------
	// Round 2 - Player 1
	// ----------

	// Check that it's player 1's turn
	if resultDatabase.Turn != 0 {
		t.Errorf("Turn number incorrect; want: %d, got: %d", 1, resultDatabase.Turn)
	}

	// Player 1 stands (does nothing)

	// Send the untouched database to the game loop
	resultDatabase, err = gameLoop(databaseToURI(resultDatabase))
	if err != nil {
		t.Error(err)
	}

	// Check that the round stayed the same
	if resultDatabase.Round != 2 {
		t.Errorf("Round number incorrect; want: %d, got: %d", 2, resultDatabase.Round)
	}

	// ----------
	// Round 2 - Player 2
	// ----------

	// Check that it's player 2's turn
	if resultDatabase.Turn != 1 {
		t.Errorf("Turn number incorrect; want: %d, got: %d", 1, resultDatabase.Turn)
	}

	// ----------
	// Round 3
	// ----------

	// ----------
	// Game finish
	// ----------

	// TODO Finish a full testing scenario
}
