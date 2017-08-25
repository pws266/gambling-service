// useful structures description
package aux

// player description for search purposes
type Player struct {
	ID      uint
	Balance float64
}

// single winner description for JSON
type Winner struct {
	PlayerId string  `json:"playerId"`
	Prize    float64 `json:"prize"`
}

// tournament results description for reading via JSON
type ResultTraits struct {
	TournamentId uint     `json:"tournamentId"`
	Winners      []Winner `json:"winners"`
}
