//
// aster :: init.go
//
//   Copyright (c) 2014 Akinori Hattori <hattya@gmail.com>
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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hattya/go.cli"
)

func init() {
	app.Add(&cli.Command{
		Name:  []string{"init"},
		Usage: "[<template>...]",
		Desc: strings.TrimSpace(`
generate an Asterfile in the current directory

    Create an Asterfile in the current directory if it does not exist,
    and add specified template files to it.

    Template files are located in:

      UNIX:    $XDG_CONFIG_HOME/aster/template/<template>
      Windows: %APPDATA%\Aster\template\<template>
`),
		Flags:  cli.NewFlagSet(),
		Action: init_,
	})
}

func init_(ctx *cli.Context) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	template := filepath.Join(dir, "template")
	if _, err := os.Stat(template); err != nil {
		return nil
	}

	f, err := os.OpenFile("Asterfile", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	newline := false
	if fi, err := f.Stat(); err == nil && 0 < fi.Size() {
		newline = true
	}

	for _, a := range ctx.Args {
		t, err := os.Open(filepath.Join(template, a))
		if err != nil {
			return fmt.Errorf("template '%v' is not found", a)
		}
		if newline {
			fmt.Fprintln(f)
		}
		scanner := bufio.NewScanner(t)
		for scanner.Scan() {
			fmt.Fprintln(f, scanner.Text())
		}
		t.Close()
		if scanner.Err() != nil {
			return fmt.Errorf("error occurred while processing template '%v'", a)
		}
		newline = true
	}
	return nil
}
