package main

import (
	"fmt"
	"log"
	"net/http"

	"poker"
)

const dbFileName = "game.db.json"

func main() {
	store, cleanup, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	server := poker.NewPlayerServer(store)

	port := ":5000"
	fmt.Printf("Starting server on port %s...", port[1:])
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
