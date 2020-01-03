package models

import "github.com/jessemillar/sabacc/internal/deck"

type Database struct {
	Round       int       `json:"round"`
	Turn        int       `json:"turn"`
	AllPlayers  []Player  `json:"players"`
	Draw        deck.Card `json:"draw"`
	AllDiscards deck.Deck `json:"discards"`
}

type Player struct {
	Email string    `json:"email"`
	Hand  deck.Deck `json:"hand"`
	Score int       `json:"score"`
}
