package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   []Player
}

func initializePlayerDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf(
			"problem getting file info from file %s, %v",
			file.Name(),
			err,
		)
	}
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}
	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initializePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf(
			"problem initializing players database file %s, %v",
			file.Name(),
			err,
		)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf(
			"problem loading player store from file %s, %v",
			file.Name(),
			err,
		)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

// league.go
func NewLeague(reader io.Reader) (league []Player, err error) {
	err = json.NewDecoder(reader).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return league, err
}

func (store *FileSystemPlayerStore) GetLeague() (league []Player) {
	slices.SortFunc(store.league, func(p1, p2 Player) int {
		if p1.Wins > p2.Wins {
			return -1
		} else if p1.Wins < p2.Wins {
			return 1
		} else {
			return 0
		}
	})
	return store.league
}

func (store *FileSystemPlayerStore) GetPlayerScore(name string) (score int, found bool) {
	player, found := findPlayer(store.league, name)
	if player != nil {
		score = player.Wins
	}
	return
}

func findPlayer(league []Player, name string) (player *Player, found bool) {
	for i, player := range league {
		if player.Name == name {
			found = true
			return &league[i], found
		}
	}
	return
}

func (store *FileSystemPlayerStore) RecordWin(name string) {
	league := store.league
	player, found := findPlayer(league, name)
	if found {
		player.Wins++
	} else {
		store.league = append(league, Player{name, 1})
	}
	store.database.Encode(league)
}
