// endpoint "/announce" implementation
package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	database "mainapp/database"
)

func AnnounceTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"announce\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"announce\" should be GET request", http.StatusBadRequest)
		return
	}

	var (
		tournamentName string = r.FormValue("tournamentId")
		depositStr     string = r.FormValue("deposit")
	)

	if len(tournamentName) == 0 || len(depositStr) == 0 {
		http.Error(w, "Error: illegal \"announce\" endpoint parameters", http.StatusBadRequest)
		return
	}

	// checking deposit value
	deposit, err := strconv.ParseFloat(depositStr, 64)

	if err != nil || deposit <= 0 {
		http.Error(w, "Error: illegal \"deposit\" value", http.StatusBadRequest)
		return
	}

	// creating record in "tournament table"
	if statusCode, err := database.CreateTournament(tournamentName, deposit); err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	// displaying data in CURL request: "OK" response is adding automatically
	msg := fmt.Sprintf("  Added tournament: \n   TournamentID: %s\n   Deposit: %s\n", tournamentName, depositStr)

	w.Write([]byte(msg))
}
