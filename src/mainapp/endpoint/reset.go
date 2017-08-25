// endpoint "/reset" implementation
package endpoint

import (
	"fmt"
	"net/http"

	database "mainapp/database"
)

func ClearDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"reset\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"reset\" should be GET request", http.StatusBadRequest)
		return
	}

	if statusCode, err := database.ClearTables(); err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	var msg string = "All data are removed!\n"
	w.Write([]byte(msg))
}
