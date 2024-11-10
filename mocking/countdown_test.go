package mocking

import (
	"bytes"
	"slices"
	"testing"
	"time"
)

type CountdownOperation string

const (
	write CountdownOperation = "write"
	sleep CountdownOperation = "sleep"
)

type SpyCountdownOperations struct {
	Calls []CountdownOperation
}

func (s *SpyCountdownOperations) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func TestCountdownOutput(t *testing.T) {
	buffer := &bytes.Buffer{}

	Countdown(buffer, &SpyCountdownOperations{})

	got := buffer.String()
	want := "3\n2\n1\nGo!\n"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestCountdownOperationOrder(t *testing.T) {
	spySleepPrinter := &SpyCountdownOperations{}
	Countdown(spySleepPrinter, spySleepPrinter)

	want := []CountdownOperation{
		write,
		sleep,
		write,
		sleep,
		write,
		sleep,
		write,
	}

	if !slices.Equal(want, spySleepPrinter.Calls) {
		t.Errorf("wanted call %v got %v", want, spySleepPrinter.Calls)
	}
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second

	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{
		duration: sleepTime,
		sleep:    spyTime.Sleep,
	}
	sleeper.sleep(sleepTime)

	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
	}
}
