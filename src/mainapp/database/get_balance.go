package database

import (
	"database/sql"

	aux "mainapp/aux"
)

// returns player's balance, http status code and error description
func GetBalance(playerName string) (float64, int, error) {
	var balance float64

	db, err := getDB()

	if err != nil {
		return 0, 500, err
	}

	defer db.Close()

	// getting player's balance value using player's name
	err = db.QueryRow("SELECT balance FROM player WHERE name=?", playerName).Scan(&balance)

	switch {
	// player isn't found
	case err == sql.ErrNoRows:
		notFoundErr := aux.CreateError("ShowBalance", "No player was found with specified name")
		return 0, 404, notFoundErr
		// OK
	case err == nil:
		return balance, 200, err
		// unknown error during DB processing
	default:
		queryErr := aux.CreateExternalError("ShowBalance", "Unexpected error", err)
		return 0, 500, queryErr
	}
}
