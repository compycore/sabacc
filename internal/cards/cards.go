package cards

import (
	"fmt"
	"log"
	"math/rand"
)

// Card holds the card suits and types in the deck
type Card struct {
	Type int
	Suit string
}

// Deck holds the cards in the deck to be shuffled
type Deck []Card

// New creates a deck of cards to be used
func New() (deck Deck) {
	types := []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	suits := []string{"Circle", "Square", "Triangle"}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(types); i++ {
		for n := 0; n < len(suits); n++ {
			card := Card{
				Type: types[i],
				Suit: suits[n],
			}
			deck = append(deck, card)
		}
	}

	// Add the zero cards to the deck
	for n := 0; n < 2; n++ {
		card := Card{
			Type: 0,
			Suit: "Zero",
		}

		deck = append(deck, card)
	}

	return
}

// Shuffle the deck
func (d *Deck) Shuffle() {
	for i := 1; i < len(*d); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			(*d)[r], (*d)[i] = (*d)[i], (*d)[r]
		}
	}
}

// Deal a specified amount of cards
func Deal(d Deck, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(d[i])
	}
}

// Debug helps debugging the deck of cards
func (d *Deck) Debug() {
	for i := 0; i < len(*d); i++ {
		log.Println(i+1, (*d)[i].Type, (*d)[i].Suit)
	}
}
