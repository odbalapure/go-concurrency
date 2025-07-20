package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestIncome(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	if !strings.Contains(output, "$34320.00") {
		t.Errorf("Wrong balance returned, expected $34320.00 got %s", output)
	}

	os.Stdout = stdOut
}
