package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (store *StubPlayerStore) GetLeague() []Player {
	players := make([]Player, 0, len(store.scores))
	for name, score := range store.scores {
		players = append(players, Player{name, score})
	}
	return players
}

func (store *StubPlayerStore) GetPlayerScore(name string) (score int, isSet bool) {
	score, isSet = store.scores[name]
	return score, isSet
}

func (store *StubPlayerStore) RecordWin(name string) {
	store.winCalls = append(store.winCalls, name)
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}
