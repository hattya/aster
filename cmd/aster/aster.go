//
// aster/cmd/aster :: aster.go
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
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
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/hattya/aster"
	"github.com/hattya/go.cli"
	"github.com/hattya/go.notify"
	"github.com/hattya/go.notify/gntp"
)

var app = cli.NewCLI()

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := app.Run(os.Args[1:]); err != nil {
		switch err.(type) {
		case cli.FlagError:
			os.Exit(2)
		case cli.Interrupt:
			os.Exit(128 + 2)
		}
		os.Exit(1)
	}
}

func init() {
	app.Version = aster.Version
	app.Usage = []string{
		"[options]",
		"init [<template>...]",
	}
	app.Epilog = strings.TrimSpace(cli.Dedent(`
		<duration> is an integer and time unit. Valid time units are "ns", "us", "ms",
		"s", "m", and "h"
	`))
	var g aster.GNTPValue
	app.Flags.Var("g", &g, "notify to Growl (default: localhost:23053)")
	app.Flags.MetaVar("g", "[=<host>[:<port>]]")
	app.Flags.Duration("s", 727*time.Millisecond, "squash events during <duration> (default: %v)")
	app.Flags.MetaVar("s", " <duration>")
	app.Action = cli.Option(watch)
}

func watch(ctx *cli.Context) error {
	var n notify.Notifier
	if ctx.String("g") != "" {
		c := gntp.New()
		c.Server = ctx.String("g")
		c.Name = "Aster"
		n = gntp.NewNotifier(c)
	}
	a, err := aster.New(ctx.UI, n)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		ctx.Interrupt()
	}()

	w, err := aster.NewWatcher(ctx.Context(), a)
	if err != nil {
		return err
	}
	defer w.Close()

	w.Squash = ctx.Duration("s")
	return w.Watch()
}
