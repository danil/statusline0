// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statusline0_test

import (
	"regexp"
	"testing"

	"github.com/danil/statusline0"
)

func TestBatteryPercent(t *testing.T) {
	b := statusline0.Battery{
		Batteries: "/sys/class/power_supply/BAT",
		Online:    "/sys/class/power_supply/AC/online",
	}
	s, _ := b.Percent()
	ok, err := regexp.MatchString("^[0-9]", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected battery %s", s)
	}
}

func TestLoadAverage(t *testing.T) {
	s := statusline0.LoadAverage()
	ok, err := regexp.MatchString("^[0-9]\\.[0-9]{2}", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected load average %s", s)
	}
}

func TestTemperature(t *testing.T) {
	tPth, err := statusline0.FileName("/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp1_input")
	if err != nil {
		t.Fatal(err)
	}
	if tPth == "" {
		t.Skipf("unexpected temperature file %s", tPth)
	}
	s := statusline0.Temperature(tPth)
	ok, err := regexp.MatchString("^[0-9]{1,}$", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected temperature %s", s)
	}
}
