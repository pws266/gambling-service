// endpoint "/fund" implementation
package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

	database "mainapp/database"
)

func FundPlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"fund\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"fund\" should be GET request", http.StatusBadRequest)
		return
	}

	if len(r.FormValue("playerId")) == 0 || len(r.FormValue("points")) == 0 {
		http.Error(w, "Error: illegal \"fund\" argument key", http.StatusBadRequest)
		return
	}

	points, err := strconv.ParseFloat(r.FormValue("points"), 64)

	if err != nil || points < 0 {
		http.Error(w, "Error: illegal \"points\" value", http.StatusBadRequest)
		return
	}

	statusCode, err := database.ChangeBalance(r.FormValue("playerId"), points, true)

	if err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)

		return
	}

	msg := fmt.Sprintf("  Method: %s\n  playerId: %s\n  added: %s points\n\n",
		r.Method, r.FormValue("playerId"), r.FormValue("points"))
	w.Write([]byte(msg))
}

func main() {
	fmt.Println("Hello World!")
}
