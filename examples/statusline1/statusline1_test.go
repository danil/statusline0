// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"regexp"
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

func Test_stat(t *testing.T) {
	var buf bytes.Buffer

	for _, f := range status() {
		_, err := buf.WriteString(f())
		if err != nil {
			t.Fatalf("unexpected bytes buffer write error: %s", err)
		}
	}

	r := ".*\\.[0-9]+°.* [0-9]+[a-zA-Z]~[0-9]+[a-zA-Z]+ [0-9]+[a-zA-Z]+ [a-zA-Z]+-[0-9]+-[a-z-A-Z]+ [A-Z]+ [0-9]+:[0-9]+"

	ok, err := regexp.MatchString(r, buf.String())
	if err != nil {
		t.Fatalf("regexp match string error: %s", err)
	}

	if !ok {
		t.Errorf("status line does not match regexp, line: %q, regexp: %q", buf.String(), r)
	}
}
