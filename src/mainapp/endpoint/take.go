// endpoint "/take" implementation
package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	database "mainapp/database"
)

func TakePoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"take\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"take\" should be GET request", http.StatusBadRequest)
		return
	}

	if len(r.FormValue("playerId")) == 0 || len(r.FormValue("points")) == 0 {
		http.Error(w, "Error: illegal \"take\" argument key", http.StatusBadRequest)
		return
	}

	points, err := strconv.ParseFloat(r.FormValue("points"), 64)

	if err != nil || points < 0 {
		http.Error(w, "Error: illegal \"points\" value", http.StatusBadRequest)
		return
	}

	statusCode, err := database.ChangeBalance(r.FormValue("playerId"), points, false)

	if err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	msg := fmt.Sprintf("  Method: %s\n  playerId: %s\n  taken: %s points\n\n",
		r.Method, r.FormValue("playerId"), r.FormValue("points"))
	w.Write([]byte(msg))
}
