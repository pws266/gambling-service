package database

import (
	aux "mainapp/aux"
)

// creates tournament with specified name and deposit
// returns http-status code and error description
func CreateTournament(tournamentId int, deposit float64) (int, error) {
	// getting DB descriptor, checking the conneection
	db, err := getDB()

	if err != nil {
		return 500, err
	}

	defer db.Close()

	// searching for tournament with specified name
	// getting player's balance value using player's name
	var rowsNumber int

	err = db.QueryRow("SELECT count(*) FROM tournament WHERE tournament_id=?", tournamentId).Scan(&rowsNumber)

	if err != nil {
		queryErr := aux.CreateExternalError("CreateTournament", "Unexpected error while searching tournament ID", err)
		return 500, queryErr
	}

	// tournament ID should be unique, otherwise return with error
	if rowsNumber != 0 {
		queryErr := aux.CreateError("CreateTournament", "The tournament with specified ID also present in DB")
		return 400, queryErr
	} else {
		// adding new tournament to "tournament" table
		statement, insertErr := db.Prepare("INSERT tournament SET tournament_id=?, deposit=?")

		if insertErr != nil {
			stmtErr := aux.CreateExternalError("CreateTournament", "Unable to prepare statement for adding new tournament", insertErr)
			return 500, stmtErr
		}

		_, insertErr = statement.Exec(tournamentId, deposit)

		if insertErr != nil {
			execErr := aux.CreateExternalError("CreateTournament", "Unable to add new tournament", insertErr)
			return 500, execErr
		}
	}

	return 200, nil
}
