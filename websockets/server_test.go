package poker

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestGETPlayers(t *testing.T) {
	scores := map[string]int{
		"Pepper": 20,
		"Floyd":  10,
	}
	winCalls := make([]string, 0)
	store := &StubPlayerStore{scores, winCalls}
	server := mustMakePlayerServer(t, store)
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
	store := &StubPlayerStore{
		map[string]int{},
		[]string{},
	}
	server := mustMakePlayerServer(t, store)

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
	store := &StubPlayerStore{scores, []string{}}
	server := mustMakePlayerServer(t, store)

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

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{})
		request := newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response, http.StatusOK)
	})
	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Ruth"
		playerServer := mustMakePlayerServer(t, store)
		server := httptest.NewServer(playerServer)
		defer server.Close()

		socketURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		socket, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", socketURL, err)
		}
		defer socket.Close()

		if err := socket.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
			t.Fatalf("could not send message over websocket connection %v", err)
		}

		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, store, winner)
	})

}

func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
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

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game", nil)
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
