package main

import (
	"sync"
	"testing"
)

func TestUpdateMessage(t *testing.T) {
	msg = "Hello World!"
	var mutex sync.Mutex

	wg.Add(1)
	go updateMessage("Goodbye!", &mutex)
	wg.Wait()

	if msg != "Goodbye!" {
		t.Errorf("Expected Goodbye but got: %s", msg)
	}
}
