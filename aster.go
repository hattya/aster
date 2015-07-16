//
// aster :: aster.go
//
//   Copyright (c) 2014-2015 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

package main

import (
	"os"
	"runtime"
	"strings"

	"github.com/hattya/go.cli"
)

const version = "0.0+"

var app = cli.NewCLI()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Version = version
	app.Usage = []string{
		"[options]",
		"init [<template>...]",
	}
	app.Epilog = strings.TrimSpace(`
<duration> is an integer and time unit. Valid time units are "ns", "us", "ms",
"s", "m" and "h"
`)
	app.Action = cli.Option(watch)

	if err := app.Run(os.Args[1:]); err != nil {
		switch err.(type) {
		case cli.FlagError:
			os.Exit(2)
		}
		os.Exit(1)
	}
}

func watch(*cli.Context) error {
	a, err := newAsterfile()
	if err != nil {
		return err
	}

	w, err := newWatcher(a)
	if err != nil {
		return err
	}
	defer w.Close()

	w.Watch()
	return nil
}
