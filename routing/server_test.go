package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (store *StubPlayerStore) GetLeague() []Player {
	players := make([]Player, 0, len(store.scores))
	for name, score := range store.scores {
		players = append(players, Player{name, score})
	}
	return players
}

func (store *StubPlayerStore) GetPlayerScore(name string) (score int, isSet bool) {
	score, isSet = store.scores[name]
	return score, isSet
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.winCalls = append(store.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	scores := map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}
	winCalls := make([]string, 0)
	store := &StubPlayerStore{scores, winCalls}
	server := NewPlayerServer(store)
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertResponseBody(t, response, "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertResponseBody(t, response, "10")
	})

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := newGetScoreRequest("Appolo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
	}
	server := NewPlayerServer(&store)

	t.Run("it records win on POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusAccepted)

		got := store.winCalls
		want := []string{"Pepper"}
		if !slices.Equal(got, want) {
			t.Errorf("Did not record expected wins, got %v want %v", got, want)
		}
	})
}

func TestLeague(t *testing.T) {
	scores := map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}
	store := StubPlayerStore{scores, []string{}}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on /league", func(t *testing.T) {
		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		got := getLeagueFromResponse(t, response)
		want := []Player{
			{"Pepper", 20},
			{"Floyd", 10},
		}
		assertLeague(t, got, want)
	})
}

func getLeagueFromResponse(t testing.TB, response *httptest.ResponseRecorder) (league []Player) {
	t.Helper()
	err := json.NewDecoder(response.Body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
	}
	return
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return request
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got []Player, want []Player) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertResponseBody(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	got := response.Body.String()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertStatusCode(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
