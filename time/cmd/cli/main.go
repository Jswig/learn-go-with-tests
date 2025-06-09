package main

import (
	"fmt"
	"log"
	"os"
	"poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, cleanup, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	alerter := poker.BlindAlerterFunc(poker.StdOutAlerter)
	game := poker.NewCLI(os.Stdin, os.Stdout, poker.NewTexasHoldEm(store, alerter))
	game.PlayPoker()
}
