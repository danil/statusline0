// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plainstatus

import (
	"io/ioutil"
	"os/exec"
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

func BatteryPercent() (string, int8) {
	const powerSupply = "/sys/class/power_supply/"
	var enFull, enNow, enPerc int
	var plugged, err = ioutil.ReadFile(powerSupply + "AC/online")
	if err != nil {
		return "ERR:" + err.Error(), 0
	}
	batts, err := ioutil.ReadDir(powerSupply)
	if err != nil {
		return "ERR:" + err.Error(), 0
	}
	readval := func(name, field string) int {
		var path = powerSupply + name + "/"
		var file []byte
		if tmp, err := ioutil.ReadFile(path + "energy_" + field); err == nil {
			file = tmp
		} else if tmp, err := ioutil.ReadFile(path + "charge_" + field); err == nil {
			file = tmp
		} else {
			return 0
		}
		if ret, err := strconv.Atoi(strings.TrimSpace(string(file))); err == nil {
			return ret
		}
		return 0
	}
	for _, batt := range batts {
		name := batt.Name()
		if !strings.HasPrefix(name, "BAT") {
			continue
		}
		enFull += readval(name, "full")
		enNow += readval(name, "now")
	}
	if enFull == 0 { // Battery found but no readable full file.
		return "ERR", 0
	}
	enPerc = enNow * 100 / enFull
	if string(plugged) == "1\n" {
		return strconv.Itoa(enPerc), +1
	}
	return strconv.Itoa(enPerc), -1
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
