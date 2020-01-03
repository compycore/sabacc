package deck

import (
	"math/rand"

	"github.com/davecgh/go-spew/spew"
)

// Card holds the card staves and values in the deck
type Card struct {
	Value int    `json:"value"`
	Stave string `json:"stave"`
}

// Deck holds the cards in the deck
type Deck []Card

// New creates a deck of cards to be used
func New() Deck {
	allValues := []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	allStaves := []string{"circle", "square", "triangle"}

	deck := Deck{}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(allValues); i++ {
		for n := 0; n < len(allStaves); n++ {
			card := Card{
				Value: allValues[i],
				Stave: allStaves[n],
			}
			deck = append(deck, card)
		}
	}

	// Add the zero cards to the deck
	// TODO Make the zero card display properly in the UI
	for n := 0; n < 2; n++ {
		card := Card{
			Value: 0,
			Stave: "zero",
		}

		deck = append(deck, card)
	}

	return deck
}

// Remove a specific card from the deck
func Remove(deck Deck, card Card) Deck {
	for i, curCard := range deck {
		if curCard.Stave == card.Stave && curCard.Value == card.Value {
			deck = append((deck)[:i], (deck)[i+1:]...)
			return deck
		}
	}

	return deck
}

func Contains(deck Deck, card Card) bool {
	for _, curCard := range deck {
		if curCard.Stave == card.Stave && curCard.Value == card.Value {
			return true
		}
	}

	return false
}

// Shuffle the deck
func Shuffle(deck Deck) Deck {
	for i := 1; i < len(deck); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random int we generated then we'll switch them out
		if i != r {
			(deck)[r], (deck)[i] = (deck)[i], (deck)[r]
		}
	}

	return deck
}

// Deal a specified amount of cards
func Deal(deck Deck, n int) (Deck, Deck) {
	hand := Deck{}

	for i := 0; i < n; i++ {
		hand = append(hand, (deck)[i])
	}

	// Remove cards from the deck
	deck = (deck)[n:len(deck)]

	return deck, hand
}

func DealSingle(deck Deck) (Deck, Card) {
	deck, hand := Deal(deck, 1)
	return deck, hand[0]
}

// Debug helps debugging the deck of cards
func Debug(deck Deck) {
	spew.Dump(deck)
}
