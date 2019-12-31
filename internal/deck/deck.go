package deck

import (
	"log"
	"math/rand"
)

// Card holds the card staves and values in the deck
type Card struct {
	Value int    `json:"value"`
	Stave string `json:"stave"`
}

// Deck holds the cards in the deck
type Deck []Card

// New creates a deck of cards to be used
func New() (deck Deck) {
	allValues := []int{-10, -9, -8, -7, -6, -5, -4, -3, -2, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	allStaves := []string{"circle", "square", "triangle"}

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

	return
}

// Remove a specific card from the deck
func (d *Deck) Remove(card Card) {
	for i, curCard := range *d {
		if curCard.Stave == card.Stave && curCard.Value == card.Value {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return
		}
	}
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
func (d *Deck) Deal(n int) Deck {
	hand := Deck{}

	for i := 0; i < n; i++ {
		hand = append(hand, (*d)[i])
	}

	// Remove cards from the deck
	*d = (*d)[n:len(*d)]

	return hand
}

// Debug helps debugging the deck of cards
func (d *Deck) Debug() {
	for i := 0; i < len(*d); i++ {
		log.Println(i+1, (*d)[i].Value, (*d)[i].Stave)
	}
}
