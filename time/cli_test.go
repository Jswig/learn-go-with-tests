package poker_test

import (
	"bytes"
	"fmt"
	"poker"
	"slices"
	"strings"
	"testing"
	"time"
)

type GameSpy struct {
	NumPlayers int
	Winner     string
}

func (g *GameSpy) Start(numPlayers int) {
	g.NumPlayers = numPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.Winner = winner
}

func TestCli(t *testing.T) {
	t.Run("Prompts for number of players in the game and starts it", func(t *testing.T) {
		input := strings.NewReader("7\n")
		output := &bytes.Buffer{}
		game := &GameSpy{}

		cli := poker.NewCLI(input, output, game)
		cli.PlayPoker()

		gotPrompt := output.String()
		if gotPrompt != poker.PlayerPrompt {
			t.Errorf("got %q, want %q", gotPrompt, poker.PlayerPrompt)
		}

		gotNumPlayers := game.NumPlayers
		wantNumPlayers := 7
		if gotNumPlayers != wantNumPlayers {
			t.Errorf("wanted Start() called with %d but got %d", wantNumPlayers, gotNumPlayers)
		}
	})

	t.Run("Parses the winner and finishes the game", func(t *testing.T) {
		input := strings.NewReader("2\nMatthew wins")
		dummyOutput := &bytes.Buffer{}
		game := &GameSpy{}

		cli := poker.NewCLI(input, dummyOutput, game)
		cli.PlayPoker()

		got := game.Winner
		want := "Matthew"
		if got != want {
			t.Errorf("wanted Finish called with %s but got %s", want, got)
		}
	})
}

type Alert struct {
	scheduledAt time.Duration
	amount      int
}

func (a Alert) String() string {
	return fmt.Sprintf("%d at %v", a.amount, a.scheduledAt)
}

type SpyBlindAlerter struct {
	alerts []Alert
}

func (a *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	a.alerts = append(a.alerts, Alert{duration, amount})
}

func TestTexasHoldEm(t *testing.T) {
	t.Run("finishing records the winner", func(t *testing.T) {
		winners := []string{"Chris", "Cleo"}
		for _, winner := range winners {
			playerStore := &poker.StubPlayerStore{}
			dummyAlerter := &SpyBlindAlerter{}
			game := poker.NewTexasHoldEm(playerStore, dummyAlerter)
			game.Finish(winner)
			poker.AssertPlayerWin(t, playerStore, winner)
		}
	})

	t.Run("starting schedules printing blind values", func(t *testing.T) {
		cases := []struct {
			numPlayers     int
			expectedAlerts []Alert
		}{
			{
				5,
				[]Alert{
					{0 * time.Second, 100},
					{10 * time.Minute, 200},
					{20 * time.Minute, 300},
					{30 * time.Minute, 400},
					{40 * time.Minute, 500},
					{50 * time.Minute, 600},
					{60 * time.Minute, 800},
					{70 * time.Minute, 1000},
					{80 * time.Minute, 2000},
					{90 * time.Minute, 4000},
					{100 * time.Minute, 8000},
				},
			},
			{
				7,
				[]Alert{
					{0 * time.Second, 100},
					{12 * time.Minute, 200},
					{24 * time.Minute, 300},
					{36 * time.Minute, 400},
					{48 * time.Minute, 500},
					{60 * time.Minute, 600},
					{72 * time.Minute, 800},
					{84 * time.Minute, 1000},
					{96 * time.Minute, 2000},
					{108 * time.Minute, 4000},
					{120 * time.Minute, 8000},
				},
			},
		}

		for _, c := range cases {
			dummyPlayerStore := &poker.StubPlayerStore{}
			alerter := &SpyBlindAlerter{}
			game := poker.NewTexasHoldEm(dummyPlayerStore, alerter)
			game.Start(c.numPlayers)

			if len(alerter.alerts) != len(c.expectedAlerts) {
				t.Errorf("Expected %d alerts, but %d alerts were schedule", len(alerter.alerts), len(c.expectedAlerts))
			}

			if !slices.Equal(alerter.alerts, c.expectedAlerts) {
				t.Errorf("Expected alerts %v, got %v", c.expectedAlerts, alerter.alerts)
			}
		}
	})
}
