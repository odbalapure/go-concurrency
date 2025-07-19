## Go keyword, WatiGroup

### Wait Group

This program won't print both the print statements all the time.  
The main goroutine may finish (and exit the program) before the goroutine `go printSomething(...)` even gets a chance to run.  
Scheduling doesn’t guarantee that the goroutine will execute before main returns.  


```go
func printSomething(s string) {
	fmt.Println(s)
}

func main() {
	go printSomething("Hello Om!")
	printSomething("Hello Om!")
}
```

A hacky way to print both the lines would be ues `Sleep` or yield the CPU using `Gosched()`.

```go
go printSomething("Hello Om!")
// time.Sleep(time.Second * 1)
runtime.Gosched()
printSomething("Hello Om!")
```

We can even loop through the results to print something

```go
func main() {
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

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, word))
	}

	time.Sleep(time.Second * 1)
	fmt.Println("This is the last line")
}
```

> 1 second is enough for these tiny goroutines to finish; if the sleep time is reduced, we'll likely miss some lines.

So that's where `WaitGroup` comes into the picture.

**wg.Add(n)**
- Increases (or decreases) the internal counter by n.
- It's used to tell the WaitGroup how many goroutines you're going to wait for.
- Important: You typically call wg.Add(1) before starting each goroutine.

**wg.Done()**
- Decrements the counter by 1.
- It’s a shorthand for wg.Add(-1).
- You usually call this at the end of a goroutine (via defer).

**wg.Wait()**
- Blocks the calling goroutine (e.g., main) until the counter reaches 0.
- It does NOT block just because you're in main. It will block in any goroutine where you call it.
- Once the counter is zero, Wait() unblocks and continues executing.

### File descriptors

- `STDIN`: process input (default: keyboard).
- `STDOUT`: normal output (default: terminal).
- `STDERR`: error messages (default: terminal, separate from stdout).

These 3 are called as file descriptors:
- A network port is an address for network communication.
- A file descriptor (FD) is a handle (number) for I/O inside a process.

> Pipe is an `stdout` of one process and the `stdin` of another process. A typical example can be `ls | grep "foo"`.

Where, `ls` and `grep` are two processes.

### Writing tests for WatiGroup

```go
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
```

**NOTE**: StdIn, StdOut, StdErr are constants
```go
var (
	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
	Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```

> `/dev/stdout` is not some magical, per-process file. It’s just a symbolic link to `/proc/self/fd/1`. And the key here is `self`.

So for a process with PID; is the directory "created" by PID 1 itself?
No. The kernel maintains `/proc`.
When **PID 1** is born, the kernel just exposes its data under `/proc/1`. The process doesn’t “create” that directory; it’s a reflection of what’s already in the kernel.

`/proc` is basically like CNN for the kernel—
it broadcasts live, real-time information about everything the kernel is tracking: processes, CPU stats, memory usage, network info, etc.
The difference? You can not only watch the news, but also poke the news anchor in the face (e.g., by writing to some files in `/proc` to change kernel parameters).

While `/proc` shows processes and kernel state, `/sys` shows devices, drivers, and kernel settings in a much cleaner way. It’s part of *sysfs*, another virtual filesystem.

**TL;DR**,
- Driver: A driver is just a piece of code in the kernel that knows how to talk to a specific piece of hardware (or virtual device).
- Bus: Is like a communication highway between the CPU and devices. It’s how devices are connected and discovered by the kernel.

> Driver + Bus = Device Access
