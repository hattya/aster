//
// aster :: otto_test_cmd.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
)

var code int

func main() {
	flag.IntVar(&code, "code", 0, "")
	flag.Parse()

	if code == 0 {
		fmt.Fprintln(os.Stdout, "stdout")
	} else {
		fmt.Fprintln(os.Stderr, "stderr")
	}
	os.Exit(code)
}
