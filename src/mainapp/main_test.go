package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	//	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// It's testing on real running server without stubs and mocks. Therefore the server should be worked on 3000 port.
func host() string {
	envHost := os.Getenv("HOST")
	if envHost == "" {
		envHost = "http://localhost:7090"
	}
	return envHost
}

type PlayerBalance struct {
	PlayerId string
	//Balance  string
	Balance float64
}

type Winner struct {
	PlayerId string  `json:"playerId,omitempty"`
	Prize    float64 `json:"prize,omitempty"`
}

type Tournament struct {
	TournamentId string   `json:"tournamentId,omitempty"`
	Winners      []Winner `json:"winners,omitempty"`
}

func cleanDB(t *testing.T) {
	getRequest(t, "/reset")
}

func getRequest(t *testing.T, uri string) (*http.Response, string) {
	response, err := http.Get(host() + uri)
	if err != nil {
		t.Fatal(err)
	}

	return response, getResponceBody(t, response)
}

func postRequest(t *testing.T, uri string, data interface{}) (*http.Response, string) {
	postJson, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	response, err := http.Post(host()+uri, "application/json", bytes.NewBuffer(postJson))
	if err != nil {
		t.Fatal(err)
	}

	return response, getResponceBody(t, response)
}

func getResponceBody(t *testing.T, response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func checkPlayerBalance(t *testing.T, playerId string, balance float64) {
	res, body := getRequest(t, "/balance?playerId="+playerId)
	player := parseJsonBodyPlayer(t, body)

	//	playerBalance, err := strconv.ParseFloat(player.Balance, 64)
	playerBalance := player.Balance
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, balance, playerBalance)
}

func parseJsonBodyPlayer(t *testing.T, body string) PlayerBalance {
	var data PlayerBalance
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestNotFoundPage(t *testing.T) {
	res, _ := getRequest(t, "/bla-bla-bla")
	assert.Equal(t, 404, res.StatusCode)
}

func TestUserBalance(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/balance?playerId=P1")
	assert.Equal(t, 404, res.StatusCode)
	getRequest(t, "/fund?playerId=P1&points=300")
	checkPlayerBalance(t, "P1", 300.00)
}

func TestCreatePlayer(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/balance?playerId=P1")
	assert.Equal(t, 404, res.StatusCode)

	res, _ = getRequest(t, "/fund?playerId=P1&points=1000")
	assert.Equal(t, 200, res.StatusCode)

	checkPlayerBalance(t, "P1", 1000.00)
}

func TestFundNegativePointsAndTransaction(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/fund?playerId=P1&points=-1000")
	assert.Equal(t, 400, res.StatusCode)

	res, _ = getRequest(t, "/balance?playerId=P1")
	assert.Equal(t, 404, res.StatusCode)
}

func TestTryToMakeNegativeBalance(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/take?playerId=P1&points=600")
	assert.Equal(t, 404, res.StatusCode)

	getRequest(t, "/fund?playerId=P1&points=1000")

	res, _ = getRequest(t, "/take?playerId=P1&points=600")
	assert.Equal(t, 200, res.StatusCode)

	res, _ = getRequest(t, "/take?playerId=P1&points=600")
	assert.Equal(t, 400, res.StatusCode)

	checkPlayerBalance(t, "P1", 400.00)
}

func TestAnnounceTournamentWithNegativeDepositAndTransactionChecking(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=10000")
	res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=-1000")
	assert.Equal(t, 400, res.StatusCode)

	res, _ = getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
	assert.Equal(t, 404, res.StatusCode)
}

func TestTryAnnounceTournamentTwoTimes(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
	assert.Equal(t, 200, res.StatusCode)

	res, _ = getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
	assert.Equal(t, 400, res.StatusCode)
}

func TestTournamentWithOnePlayer(t *testing.T) {
	cleanDB(t)
	res, _ := getRequest(t, "/balance?playerId=P1")
	assert.Equal(t, 404, res.StatusCode)

	getRequest(t, "/fund?playerId=P1&points=1500")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")

	checkPlayerBalance(t, "P1", 500.00)

	winner_1 := Winner{PlayerId: "P1", Prize: 1000}
	result := Tournament{TournamentId: "1", Winners: []Winner{winner_1}}

	res, _ = postRequest(t, "/resultTournament", result)

	assert.Equal(t, 200, res.StatusCode)

	checkPlayerBalance(t, "P1", 1500.00)
}

func TestTournamentWithTwoPlayers(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=1200")
	getRequest(t, "/fund?playerId=P2&points=500")
	getRequest(t, "/fund?playerId=P3&points=500")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P2&backerId=P3")

	winner := Winner{PlayerId: "P2", Prize: 900}
	result := Tournament{TournamentId: "1", Winners: []Winner{winner}}

	res, _ := postRequest(t, "/resultTournament", result)
	assert.Equal(t, 200, res.StatusCode)

	checkPlayerBalance(t, "P1", 200.00)
	checkPlayerBalance(t, "P2", 450)
	checkPlayerBalance(t, "P3", 450)
}

func TestTournamentWithTwoPlayersAndTwoWinners(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=300")
	getRequest(t, "/fund?playerId=P2&points=300")
	getRequest(t, "/fund?playerId=P3&points=300")
	getRequest(t, "/fund?playerId=P4&points=500")
	getRequest(t, "/fund?playerId=P5&points=1100")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=1001")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P5")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1&backerId=P2&backerId=P3&backerId=P4")

	winner_1 := Winner{PlayerId: "P1", Prize: 1000}
	winner_5 := Winner{PlayerId: "P5", Prize: 500}
	result := Tournament{TournamentId: "1", Winners: []Winner{winner_1, winner_5}}

	res, _ := postRequest(t, "/resultTournament", result)
	assert.Equal(t, 200, res.StatusCode)

	checkPlayerBalance(t, "P1", 299.75)
	checkPlayerBalance(t, "P4", 499.75)
	checkPlayerBalance(t, "P5", 599.00)
}

