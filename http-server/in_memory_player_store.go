package main

import "sync"

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		map[string]int{},
		sync.Mutex{},
	}
}

type InMemoryPlayerStore struct {
	scores map[string]int
	mutex  sync.Mutex
}

func (store *InMemoryPlayerStore) RecordWin(name string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.scores[name]++
}

func (store *InMemoryPlayerStore) GetPlayerScore(name string) (score int, isSet bool) {
	score, isSet = store.scores[name]
	return score, isSet
}
