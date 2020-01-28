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
	f := []func() string{
		func() string { return fmt.Sprintf("%4s", plainstatus.LoadAverage1(plainstatus.LoadAverage())) },
		func() string { return fmt.Sprintf(" %d", bytefmt.New(memFree())) },
		func() string { df := diskAvail(); return fmt.Sprintf(" %d", df) },
		func() string { return time.Now().Local().Format(" Jan-02 MST 15:04") },
	}
	if *once {
		plainstatus.Run(os.Stdout, f...)
	} else {
		plainstatus.Run(plainstatus.Xsetroot{Interval: 1 * time.Second}, f...)
	}
}
