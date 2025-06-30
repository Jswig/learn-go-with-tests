package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
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
	router.Handle("/game", http.HandlerFunc(server.handleGame))
	router.Handle("/ws", http.HandlerFunc(server.webSocket))
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

func (server *PlayerServer) handleGame(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("game.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("error parsing template: %v", err.Error()), http.StatusInternalServerError)
	}
	tmpl.Execute(writer, nil)
}

func (server *PlayerServer) webSocket(writer http.ResponseWriter, request *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("error creating websocket connection: %v", err.Error()), http.StatusInternalServerError)
	}
	messageType, content, err := conn.ReadMessage()
	if err != nil {
		http.Error(writer, fmt.Sprintf("error reading from websocket: %v", err.Error()), http.StatusInternalServerError)
	}
	if messageType == websocket.TextMessage {
		player := string(content)
		server.store.RecordWin(player)
	}
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
