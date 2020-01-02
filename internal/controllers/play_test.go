package controllers

import (
	"encoding/json"
	"testing"

	"github.com/jessemillar/sabacc/internal/deck"
	"github.com/jessemillar/sabacc/internal/models"
)

var databaseStruct = models.Database{
	Round: 1,
	Turn:  0,
	AllPlayers: []models.Player{
		{
			Email: "hellojessemillar@gmail.com",
			Hand: deck.Deck{
				{1, "triangle"},
				{3, "circle"},
			},
		},
		{
			Email: "penguinshatestuff@gmail.com",
			Hand: deck.Deck{
				{-3, "square"},
				{5, "circle"},
			},
		},
	},
	Draw: deck.Card{5, "triangle"},
	AllDiscards: deck.Deck{
		{10, "triangle"},
		{10, "circle"},
	},
}

func TestParseDatabase(t *testing.T) {
	encodedDatabase, err := encodeDatabase(databaseStruct)
	if err != nil {
		t.Error(err)
	}

	resultDatabase, err := parseDatabase(encodedDatabase)
	// database, err := ParseDatabase("%7B%22round%22%3A1%2C%22turn%22%3A1%2C%22players%22%3A%5B%7B%22email%22%3A%22hellojessemillar%40gmail.com%22%2C%22hand%22%3A%5B%7B%22value%22%3A10%2C%22stave%22%3A%22triangle%22%7D%2C%7B%22value%22%3A-9%2C%22stave%22%3A%22circle%22%7D%5D%7D%2C%7B%22email%22%3A%22penguinshatestuff%40gmail.com%22%2C%22hand%22%3A%5B%7B%22value%22%3A-8%2C%22stave%22%3A%22circle%22%7D%2C%7B%22value%22%3A10%2C%22stave%22%3A%22circle%22%7D%5D%7D%5D%2C%22draw%22%3A%7B%22value%22%3A8%2C%22stave%22%3A%22square%22%7D%2C%22discards%22%3Anull%7D")
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
