## Race conditions, Channels, Mutex


### Race detector

This enables the race detector in Go when running your program.

```go
go run -race .
```

It detects data races at runtime; situations where two or more goroutines access the same memory concurrently, and at least one of the accesses is a write, without proper synchronization (e.g., without a mutex).  

If a race is detected, it prints detailed debugging info showing which goroutines caused the race and where in the code it happened.

> We can do the same while running tests; `go test -race .`


```zsh
ombalapure@Oms-MacBook-Air race-condition-channel % go run -race . 
==================
WARNING: DATA RACE
Write at 0x0001029b0a00 by goroutine 6:
  main.updateMessage()
    /golang/go-concurrency/race-condition/main.go:15 +0x74
  main.main.gowrap1()
    /golang/go-concurrency/race-condition/main.go:25 +0x40

Previous write at 0x0001029b0a00 by goroutine 7:
  main.updateMessage()
    /golang/go-concurrency/race-condition/main.go:15 +0x74
  main.main.gowrap2()
    /golang/go-concurrency/race-condition/main.go:26 +0x40

Goroutine 6 (running) created at:
  main.main()
    /golang/go-concurrency/race-condition/main.go:25 +0xf4

Goroutine 7 (finished) created at:
  main.main()
    /golang/go-concurrency/race-condition/main.go:26 +0x158
==================
Hello universe
Found 1 data race(s)
exit status 66
```

```go
package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	msg = "Hello world!"

	wg.Add(2)
	updateMessage("Hello universe")
	updateMessage("Hello cosmos")
	wg.Wait()

	fmt.Println(msg)
}
```

To fix this, we could use a mutex.

```go
m.Lock()
msg = s
m.Unlock()
```

Now, this fixes the race warnings because the race detector no longer sees unsynchronized concurrent access to msg.

### Writing a test

```go
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
```
