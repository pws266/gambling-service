package database

import (
	"database/sql"

	aux "mainapp/aux"
)

// increases/reduces player's balance
// returns http-status code and error description
func ChangeBalance(playerName string, points float64, isIncrease bool) (int, error) {
	var player aux.Player

	// getting DB descriptor, checking the connection
	db, err := getDB()

	if err != nil {
		return 500, err
	}

	defer db.Close()

	// searching for player with specified name
	err = db.QueryRow("SELECT player_id, balance FROM player WHERE name=?", playerName).Scan(&player.ID, &player.Balance)

	switch {
	case err == sql.ErrNoRows:
		// can't reduce balance of non-existent player
		if !isIncrease {
			notFoundErr := aux.CreateError("ChangeBalance", "Player is not found")
			return 404, notFoundErr
		}

		// adding new player in the case of increasing balance of non-existent player
		statement, insertErr := db.Prepare("INSERT player SET name=?, balance=?")

		if insertErr != nil {
			stmtErr := aux.CreateExternalError("ChangeBalance", "Unable to prepare statement for adding new player", insertErr)
			return 500, stmtErr
		}

		_, insertErr = statement.Exec(playerName, points)

		if insertErr != nil {
			execErr := aux.CreateExternalError("ChangeBalance", "Unable to add new player", insertErr)
			return 500, execErr
		}

	case err == nil:
		var sign int = 1
		// can't reduce player's balance on sum greater than balance value
		if !isIncrease {
			if points > player.Balance {
				decreaseErr := aux.CreateError("ChangeBalance", "Unable to decrease balance. Illegal points value")
				return 400, decreaseErr
			}
			sign = -1
		}

		// updating player's balance
		player.Balance += float64(sign) * points
		statement, updateErr := db.Prepare("UPDATE player SET balance=? WHERE player_id=?")

		if updateErr != nil {
			stmtErr := aux.CreateExternalError("ChangeBalance", "Unable to prepare statement for update player's balance", updateErr)
			return 500, stmtErr
		}

		_, updateErr = statement.Exec(player.Balance, player.ID)

		if updateErr != nil {
			execErr := aux.CreateExternalError("ChangeBalance", "Unable to change player's balance", updateErr)
			return 500, execErr
		}

	default:
		queryErr := aux.CreateExternalError("ChangeBalance", "Error while looking for specified player", err)
		return 500, queryErr
	}

	return 200, nil
}
