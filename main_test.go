package plainstatus_test

import "runtime"

func getLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
