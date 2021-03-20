// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/danil/bytefmt"
	"github.com/danil/statusline1"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	once := flag.Bool("1", false, "print to stdout and exit")
	flag.Parse()

	if *once {
		statusline1.Run(os.Stdout, status()...)
	} else {
		statusline1.Run(statusline1.Xsetroot{Interval: 1 * time.Second}, status()...)
	}
}

func status() []func() string {
	batts := statusline1.BatterySign{
		Plus:  "＋",
		Minus: "−",
		Icon:  "⚡",
	}
	batt := statusline1.Battery{
		Batteries: "/sys/class/power_supply/BAT",
		Online:    "/sys/class/power_supply/AC/online",
	}

	temp := statusline1.DegreesPrefix{Degree: "°"}

	memFreeAvail := func() [2]interface{} {
		var a [2]interface{}
		m, err := mem.VirtualMemory()
		if err != nil {
			return a
		}
		a[0] = bytefmt.New(m.Free)
		a[1] = bytefmt.New(m.Available)
		return a
	}

	diskAvail := func() bytefmt.Bytes {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs("/", &fs)
		if err != nil {
			return bytefmt.New(0)
		}
		return bytefmt.New(fs.Bavail * uint64(fs.Bsize))
	}

	tPth, _ := statusline1.FileName("/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp1_input")

	return []func() string{
		func() string { batts.Power, batts.Sign = batt.Percent(); return fmt.Sprint(batts) },
		func() string { return fmt.Sprintf("%4s", statusline1.LoadAverage1(statusline1.LoadAverage())) },
		func() string { temp.Value = statusline1.Temperature(tPth); return fmt.Sprint(temp) },
		func() string { mfa := memFreeAvail(); return fmt.Sprintf(" %d~%d", mfa[:]...) },
		func() string { df := diskAvail(); return fmt.Sprintf(" %d", df) },
		func() string { return time.Now().Local().Format(" Jan-02-Mon MST 15:04") },
	}
}
