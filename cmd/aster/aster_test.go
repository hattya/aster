//
// aster/cmd/aster :: aster_test.go
//
//   Copyright (c) 2017-2022 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package main

import (
	"io"
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
	ui.Stderr = io.Discard
	return ui
}
