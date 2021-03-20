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
	batt := statusline1.BatterySign{Plus: "＋", Minus: "−", Icon: "⚡"}
	temp := statusline1.DegreesPrefix{Degree: "°"}
	memFree := func() uint64 {
		m, err := mem.VirtualMemory()
		if err != nil {
			return 0
		}
		return m.Free
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
	f := []func() string{
		func() string { batt.Power, batt.Sign = statusline1.BatteryPercent(); return fmt.Sprint(batt) },
		func() string { return fmt.Sprintf("%4s", statusline1.LoadAverage1(statusline1.LoadAverage())) },
		func() string { temp.Value = statusline1.Temperature(tPth); return fmt.Sprint(temp) },
		func() string { return fmt.Sprintf(" %d", bytefmt.New(memFree())) },
		func() string { df := diskAvail(); return fmt.Sprintf(" %d", df) },
		func() string { return time.Now().Local().Format(" Jan-02 MST 15:04") },
	}
	if *once {
		statusline1.Run(os.Stdout, f...)
	} else {
		statusline1.Run(statusline1.Xsetroot{Interval: 1 * time.Second}, f...)
	}
}
