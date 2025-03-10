package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	port := ":5000"
	fmt.Printf("Starting server on port %s...", port[1:])
	log.Fatal(http.ListenAndServe(":5000", server))
}
