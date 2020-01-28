package plainstatus_test

import (
	"regexp"
	"testing"

	"github.com/danil/plainstatus"
)

func TestBatteryPercent(t *testing.T) {
	s, _ := plainstatus.BatteryPercent()
	ok, err := regexp.MatchString("^[0-9]", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected battery %s", s)
	}
}

func TestLoadAverage(t *testing.T) {
	s := plainstatus.LoadAverage()
	ok, err := regexp.MatchString("^[0-9]\\.[0-9]{2}", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected load average %s", s)
	}
}

func TestTemperature(t *testing.T) {
	tPth, err := plainstatus.FileName("/sys/devices/platform/coretemp.0/hwmon/hwmon*/temp1_input")
	if err != nil {
		t.Fatal(err)
	}
	if tPth == "" {
		t.Skipf("unexpected temperature file %s", tPth)
	}
	s := plainstatus.Temperature(tPth)
	ok, err := regexp.MatchString("^[0-9]{1,}$", s)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skipf("unexpected temperature %s", s)
	}
}
