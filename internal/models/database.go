package models

import "github.com/compycore/sabacc/internal/deck"

type Database struct {
	Codename    string    `json:"codename"`
	Round       int       `json:"round"`
	Turn        int       `json:"turn"`
	Dealer      int       `json:"dealer"`
	Rolled      bool      `json:"rolled"`
	Dice        []int     `json:"dice"`
	AllPlayers  []Player  `json:"players"`
	Draw        deck.Card `json:"draw"`
	AllDiscards deck.Deck `json:"discards"`
	Rematch     []Player  `json:"rematch"`
	Result      string    `json:"result"`
}

type Player struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Hand     deck.Deck `json:"hand"`
	HandSize int       `json:"handSize"`
	Score    int       `json:"score"`
}
