package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	// Best to add it, on the 1st line
	// A function might "panic" and Done() might never get called
	// The program will deadlock because wg.Wait() will block forever
	defer wg.Done()

	fmt.Println(s)

	// wg.Done()
}

func main() {
	var wg sync.WaitGroup
	words := []string{
		"alpha",
		"beta",
		"gamma",
		"pi",
		"zeta",
		"eta",
		"theta",
		"eplison",
	}

	wg.Add(len(words))

	for i, word := range words {
		// Its best to pass `WaitGroup` as a pointer
		// We shouldn't go around copying and modifying them
		go printSomething(fmt.Sprintf("%d: %s", i, word), &wg)
	}

	// Wait() blocks the `main` thread until all goroutines have finished.
	wg.Wait()

	// This line will be printed the last
	fmt.Println("This is the last line")
}
