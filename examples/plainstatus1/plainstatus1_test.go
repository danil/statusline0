package main

import (
	"os"
	"testing"
)

func Test_main(t *testing.T) {
	os.Args = append(os.Args, "-1")
	main()
}
