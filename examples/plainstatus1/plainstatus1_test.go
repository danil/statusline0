package main

import (
	"os"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	done := make(chan struct{})

	go func() {
		os.Args = append(os.Args, "-1")
		main()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("test takes too long")
	}
}
