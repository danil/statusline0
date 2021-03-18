// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
