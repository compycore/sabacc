package models

import "github.com/jessemillar/sabacc/internal/deck"

type Database struct {
	Codename    string    `json:"codename"`
	Round       int       `json:"round"`
	Turn        int       `json:"turn"`
	Dealer      int       `json:"dealer"`
	AllPlayers  []Player  `json:"players"`
	Draw        deck.Card `json:"draw"`
	AllDiscards deck.Deck `json:"discards"`
	Rematch     []Player  `json:"rematch"`
	Result      string    `json:"result"`
}

type Player struct {
	Email string    `json:"email"`
	Hand  deck.Deck `json:"hand"`
	Score int       `json:"score"`
}
