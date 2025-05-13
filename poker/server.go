package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type Player struct {
	Name string
	Wins int
}

type PlayerStore interface {
	GetLeague() []Player
	GetPlayerScore(name string) (score int, scoreAvailable bool)
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(server.handlePlayers))
	router.Handle("/league", http.HandlerFunc(server.handleLeague))
	server.Handler = router
	return server
}

func (server *PlayerServer) handlePlayers(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.Path, "/players/")
	switch request.Method {
	case http.MethodPost:
		server.processWin(writer, player)
	case http.MethodGet:
		server.showScore(writer, player)
	}
}

func (server *PlayerServer) handleLeague(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", jsonContentType)
	players := server.store.GetLeague()
	json.NewEncoder(writer).Encode(players)
}

func (server *PlayerServer) showScore(writer http.ResponseWriter, player string) {
	score, found := server.store.GetPlayerScore(player)
	if found {
		fmt.Fprint(writer, score)
		writer.WriteHeader(http.StatusOK)

	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func (server *PlayerServer) processWin(writer http.ResponseWriter, player string) {
	server.store.RecordWin(player)
	writer.WriteHeader(http.StatusAccepted)
}
