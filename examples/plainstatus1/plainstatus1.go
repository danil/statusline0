package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/danil/bytefmt"
	"github.com/danil/plainstatus"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	once := flag.Bool("1", false, "print to stdout and exit")
	flag.Parse()
	batt := plainstatus.BatterySign{Plus: "＋", Minus: "−", Icon: "⚡"}
	temp := plainstatus.DegreesPrefix{Degree: "°"}
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
	tPth, _ := plainstatus.FileName("/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp1_input")
	f := []func() string{
		func() string { batt.Power, batt.Sign = plainstatus.BatteryPercent(); return fmt.Sprint(batt) },
		func() string { return fmt.Sprintf("%4s", plainstatus.LoadAverage1(plainstatus.LoadAverage())) },
		func() string { temp.Value = plainstatus.Temperature(tPth); return fmt.Sprint(temp) },
		func() string { mfa := memFreeAvail(); return fmt.Sprintf(" %d~%d", mfa[:]...) },
		func() string { df := diskAvail(); return fmt.Sprintf(" %d", df) },
		func() string { return time.Now().Local().Format(" Jan-02-Mon MST 15:04") },
		func() string { return "      " },
	}
	if *once {
		plainstatus.Run(os.Stdout, f...)
	} else {
		plainstatus.Run(plainstatus.Xsetroot{Interval: 1 * time.Second}, f...)
	}
}
