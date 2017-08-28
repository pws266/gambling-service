// endpoint "/join" implementation
// format: "srv_ip:port/join?playerId=P1&backerId=B1&backerId=B2 ... "
package endpoint

import (
	"fmt"
	"net/http"
	//	"strconv"
	"strings"

	database "mainapp/database"
)

func JoinTournament(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Endpoint \"join\" response")

	if r.Method != "GET" {
		http.Error(w, "Error: \"join\" should be GET request", http.StatusBadRequest)
		return
	}

	err := r.ParseForm()
	if err != nil {
		errMsg := fmt.Sprintf("Error: unable to parse request arguments in \"join\"\nDescription: %s", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	var (
		tournamentName []string = r.Form["tournamentId"]

		players []string = r.Form["playerId"]
		backers []string = r.Form["backerId"]
	)

	if len(tournamentName[0]) == 0 || len(players[0]) == 0 {
		http.Error(w, "Error: illegal parameters set in \"join\" endpoint", http.StatusBadRequest)
		return
	}

	var (
	//		tournamentId int
	)
	/*
		// getting tournament ID
		tournamentId, err = strconv.Atoi(r.Form["tournamentId"][0])

		if err != nil || tournamentId <= 0 {
			http.Error(w, "Error: illegal \"tournamentId\" value", http.StatusBadRequest)
			return
		}
	*/
	if len(backers) > 0 {
		for _, backer := range backers {
			if strings.Compare(backer, players[0]) == 0 {
				http.Error(w, "Error: one of \"backerId\" matched with \"playerId\" value", http.StatusBadRequest)
				return
			}
		}
	}

	deposit, fee, statusCode, err := database.SetupTournament(tournamentName[0], players[0], backers)
	if err != nil {
		errMsg := fmt.Sprintf("Error in DB processing: %s", err)
		http.Error(w, errMsg, statusCode)
		return
	}

	fmt.Printf("%f %f\n", deposit, fee)

	for i, player := range players {
		fmt.Printf("Player #%d: %s\n", i, player)
	}

	fmt.Println()

	for i, backer := range r.Form["backerId"] {
		fmt.Printf("Backer #%d: %s\n", i, backer)
	}

	//	if len(r.FormValue("playerId")) == 0 {
	//		http.Error(w, "Error: \"join\" endpoint should contain player", http.StatusBadRequest)
	//		return
	//	}
	/*
		points, err := strconv.Atoi(r.FormValue("points"))

		if err != nil || points < 0 {
			http.Error(w, "Error: illegal \"points\" value", http.StatusBadRequest)
			return
		}

		err = database.ChangeBalance(r.FormValue("playerId"), points, true)

		if err != nil {
			errMsg := fmt.Sprintf("Error in DB processing: %s", err)
			http.Error(w, errMsg, http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)

		msg := fmt.Sprintf("  Method: %s\n  playerId: %s\n  added: %s points\n\n",
			r.Method, r.FormValue("playerId"), r.FormValue("points"))
		w.Write([]byte(msg))
	*/
}
