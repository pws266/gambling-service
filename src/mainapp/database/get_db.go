package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	aux "mainapp/aux"
)

// obtains database handler and doesn't open connection
func getDB() (*sql.DB, error) {
	// forming DSN (data source name) string
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8", aux.InParam.UserLogin, aux.InParam.UserPassword, aux.InParam.DBName)

	pdb, err := sql.Open("mysql", dsn)

	if err != nil {
		openErr := aux.CreateExternalError("GetHandle", "Unable to get database handler", err)
		return pdb, openErr
	}

	// validating connection
	err = pdb.Ping()

	if err != nil {
		pingErr := aux.CreateExternalError("GetHandle", "Connection test is failed. Check your DSN", err)
		return pdb, pingErr
	}

	return pdb, nil
}
