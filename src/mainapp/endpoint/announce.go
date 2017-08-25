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

	if len(r.FormValue("tournamentId")) == 0 || len(r.FormValue("deposit")) == 0 {
		http.Error(w, "Error: illegal \"announce\" endpoint parameters", http.StatusBadRequest)
		return
	}

	// checking tournament ID value
	tournament_id, err := strconv.Atoi(r.FormValue("tournamentId"))
	if err != nil || tournament_id <= 0 {
		http.Error(w, "Error: illegal \"tournamentId\" value", http.StatusBadRequest)
		return
	}

	// checking deposit value
	deposit, err := strconv.ParseFloat(r.FormValue("deposit"), 64)

	if err != nil || deposit <= 0 {
		http.Error(w, "Error: illegal \"deposit\" value", http.StatusBadRequest)
		return
	}

	// creating record in "tournament table"
	if statusCode, err := database.CreateTournament(tournament_id, deposit); err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	// displaying data in CURL request: "OK" response is adding automatically
	msg := fmt.Sprintf("  Added tournament: \n   TournamentID: %s\n   Deposit: %s\n",
		r.FormValue("tournamentId"), r.FormValue("deposit"))

	w.Write([]byte(msg))
}
