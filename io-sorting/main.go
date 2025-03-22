package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	database, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s: %v", dbFileName, err)
	}
	store, err := NewFileSystemPlayerStore(database)
	if err != nil {
		log.Fatalf("problem creating file system player store: %v", err)
	}
	server := NewPlayerServer(store)

	port := ":5000"
	fmt.Printf("Starting server on port %s...", port[1:])
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
