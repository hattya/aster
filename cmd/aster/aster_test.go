//
// aster/cmd/aster :: aster_test.go
//
//   Copyright (c) 2017 Akinori Hattori <hattya@gmail.com>
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
	"io/ioutil"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/aster/internal/test"
	"github.com/hattya/go.cli"
)

func TestInterrupt(t *testing.T) {
	app := clone()
	time.AfterFunc(101*time.Millisecond, app.Interrupt)
	err := test.Sandbox(func() {
		if err := sh.Touch("Asterfile"); err != nil {
			t.Fatal(err)
		}
		switch err := app.Run(nil).(type) {
		case cli.Interrupt:
		default:
			t.Errorf("expected cli.Interrupt, got %#v", err)
		}
	})
	if err != nil {
		t.Error()
	}
}

func TestNotifier(t *testing.T) {
	tests := [][]string{
		{"-g"},
		{"-n", "freedesktop"},
		{"-n", "gntp"},
	}
	if runtime.GOOS == "windows" {
		tests = append(tests, []string{"-n", "windows"})
	}
	for _, args := range tests {
		app := clone()
		time.AfterFunc(101*time.Millisecond, app.Interrupt)
		err := test.Sandbox(func() {
			if err := sh.Touch("Asterfile"); err != nil {
				t.Fatal(err)
			}
			switch err := app.Run(args).(type) {
			case *net.OpError:
			case cli.Interrupt:
			default:
				t.Errorf("expected cli.Interrupt, got %#v", err)
			}
		})
		if err != nil {
			t.Error()
		}
	}
}

func clone() *cli.CLI {
	ui := cli.NewCLI()
	app.Flags.Reset()
	app.Flags.VisitAll(ui.Flags.Add)
	ui.Action = cli.Option(watch)
	ui.Stderr = ioutil.Discard
	return ui
}
