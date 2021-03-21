// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statusline0

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func LoadAverage() string {
	p, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return "ERR:" + err.Error()
	}
	return string(p)
}

type Battery struct {
	Batteries string
	Online    string
}

func (batt Battery) Percent() (string, int8) {
	var full, now, perc int

	plugged, err := ioutil.ReadFile(batt.Online)
	if err != nil {
		return "ERR:" + err.Error(), 0
	}

	if batt.Batteries[len(batt.Batteries)-1] == os.PathSeparator {
		return "ERR:batteries path trailing slash", 0
	}

	dir, prefix := path.Split(batt.Batteries)
	if dir == "" {
		return "ERR:missing batteries directory", 0
	}
	if prefix == "" {
		return "ERR:missing batteries prefix", 0
	}

	batts, err := ioutil.ReadDir(dir)
	if err != nil {
		return "ERR:" + err.Error(), 0
	}

	for _, b := range batts {
		name := b.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		full += batteryPercentRead(dir, name, "full")
		now += batteryPercentRead(dir, name, "now")
	}

	if full == 0 { // Battery found but no readable full file.
		return "ERR", 0
	}

	perc = now * 100 / full
	if string(plugged) == "1\n" {
		return strconv.Itoa(perc), +1
	}

	return strconv.Itoa(perc), -1
}

func batteryPercentRead(dir, name, field string) int {
	var path = dir + name + "/"
	var file []byte
	if tmp, err := ioutil.ReadFile(path + "energy_" + field); err == nil {
		file = tmp
	} else if tmp, err := ioutil.ReadFile(path + "charge_" + field); err == nil {
		file = tmp
	} else {
		return 0
	}
	if perc, err := strconv.Atoi(strings.TrimSpace(string(file))); err == nil {
		return perc
	}
	return 0
}

func Temperature(p string) string {
	a, err := ioutil.ReadFile(p)
	if err != nil {
		return "ERR:" + err.Error()
	}
	i, err := strconv.Atoi(strings.TrimSpace(string(a)))
	if err != nil {
		return "ERR:" + err.Error()
	}
	return strconv.Itoa(i / 1000)
}

type Xsetroot struct{ Interval time.Duration }

// Write updates the xsetroot (each write overwrites previous value),
// and then sleeps until the start of the next second, plus the time interval
func (x Xsetroot) Write(p []byte) (int, error) {
	var now = time.Now()
	err := exec.Command("xsetroot", "-name", string(p)).Run()
	if err != nil {
		return 0, err
	}
	time.Sleep(now.Truncate(time.Second).Add(x.Interval).Sub(now))
	return len(p), nil
}
