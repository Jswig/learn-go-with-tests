package poker_test

import (
	"bytes"
	"fmt"
	"poker"
	"strings"
	"testing"
	"time"
)

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

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		dummyAlerter := &SpyBlindAlerter{}
		dummyStdout := &bytes.Buffer{}

		game := poker.NewGame(playerStore, dummyAlerter)
		cli := poker.NewCLI(in, dummyStdout, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("3\nCleo wins\n")
		playerStore := &poker.StubPlayerStore{}
		dummyAlerter := &SpyBlindAlerter{}
		dummyStdout := &bytes.Buffer{}

		game := poker.NewGame(playerStore, dummyAlerter)
		cli := poker.NewCLI(in, dummyStdout, game)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		dummyStdout := &bytes.Buffer{}

		game := poker.NewGame(playerStore, blindAlerter)
		cli := poker.NewCLI(in, dummyStdout, game)
		cli.PlayPoker()

		cases := []Alert{
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
		}

		for i, want := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", want.amount, want.scheduledAt), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		dummyPlayerStore := &poker.StubPlayerStore{}
		in := strings.NewReader("7\n")
		stdout := &bytes.Buffer{}

		game := poker.NewGame(dummyPlayerStore, blindAlerter)
		cli := poker.NewCLI(in, stdout, game)
		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []Alert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}

func TestGame(t *testing.T) {
	t.Run("finishing records the winner", func(t *testing.T) {
		winners := []string{"Chris", "Cleo"}
		for _, winner := range winners {
			playerStore := &poker.StubPlayerStore{}
			dummyAlerter := &SpyBlindAlerter{}
			game := poker.NewGame(playerStore, dummyAlerter)
			game.Finish(winner)
			poker.AssertPlayerWin(t, playerStore, winner)
		}
	})
}

func assertScheduledAlert(t *testing.T, got Alert, want Alert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}
	if got.scheduledAt != want.scheduledAt {
		t.Errorf("got scheduled time of %v, want %v", got.scheduledAt, want.scheduledAt)
	}
}
