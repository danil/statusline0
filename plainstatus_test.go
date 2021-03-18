// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plainstatus_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/danil/plainstatus"
)

var LoadAverage1FormatTests = []struct {
	name     string
	index    int
	input    string
	expected string
}{
	{
		name:     "non-zero fractional part and zero integer part",
		index:    getLine(),
		input:    "0.42 0.21 0.26 1/1331 29136",
		expected: ".42",
	},
	{
		name:     "non-zero fractional part and non-zero integer part",
		index:    getLine(),
		input:    "1.40 0.21 0.26 1/1331 29136",
		expected: "1.40",
	},
	{
		name:     "zero",
		index:    getLine(),
		input:    "0 0.21 0.26 1/1331 29136",
		expected: ".00",
	},
	{
		name:     "zero fractional part and zero integer part",
		index:    getLine(),
		input:    "0.0 0.21 0.26 1/1331 29136",
		expected: ".00",
	},
	{
		name:     "empty",
		index:    getLine(),
		input:    "",
		expected: "ERR:EOF",
	},
	{
		name:     "blank",
		index:    getLine(),
		input:    " ",
		expected: "ERR:EOF",
	},
	{
		name:     "invalid syntax",
		index:    getLine(),
		input:    "foobar",
		expected: "ERR:str", // is an abbreviation of the `strconv.ParseFloat: parsing "": invalid syntax`
	},
}

func Test_LoadAverage1Format(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, test := range LoadAverage1FormatTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, test.index)
			s := fmt.Sprintf("%s", plainstatus.LoadAverage1(test.input))
			if s != test.expected {
				t.Errorf("unexpected load average %#v, expected %#v %s", s, test.expected, linkToExample)
			}
		})
	}
}

var BatterySignFormatTests = []struct {
	name     string
	index    int
	input    plainstatus.BatterySign
	expected string
}{
	{
		name:     "three digits",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: "100", Icon: "⚡"},
		expected: "100⚡",
	},
	{
		name:     "two digits, charging",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: "42＋", Icon: "⚡"},
		expected: "42＋⚡",
	},
	{
		name:     "one digit, discharging",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: "1-", Icon: "⚡"},
		expected: "1-⚡",
	},
	{
		name:     "empty",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: "", Icon: "⚡"},
		expected: "ERR:n/a⚡",
	},
	{
		name:     "blank",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: " ", Icon: "⚡"},
		expected: "ERR:n/a⚡",
	},
	{
		name:     "long",
		index:    getLine(),
		input:    plainstatus.BatterySign{Power: "something went wrong", Icon: "⚡"},
		expected: "somethi⚡",
	},
}

func Test_BatterySignFormat(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, test := range BatterySignFormatTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, test.index)
			s := fmt.Sprintf("%s", test.input)
			if s != test.expected {
				t.Errorf("unexpected battery %#v, expected %#v %s", s, test.expected, linkToExample)
			}
		})
	}
}

var DegreesPrefixFormatTests = []struct {
	name     string
	index    int
	input    plainstatus.DegreesPrefix
	expected string
}{
	{
		name:     "three digits",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: "100", Degree: "°"},
		expected: "°100",
	},
	{
		name:     "two digits",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: "42", Degree: "°"},
		expected: "°42",
	},
	{
		name:     "one digit, minus",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: "-1", Degree: "°"},
		expected: "°-1",
	},
	{
		name:     "empty",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: "", Degree: "°"},
		expected: "ERR:°n/a",
	},
	{
		name:     "blank",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: " ", Degree: "°"},
		expected: "ERR:°n/a",
	},
	{
		name:     "long",
		index:    getLine(),
		input:    plainstatus.DegreesPrefix{Value: "something went wrong", Degree: "°"},
		expected: "°somethi",
	},
}

func Test_DegreesPrefixFormat(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, test := range DegreesPrefixFormatTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, test.index)
			s := fmt.Sprintf("%s", plainstatus.DegreesPrefix(test.input))
			if s != test.expected {
				t.Errorf("unexpected battery %#v, expected %#v %s", s, test.expected, linkToExample)
			}
		})
	}
}

var ConcatUpToMaxRunesTests = []struct {
	name     string
	index    int
	max      int
	strings  []string
	expected string
}{
	{
		name:     "3 digits and 1 unicode char in one string",
		index:    getLine(),
		max:      4,
		strings:  []string{"100٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in one string",
		index:    getLine(),
		max:      42,
		strings:  []string{"100٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in one string",
		index:    getLine(),
		max:      3,
		strings:  []string{"100٪"},
		expected: "100",
	},
	{
		name:     "3 digits + 1 unicode char in two strings",
		index:    getLine(),
		max:      4,
		strings:  []string{"10", "0٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in two strings",
		index:    getLine(),
		max:      42,
		strings:  []string{"10", "0٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in two strings",
		index:    getLine(),
		max:      3,
		strings:  []string{"10", "0٪"},
		expected: "100",
	},
	{
		name:     "3 digits + 1 unicode char in four strings",
		index:    getLine(),
		max:      4,
		strings:  []string{"1", "0", "0", "٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in four strings",
		index:    getLine(),
		max:      42,
		strings:  []string{"1", "0", "0", "٪"},
		expected: "100٪",
	},
	{
		name:     "3 digits + 1 unicode char in four strings",
		index:    getLine(),
		max:      3,
		strings:  []string{"1", "0", "0", "٪"},
		expected: "100",
	},
}

func TestConcatUpToMaxRunes(t *testing.T) {
	_, testFile, _, _ := runtime.Caller(0)
	for _, test := range ConcatUpToMaxRunesTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			linkToExample := fmt.Sprintf("%s:%d", testFile, test.index)
			s := plainstatus.ConcatUpToMaxRunes(test.max, test.strings...)
			if s != test.expected {
				t.Errorf("unexpected string %#v, expected %#v %s",
					s, test.expected, linkToExample)
			}
		})
	}
}
