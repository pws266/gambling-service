// endpoint "/resultTournament" implementation
package endpoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aux "mainapp/aux"
	database "mainapp/database"
)

func ResultTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"resultTournament\" response")

	if r.Method != "POST" {
		http.Error(w, "Error: \"resultTournament\" should be POST request", http.StatusBadRequest)
		return
	}

	// reading request paramenets
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg := aux.CreateExternalError("ResultTournament", "Error while reading \"resultTournament\" - endpoint parameters", err)
		http.Error(w, errMsg.Error(), http.StatusBadRequest)
		return
	}

	// reading tournament results using JSON
	var tourRes aux.ResultTraits
	err = json.Unmarshal(body, &tourRes)

	if err != nil {
		errMsg := aux.CreateExternalError("ResultTournament", "Error while parsing \"resultTournament\" - endpoint parameters", err)
		http.Error(w, errMsg.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Results: %v\n", tourRes)

	statusCode, err := database.SetResult(tourRes.TournamentId, tourRes.Winners)

	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
}
