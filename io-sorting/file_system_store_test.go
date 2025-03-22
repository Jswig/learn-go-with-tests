package main

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		contents := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`
		database, cleanup := makeTestFile(t, contents)
		defer cleanup()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		got = store.GetLeague()

		assertLeague(t, got, want)
	})

	//file_system_store_test.go
	t.Run("get player score", func(t *testing.T) {
		contents := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`
		database, cleanup := makeTestFile(t, contents)
		defer cleanup()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got, _ := store.GetPlayerScore("Chris")
		want := 33

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("record player win for existing player", func(t *testing.T) {
		contents := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`
		database, cleanup := makeTestFile(t, contents)
		defer cleanup()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Cleo")

		got, _ := store.GetPlayerScore("Cleo")
		want := 11

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("record player win for new player", func(t *testing.T) {
		contents := `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`
		database, cleanup := makeTestFile(t, contents)
		defer cleanup()

		store, err := NewFileSystemPlayerStore(database)
		if err != nil {
			t.Errorf("problem creating filesystem player store, %v", err)
		}

		store.RecordWin("Gerard")
		score, found := store.GetPlayerScore("Gerard")
		wantFound := true
		if found != wantFound {
			t.Errorf("got %t want %t", found, wantFound)
		}
		wantScore := 1
		if score != wantScore {
			t.Errorf("got %d want %d", score, wantScore)
		}
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanup := makeTestFile(t, "")
		defer cleanup()

		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanup := makeTestFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanup()

		store, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)

		got := store.GetLeague()

		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func makeTestFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()
	file, err := os.CreateTemp("", "database")
	if err != nil {
		t.Fatalf("Error creating temp database file: %v", err)
	}
	file.WriteString(initialData)

	cleanup := func() {
		file.Close()
		os.Remove(file.Name())
	}
	return file, cleanup
}
