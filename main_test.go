// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plainstatus_test

import "runtime"

func getLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
