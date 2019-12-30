package cards

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Card holds the card suits and types in the deck
type Card struct {
	Type string
	Suit string
}

// Deck holds the cards in the deck to be shuffled
type Deck []Card

// New creates a deck of cards to be used
func New() (deck Deck) {

	// Valid types include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	types := []string{"Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace"}
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
	return
}

// Shuffle the deck
func Shuffle(d Deck) Deck {
	for i := 1; i < len(d); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d[r], d[i] = d[i], d[r]
		}
	}
	return d
}

// Deal a specified amount of cards
func Deal(d Deck, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(d[i])
	}
}

// Debug helps debugging the deck of cards
func Debug(d Deck) {
	if os.Getenv("DEBUG") != "" {
		for i := 0; i < len(d); i++ {
			fmt.Printf("Card #%d is a %s of %ss\n", i+1, d[i].Type, d[i].Suit)
		}
	}
}

// Seed our randomness with the current time
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	deck := New()
	Debug(deck)
	Shuffle(deck)
	Debug(deck)
	Deal(deck, 2)
}
