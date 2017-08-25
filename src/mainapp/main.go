// Social Tournament Service
// HTTP-server implementation
package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	aux "mainapp/aux"

	endpoint "mainapp/endpoint"
)

func main() {
	fmt.Println("Social Tournament Service")

	cmdLine := strings.Join(os.Args[1:], " ")

	err := aux.InParam.ReadArgs(cmdLine)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("New port number: %s\n", aux.InParam.PortNumber)

	fmt.Printf("> Use http://localhost:%s/endpoint_name?par1=val1&par2=val2\n",
		aux.InParam.PortNumber)
	fmt.Printf("> Supported endpoints:\n  - /take\n  - /fund\n  - /balance\n  - /reset\n  - /joinTournament\n  - /announceTournament\n  - /resultTournament\n")

	http.HandleFunc("/take", endpoint.TakePoints)
	http.HandleFunc("/fund", endpoint.FundPlayer)
	http.HandleFunc("/balance", endpoint.ShowBalance)
	http.HandleFunc("/reset", endpoint.ClearDB)
	http.HandleFunc("/announceTournament", endpoint.AnnounceTournament)
	http.HandleFunc("/joinTournament", endpoint.JoinTournament)
	http.HandleFunc("/resultTournament", endpoint.ResultTournament)

	http.ListenAndServe(":"+aux.InParam.PortNumber, nil)
}
