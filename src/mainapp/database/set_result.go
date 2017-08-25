package database

import (
	"container/list"
	"database/sql"
	"fmt"

	aux "mainapp/aux"
)

// searches for tournament with specified ID
// returns status code and error description
func searchTournament(db *sql.DB, tournamentId uint) (int, error) {
	var tourNumber int

	err := db.QueryRow("SELECT count(*) FROM tournament WHERE tournament_id=?", tournamentId).Scan(&tourNumber)

	if err != nil {
		queryErr := aux.CreateExternalError("searchTournament", "Unexpected error while searching tournament ID", err)
		return 500, queryErr
	}

	// generating error if tournament with specified ID wasn't found
	if tourNumber == 0 {
		queryErr := aux.CreateError("searchTournament", "The tournament with specified ID isn't found in DB")
		return 404, queryErr
	}

	return 200, nil
}

// reads specified tournament player and backers IDs and balance
// returns list of players/backers, http-status code and error description
func getPlayerAndBackers(db *sql.DB, tournamentId uint, winner aux.Winner, players *list.List) (int, error) {
	// players and backers
	var currentPlayer aux.Player

	// looking for specified winner ID in "player" table
	err := db.QueryRow("SELECT player_id, balance FROM player WHERE name=?", winner.PlayerId).Scan(&currentPlayer.ID, &currentPlayer.Balance)

	if err != nil {
		var (
			errMsg     string = "Unexpected error while searching specified winner"
			statusCode int    = 500
		)

		if err == sql.ErrNoRows {
			errMsg = "Specified winner is not found"
			statusCode = 404
		}

		queryErr := aux.CreateExternalError("getPlayerAndBackers", errMsg, err)
		return statusCode, queryErr
	}

	// looking for winner as specified tournament participant
	fmt.Printf("Tournament ID: %d Player ID: %d\n", tournamentId, currentPlayer.ID)

	var playerNumber int
	err = db.QueryRow("SELECT count(*) FROM participant WHERE (tournament_id=? AND player_id=?)", tournamentId, currentPlayer.ID).Scan(&playerNumber)

	// generating error if tournament with specified ID wasn't found
	if playerNumber == 0 {
		queryErr := aux.CreateError("getPlayerAndBackers", "The player with specified ID isn't participant of specified tournament")
		return 404, queryErr
	}

	// adding player to resulting list
	players.PushBack(currentPlayer)

	// looking for backers
	rows, err := db.Query("SELECT backer_id FROM backer WHERE (tournament_id=? AND player_id=?)", tournamentId, currentPlayer.ID)

	if err != nil {
		queryErr := aux.CreateExternalError("getPlayerAndBackers", "Problems while searching for backers", err)
		return 500, queryErr
	}

	defer rows.Close()

	// saving backers to the player list
	for rows.Next() {
		// getting backer ID
		if err = rows.Scan(&currentPlayer.ID); err != nil {
			readErr := aux.CreateExternalError("getPlayerAndBackers", "Unable to get backer ID", err)
			return 500, readErr
		}

		// getting backer balance
		err = db.QueryRow("SELECT balance FROM player WHERE player_id=?", currentPlayer.ID).Scan(&currentPlayer.Balance)

		if err != nil {
			queryErr := aux.CreateExternalError("getPlayerAndBackers", "Unable to get backer balance", err)
			return 500, queryErr
		}

		players.PushBack(currentPlayer)
	}

	if err = rows.Err(); err != nil {
		rowErr := aux.CreateExternalError("getPlayerAndBackers", "Problem while processing query rows for backers", err)
		return 500, rowErr
	}

	fmt.Println("--> Player and backers list: ")
	for elm := players.Front(); elm != nil; elm = elm.Next() {
		fmt.Printf("ID: %3d Balance %8.2f\n", elm.Value.(aux.Player).ID, elm.Value.(aux.Player).Balance)
	}

	return 200, nil
}

// updates winners (player and backers) balance according to the prize value
// returns http-status code and error description
func updateWinnersBalance(db *sql.DB, players *list.List, prize float64) (int, error) {
	var payout float64
	var newBalance float64

	if players.Len() == 0 {
		listError := aux.CreateError("updateWinnerBalance", "Empty list of player/backers")
		return 500, listError
	}

	// getting payout for each backer
	payout = prize / float64(players.Len())

	// updating each winner (player or backer) balance
	for elm := players.Front(); elm != nil; elm = elm.Next() {
		newBalance = elm.Value.(aux.Player).Balance + payout

		statement, updateErr := db.Prepare("UPDATE player SET balance=? WHERE player_id=?")

		if updateErr != nil {
			stmtErr := aux.CreateExternalError("updateWinnersBalance", "Unable to prepare statement for update player's balance", updateErr)
			return 500, stmtErr
		}

		_, updateErr = statement.Exec(newBalance, elm.Value.(aux.Player).ID)

		if updateErr != nil {
			execErr := aux.CreateExternalError("updateWinnersBalance", "Unable to change player's balance", updateErr)
			return 500, execErr
		}
	}

	return 200, nil
}

// returns status code and error description
func SetResult(tournamentId uint, winners []aux.Winner) (int, error) {
	// getting DB descriptor, checking the conneection
	db, err := getDB()

	if err != nil {
		return 500, err
	}

	defer db.Close()

	// searching for tournament with specified ID
	if statusCode, err := searchTournament(db, tournamentId); err != nil {
		return statusCode, err
	}

	var playerList list.List

	for _, currentWinner := range winners {
		playerList.Init()

		// getting player/backers list
		if statusCode, err := getPlayerAndBackers(db, tournamentId, winners[0], &playerList); err != nil {
			return statusCode, err
		}

		// updating current winner balance
		if statusCode, err := updateWinnersBalance(db, &playerList, currentWinner.Prize); err != nil {
			return statusCode, err
		}
	}

	return 200, nil
}
