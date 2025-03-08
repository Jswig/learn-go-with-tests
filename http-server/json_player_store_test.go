package main

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"testing"
)

func TestNewJSONPlayerStore(t *testing.T) {
	dir := t.TempDir()
	filepath := filepath.Join(dir, "test_database.json")

	store := NewJSONPlayerStore(filepath, Scores{})

	fileContents, err := os.ReadFile(store.fileName)
	if err != nil {
		t.Fatalf("Error reading store file: %v", err)
	}
	var got Scores
	err = json.Unmarshal(fileContents, &got)
	if err != nil {
		t.Fatalf("error decoding JSON: %v", err)
	}
	want := Scores{}
	if !maps.Equal(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestGetPlayerScore(t *testing.T) {
	dir := t.TempDir()
	filepath := filepath.Join(dir, "test_database.json")

	scores := Scores{"Pepper": 1, "Anna": 3}
	store := NewJSONPlayerStore(filepath, scores)

	t.Run("Score for player not in database", func(t *testing.T) {
		_, isSet := store.GetPlayerScore("Joan")
		want := false
		if isSet != want {
			t.Errorf("got %t want %t", isSet, want)
		}
	})

	t.Run("score for player in the database", func(t *testing.T) {
		score, isSet := store.GetPlayerScore("Pepper")
		wantIsSet := true
		if isSet != wantIsSet {
			t.Errorf("got %t want %t", isSet, wantIsSet)
		}
		wantScore := 1
		if score != wantScore {
			t.Errorf("got %d want %d", score, wantScore)
		}
	})
}

func TestRecordWin(t *testing.T) {
	dir := t.TempDir()
	filepath := filepath.Join(dir, "test_database.json")

	store := NewJSONPlayerStore(filepath, Scores{})

	t.Run("Record win for player not in database", func(t *testing.T) {
		store.RecordWin("Dimitri")
		score, isSet := store.GetPlayerScore("Dimitri")
		wantIsSet := true
		if isSet != wantIsSet {
			t.Errorf("got %t want %t", isSet, wantIsSet)
		}
		wantScore := 1
		if score != wantScore {
			t.Errorf("got %d want %d", score, wantScore)
		}
	})
	t.Run("Record win for player already in database", func(t *testing.T) {
		store.RecordWin("Dimitri")
		score, isSet := store.GetPlayerScore("Dimitri")
		wantIsSet := true
		if isSet != wantIsSet {
			t.Errorf("got %t want %t", isSet, wantIsSet)
		}
		wantScore := 2
		if score != wantScore {
			t.Errorf("got %d want %d", score, wantScore)
		}
	})
}
