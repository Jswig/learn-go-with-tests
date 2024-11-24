package racer

import (
	"fmt"
	"net/http"
	"time"
)

func Racer(firstURL string, secondURL string) (winner string, err error) {
	select {
	case <-ping(firstURL):
		return firstURL, nil
	case <-ping(secondURL):
		return secondURL, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", firstURL, secondURL)
	}
}

func ping(url string) chan struct{} {
	channel := make(chan struct{})
	go func() {
		http.Get(url)
		close(channel)
	}()
	return channel
}
