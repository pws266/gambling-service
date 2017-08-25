// endpoint "/balance" implementation
package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "mainapp/database"
)

type PlayerTraits struct {
	PlayerID string  `json:"playerId"`
	Balance  float64 `json:"balance"`
}

func ShowBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"balance\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"balance\" should be GET request", http.StatusBadRequest)
		return
	}

	if len(r.FormValue("playerId")) == 0 {
		http.Error(w, "Error: illegal \"balance\" argument key", http.StatusBadRequest)
		return
	}

	var player PlayerTraits = PlayerTraits{r.FormValue("playerId"), 0}

	balance, statusCode, err := database.GetBalance(r.FormValue("playerId"))

	if err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	// creating JSON response
	player.Balance = balance

	js, err := json.Marshal(player)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//responseStr := fmt.Sprintf("%s\n", js)

	w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(responseStr))
	w.Write([]byte(js))
}
