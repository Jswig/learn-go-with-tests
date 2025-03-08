package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(name string) (score int, scoreAvailable bool)
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
}

func (server *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.Path, "/players/")
	switch request.Method {
	case http.MethodPost:
		server.processWin(writer, player)
	case http.MethodGet:
		server.showScore(writer, player)
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