func TestAddOnePlayersAsBackerToTimesToOneTournament(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=2000")
	getRequest(t, "/fund?playerId=P2&points=2000")
	getRequest(t, "/fund?playerId=P3&points=2000")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=500")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1&backerId=P3")
	res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P2&backerId=P3")
	assert.Equal(t, 200, res.StatusCode)
}

func TestAddOnePlayersTwoTimesToOneTournament(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=3000")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=500")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
	res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
	assert.Equal(t, 400, res.StatusCode)
}

func TestOnManyThreadsIncreasePlayerBalance(t *testing.T) {
	cleanDB(t)

	const requestsCount = 10
	responses := make(chan int, requestsCount)

	var wg sync.WaitGroup
	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := getRequest(t, "/fund?playerId=P1&points=100")
			responses <- res.StatusCode
		}()
	}
	wg.Wait()

	successResponses := 0
	for i := 1; i <= requestsCount; i++ {
		if 200 == <-responses {
			successResponses += 1
		}
	}
	assert.Equal(t, requestsCount, successResponses)

	checkPlayerBalance(t, "P1", 1000.00)
}

func TestOnManyThreadsDecreasePlayerBalance(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=100")

	const requestsCount = 5
	responses := make(chan int, requestsCount)

	var wg sync.WaitGroup
	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := getRequest(t, "/take?playerId=P1&points=100")
			responses <- res.StatusCode

		}()
	}
	wg.Wait()

	successResponses := 0
	for i := 1; i <= requestsCount; i++ {
		if 200 == <-responses {
			successResponses += 1
		}
	}
	assert.Equal(t, 1, successResponses)

	checkPlayerBalance(t, "P1", 0.00)
}

func TestOnManyThreadsCreateTournament(t *testing.T) {
	cleanDB(t)

	const requestsCount = 5
	responses := make(chan int, requestsCount)

	var wg sync.WaitGroup
	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := getRequest(t, "/announceTournament?tournamentId=1&deposit=500")
			responses <- res.StatusCode

		}()
	}
	wg.Wait()

	successResponses := 0
	for i := 1; i <= requestsCount; i++ {
		if 200 == <-responses {
			successResponses += 1
		}
	}
	assert.Equal(t, 1, successResponses)
}

func TestOnManyThreadsJoinUserToTournament(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=1000")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=500")

	const requestsCount = 5
	responses := make(chan int, requestsCount)

	var wg sync.WaitGroup
	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
			responses <- res.StatusCode

		}()
	}
	wg.Wait()

	successResponses := 0
	for i := 1; i <= requestsCount; i++ {
		r := <-responses
		if 200 == r {
			successResponses += 1
		}
	}
	assert.Equal(t, 1, successResponses)

	checkPlayerBalance(t, "P1", 500.00)
}

func TestOnManyThreadsSubmitResultOfTournament(t *testing.T) {
	cleanDB(t)
	getRequest(t, "/fund?playerId=P1&points=2000")
	getRequest(t, "/fund?playerId=P2&points=2000")
	getRequest(t, "/fund?playerId=P3&points=1000")
	getRequest(t, "/announceTournament?tournamentId=1&deposit=1000")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P1")
	getRequest(t, "/joinTournament?tournamentId=1&playerId=P2&backerId=P3")

	winner_1 := Winner{PlayerId: "P1", Prize: 1000}
	winner_2 := Winner{PlayerId: "P2", Prize: 500}
	result := Tournament{TournamentId: "1", Winners: []Winner{winner_1, winner_2}}

	const requestsCount = 5
	responses := make(chan int, requestsCount)

	var wg sync.WaitGroup
	for i := 1; i <= requestsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, _ := postRequest(t, "/resultTournament", result)
			responses <- res.StatusCode
		}()
	}
	wg.Wait()

	successResponses := 0
	for i := 1; i <= requestsCount; i++ {
		r := <-responses
		if 200 == r {
			successResponses += 1
		}
	}
	assert.Equal(t, 1, successResponses)

	checkPlayerBalance(t, "P1", 2000.00)
	checkPlayerBalance(t, "P2", 1750.00)
	checkPlayerBalance(t, "P3", 750.00)
}
