package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Scores map[string]int

type JSONPlayerStore struct {
	fileName  string     // database file
	fileMutex sync.Mutex // mutex controlling access to the database file
}

func (store *JSONPlayerStore) RecordWin(name string) {
	store.fileMutex.Lock()
	defer store.fileMutex.Unlock()
	jsonContents, err := os.ReadFile(store.fileName)
	check(err)
	var scores Scores
	err = json.Unmarshal(jsonContents, &scores)
	check(err)
	_, isSet := scores[name]
	if isSet {
		scores[name]++
	} else {
		scores[name] = 1
	}
	jsonContents, err = json.Marshal(scores)
	check(err)
	err = os.WriteFile(store.fileName, jsonContents, 0644)
	check(err)
}

func (store *JSONPlayerStore) GetPlayerScore(name string) (score int, isSet bool) {
	store.fileMutex.Lock()
	defer store.fileMutex.Unlock()
	jsonContents, err := os.ReadFile(store.fileName)
	check(err)
	var scores Scores
	err = json.Unmarshal(jsonContents, &scores)
	check(err)
	score, isSet = scores[name]
	return score, isSet
}

func NewJSONPlayerStore(fileName string, scores Scores) *JSONPlayerStore {
	jsonContent, err := json.Marshal(scores)
	check(err)
	err = os.WriteFile(fileName, jsonContent, 0644)
	check(err)
	return &JSONPlayerStore{fileName, sync.Mutex{}}
}
