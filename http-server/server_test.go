package main

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
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
	server := &PlayerServer{store}
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
	server := &PlayerServer{&store}

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

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return request
}

func assertStatusCode(t testing.TB, response *httptest.ResponseRecorder, want int) {
	got := response.Code
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, response *httptest.ResponseRecorder, want string) {
	got := response.Body.String()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
