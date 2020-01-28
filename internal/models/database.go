package models

import "github.com/compycore/sabacc/internal/deck"

type Database struct {
	Password    string    `json:"password"`
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
	Template    Template  `json:",omitempty"`
}

type Player struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Hand     deck.Deck `json:"hand"`
	HandSize int       `json:"handSize"`
	Score    int       `json:"score"`
}

// These values are only for use in templates, not for returning to the user
type Template struct {
	Link           string `json:",omitempty"`
	AllPlayerNames string `json:",omitempty"`
}
