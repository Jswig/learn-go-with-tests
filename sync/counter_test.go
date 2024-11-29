package counter

import (
	"sync"
	"testing"
)

func assertCounter(t testing.TB, got *Counter, want uint) {
	t.Helper()
	val := got.Value()
	if val != want {
		t.Errorf("got %d, want %d", val, want)
	}
}

func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := NewCounter()
		counter.Inc()
		counter.Inc()
		counter.Inc()

		assertCounter(t, counter, uint(3))
	})

	t.Run("it runs safely concurrently", func(t *testing.T) {
		wantedCount := 1000
		counter := NewCounter()

		wg := sync.WaitGroup{}
		wg.Add(wantedCount)

		for range wantedCount {
			go func() {
				counter.Inc()
				wg.Done()
			}()
		}
		wg.Wait()

		assertCounter(t, counter, uint(wantedCount))
	})
}
