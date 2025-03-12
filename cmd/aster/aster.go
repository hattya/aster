//
// aster/cmd/aster :: aster.go
//
//   Copyright (c) 2014-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
		<impl> is one of the following:
		  - freedesktop, fdo
		  - gntp
		  - windows (only available on Windows)

		The -n flag takes precedence over the -g flag

		<duration> is an integer and time unit. Valid time units are "ns", "us", "ms",
		"s", "m", and "h"
	`))
	var g aster.GNTPValue
	app.Flags.Var("g", &g, "notify to Growl (default: localhost:23053)")
	app.Flags.MetaVar("g", "[=<host>[:<port>]]")
	app.Flags.PrefixChoice("n", "", impls, "notifier implementation")
	app.Flags.MetaVar("n", " <impl>")
	app.Flags.Duration("s", 727*time.Millisecond, "squash events during <duration> (default: %v)")
	app.Flags.MetaVar("s", " <duration>")
	app.Action = cli.Option(watch)
}

var (
	icon = make(map[string]notify.Icon)
	opts = map[string]any{
		"freedesktop:timeout": -1,
		"gntp:enabled":        true,
		"windows:sound":       true,
	}
)

func watch(ctx *cli.Context) error {
	n := newNotifier("Aster", ctx.String("n"))
	if n == nil {
		if ctx.String("n") == "gntp" && ctx.String("g") == "" {
			ctx.Flags.Set("g", "true")
		}
		if g := ctx.String("g"); g != "" {
			c := gntp.New()
			c.Server = g
			c.Name = "Aster"
			n = gntp.NewNotifier(c)
		}
	}
	if n != nil {
		for _, ev := range []string{"success", "failure"} {
			if err := n.Register(ev, icon[ev], opts); err != nil {
				return err
			}
		}
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
