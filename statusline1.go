// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statusline1

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Load average fmt.Formatter

type LoadAverage1 string

func (la LoadAverage1) Format(f fmt.State, c rune) {
	var flt float32
	_, err := fmt.Sscanf(string(la), "%f", &flt)
	if err != nil {
		s := ConcatUpToMaxRunes(7, "ERR:", err.Error())
		f.Write([]byte(strings.TrimSpace(s)))
		return
	}
	if flt == 0 {
		f.Write([]byte(".00"))
	} else if flt < 1 {
		f.Write([]byte(fmt.Sprintf("%.2f", flt)[1:]))
	} else {
		f.Write([]byte(fmt.Sprintf("%.2f", flt)))
	}
}

// Battery fmt.Formatter

type BatterySign struct {
	Power string
	Sign  int8
	Plus  string
	Minus string
	Icon  string
}

func (bat BatterySign) Format(f fmt.State, c rune) {
	s := strings.TrimSpace(string(bat.Power))
	if len(s) == 0 {
		f.Write([]byte("ERR:n/a" + bat.Icon))
		return
	}
	var sgn string
	if bat.Sign > 0 {
		sgn = bat.Plus
	} else if bat.Sign < 0 {
		sgn = bat.Minus
	}
	f.Write([]byte(fmt.Sprintf("%s%s%s", ConcatUpToMaxRunes(7, s), sgn, bat.Icon)))
}

// Temperature fmt.Formatter

type DegreesPrefix struct{ Value, Degree string }

func (td DegreesPrefix) Format(f fmt.State, c rune) {
	s := strings.TrimSpace(string(td.Value))
	if len(s) == 0 {
		f.Write([]byte("ERR:" + td.Degree + "n/a"))
		return
	}
	f.Write([]byte(fmt.Sprintf("%s%s", td.Degree, ConcatUpToMaxRunes(7, s))))
}

func ConcatUpToMaxRunes(max int, p ...string) string {
	var runes []rune
	for _, s := range p {
		runes = append(runes, []rune(s)...)
		if len(runes) > max {
			return string(runes[:max])
		}
	}
	return string(runes)
}

func FileName(pattern string) (string, error) {
	a, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(a) == 0 {
		return "", os.ErrNotExist
	}
	return a[0], nil
}

// Run writes to the io writer every second
// or writes to os stdout and return
func Run(w io.Writer, f ...func() string) error {
	if w == os.Stdout {
		for _, f := range f {
			_, err := fmt.Fprint(w, f())
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprint(w, "\n")
		return err
	}
	var buf bytes.Buffer
	for {
		for _, f := range f {
			buf.WriteString(f())
		}
		_, err := fmt.Fprint(w, buf.String())
		if err != nil {
			return err
		}
		buf.Reset()
	}
}
