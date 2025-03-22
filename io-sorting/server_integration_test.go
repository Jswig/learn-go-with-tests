// server_integration_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanup := makeTestFile(t, "")
	defer cleanup()
	store, err := NewFileSystemPlayerStore(database)
	if err != nil {
		t.Errorf("Problem creating filesystem player store, %v", err)
	}
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatusCode(t, response, http.StatusOK)

		assertResponseBody(t, response, "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatusCode(t, response, http.StatusOK)
		assertContentType(t, response, jsonContentType)

		want := []Player{{"Pepper", 3}}
		var got []Player
		json.NewDecoder(response.Body).Decode(&got)
		assertLeague(t, got, want)
	})
}
