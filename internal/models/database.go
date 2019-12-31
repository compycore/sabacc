package models

import "github.com/jessemillar/sabacc/internal/deck"

/*
JSON sample:

{
	"round": 1,
  "turn": 1,
  "players": [
    {
      "email": "hellojessemillar@gmail.com",
      "hand": [
        {
          "value": 1,
          "stave": "triangle"
        },
        {
          "value": -1,
          "stave": "circle"
        }
      ]
    },
    {
      "email": "fairdinkumriley@gmail.com",
      "hand": [
        {
          "value": 1,
          "stave": "triangle"
        },
        {
          "value": -1,
          "stave": "circle"
        }
      ]
    },
    {
      "email": "redscooterboy@gmail.com",
      "hand": [
        {
          "value": 1,
          "stave": "triangle"
        },
        {
          "value": -1,
          "stave": "circle"
        }
      ]
    }
  ],
  "discards": [
    {
      "value": 1,
      "stave": "triangle"
    },
    {
      "value": -1,
      "stave": "circle"
    }
  ]
}
*/

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
}
