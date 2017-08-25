package database

import (
	"database/sql"
	"fmt"

	aux "mainapp/aux"
)

// verifies if player presents in table "player"
func checkPlayer(db *sql.DB, playerName string, fee float64, player *aux.Player) (int, error) {
	// searching for specified player
	err := db.QueryRow("SELECT player_id, balance FROM player WHERE name = ?", playerName).Scan(&player.ID, &player.Balance)

	if err != nil {
		queryErr := aux.CreateExternalError("checkPlayer", "Can't find player with specified ID", err)
		return 404, queryErr
	}

	if player.Balance < fee {
		errMsg := fmt.Sprintf("Player \"%s\" balance %f is lower than tournament participation fee %f\n", playerName, player.Balance, fee)
		feeErr := aux.CreateError("checkPlayer", errMsg)

		return 500, feeErr
	}

	return 200, nil
}

// reduce balance of specified player on tournament participation fee
func reduceBalance(db *sql.DB, player aux.Player, fee float64) (int, error) {
	player.Balance -= fee

	statement, err := db.Prepare("UPDATE player SET balance=? WHERE player_id=?")

	if err != nil {
		stmtErr := aux.CreateExternalError("reduceBalance", "Unable to prepare statement for decreasing participant/backer balance", err)
		return 500, stmtErr
	}

	_, err = statement.Exec(player.Balance, player.ID)

	if err != nil {
		insErr := aux.CreateExternalError("reduceBalance", "Unable to decrease participant/backer balance", err)
		return 500, insErr
	}

	return 200, nil
}

var (
	sqlRequest []string = []string{
		"INSERT participant SET tournament_id=?, player_id=?, fee=?",
		"INSERT backer SET tournament_id=?, player_id=?, backer_id=?",
	}
)

// insert player and backers into appropriate tables
func insertInTables(db *sql.DB, tournamentId int, player []aux.Player, fee float64) (int, error) {
	var currentRequest string

	for i := range player {
		if i == 0 {
			currentRequest = sqlRequest[0]
		} else {
			currentRequest = sqlRequest[1]
		}

		// preparing statement for adding player/backer to appropriate table
		statement, err := db.Prepare(currentRequest)

		if err != nil {
			stmtErr := aux.CreateExternalError("insertInTables", "Unable to prepare statement for adding new tournament participant/backer", err)
			return 500, stmtErr
		}

		// adding participant/backer to tournament
		if i == 0 {
			_, err = statement.Exec(tournamentId, player[i].ID, fee)
		} else {
			_, err = statement.Exec(tournamentId, player[0].ID, player[i].ID)
		}

		if err != nil {
			insErr := aux.CreateExternalError("insertInTables", "Unable to add new tournament participant/backer", err)
			return 500, insErr
		}
	}

	return 200, nil
}

// assigns players and backers for specified tournament
// returns tournament deposit, backing sum, http-status code and error description
func SetupTournament(tournamentId int, playerId string, backerId []string) (float64, float64, int, error) {
	// getting DB descriptor, checking the conneection
	db, err := getDB()

	if err != nil {
		return 0, 0, 500, err
	}

	defer db.Close()

	// getting deposit value of specified tournament
	var (
		deposit float64
		fee     float64

		playerNames []string     = make([]string, len(backerId)+1)
		players     []aux.Player = make([]aux.Player, len(playerNames))
	)

	// initializing player/backer names array
	playerNames[0] = playerId
	copy(playerNames[1:], backerId)

	// getting tournament deposit
	err = db.QueryRow("SELECT deposit FROM tournament WHERE tournament_id = ?", tournamentId).Scan(&deposit)

	if err != nil {
		queryErr := aux.CreateExternalError("SetupTournament", "Can't find tournament with specified ID", err)
		return 0, 0, 404, queryErr
	}

	// getting fee for one player/backer
	fee = deposit / float64(len(backerId)+1)

	// checking presence of player/backers in database
	for i, playerName := range playerNames {
		statusCode, err := checkPlayer(db, playerName, fee, &players[i])

		if err != nil {
			return deposit, fee, statusCode, err
		}
	}

	// adding participant and backers to appropriate tables
	statusCode, err := insertInTables(db, tournamentId, players, fee)

	if err != nil {
		return deposit, fee, statusCode, err
	}

	// decreasing balance of tournament participants and backers
	for _, currentPlayer := range players {
		statusCode, err := reduceBalance(db, currentPlayer, fee)

		if err != nil {
			return deposit, fee, statusCode, err
		}
	}

	return deposit, fee, 200, nil
}
