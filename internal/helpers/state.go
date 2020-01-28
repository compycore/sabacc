package helpers

import "github.com/compycore/sabacc/internal/models"

func IsGameOver(database models.Database) bool {
	if database.Round <= 3 && database.Turn < len(database.AllPlayers) && len(database.AllPlayers) > 1 {
		return false
	}

	return true
}
