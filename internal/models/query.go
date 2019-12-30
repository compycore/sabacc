package models

type Query struct {
	Round       int      `json:"round"`
	Turn        int      `json:"turn"`
	AllPlayers  []string `json:"players"`
	AllHands    []string `json:"hands"`
	AllDiscards []string `json:"discards"`
}
