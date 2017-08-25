package database

import (
	aux "mainapp/aux"
)

var (
	clearTablesSQLCommands []string = []string{
		"SET FOREIGN_KEY_CHECKS = 0",
		"SET AUTOCOMMIT = 0",
		"TRUNCATE TABLE player",
		"TRUNCATE TABLE tournament",
		"TRUNCATE TABLE participant",
		"TRUNCATE TABLE backer",
		"SET FOREIGN_KEY_CHECKS = 1",
		"COMMIT",
		"SET AUTOCOMMIT = 1",
	}

	clearTablesErrMsg []string = []string{
		"Unable assign FOREIGN_KEY_CHECKS flag to 0",
		"Error while disabling AUTOCOMMIT",
		"Unexpected error while deleting data in \"player\" table",
		"Unexpected error while deleting data in \"tournament\" table",
		"Unexpected error while deleting data in \"participant\" table",
		"Unexpected error while deleting data in \"backer\" table",
		"Unable assign FOREIGN_KEY_CHECKS flag to 1",
		"Unable to commit transaction",
		"Error while enabling AUTOCOMMIT",
	}
)

// clears all DB
// returns http status code and error description
func ClearTables() (int, error) {
	db, err := getDB()

	if err != nil {
		return 500, err
	}

	defer db.Close()

	// opening transaction
	tx, err := db.Begin()
	if err != nil {
		clearErr := aux.CreateExternalError("ClearTables", "Error while creating transaction", err)
		return 500, clearErr
	}

	defer tx.Rollback()

	for i, currentExp := range clearTablesSQLCommands {
		if _, err := tx.Exec(currentExp); err != nil {
			clearErr := aux.CreateExternalError("ClearTables", clearTablesErrMsg[i], err)
			return 500, clearErr
		}
	}

	return 200, nil
}
