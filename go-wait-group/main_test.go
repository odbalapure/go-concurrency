package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

// Writing tests for WaitGroup
func TestPrintSomething(t *testing.T) {
	// When we run our program; it prints/writes to `os.StdOut`
	// So preserve it to avoid screwing up the rest of your test process
	stdOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	var wg sync.WaitGroup
	wg.Add(1)

	go printSomething("epsilon", &wg)

	wg.Wait()

	// Close the write end of the pipe
	// `w.Close()` signals end of output to the read side (r).
	// Without it, `io.ReadAll(r)` will block forever, waiting for more bytes.
	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "epsilon") {
		t.Errorf("Expected to find epsilon")
	}
}
