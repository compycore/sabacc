package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/jessemillar/sabacc/internal/cards"
	"github.com/jessemillar/sabacc/internal/email"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.Println("Making deck")
	deck := cards.New()
	deck.Debug()
	log.Println("Shuffling deck")
	deck.Shuffle()
	deck.Debug()

	err := email.Send()
	if err != nil {
		log.Println(err)
	}
}
