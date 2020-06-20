//
// aster/cmd/aster :: init.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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

const lf = "\n"

var configDir func() (string, error)

func init() {
	app.Add(&cli.Command{
		Name:  []string{"init"},
		Usage: "[<template>...]",
		Desc: strings.TrimSpace(cli.Dedent(`
			generate an Asterfile in the current directory

			  Create an Asterfile in the current directory if it does not exist,
			  and add specified template files to it.

			  Template files are located in:

			  UNIX:    $XDG_CONFIG_HOME/aster/template/<template>
			  Windows: %APPDATA%\Aster\template\<template>
		`)),
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

	f, err := os.OpenFile("Asterfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var newline []byte
	if fi, err := f.Stat(); err == nil && 0 < fi.Size() {
		var off int64
		if 1 < fi.Size() {
			off = -2
		} else {
			off = -1
		}
		if _, err := f.Seek(off, os.SEEK_END); err == nil {
			b := make([]byte, -off)
			if n, err := f.Read(b); err == nil {
				switch n {
				case 2:
					if b[0] == '\r' && b[1] == '\n' {
						newline = b
						break
					}
					b[0] = b[1]
					fallthrough
				case 1:
					switch b[0] {
					case '\n', '\r':
						newline = b[:1]
					}
				}
			}
		}
		if len(newline) == 0 {
			newline = []byte(lf)
		}
	}

	for _, a := range ctx.Args {
		t, err := os.Open(filepath.Join(template, a))
		if err != nil {
			return fmt.Errorf("template '%v' is not found", a)
		}
		if 0 < len(newline) {
			f.Write(newline)
		} else {
			newline = []byte(lf)
		}
		s := bufio.NewScanner(t)
		for s.Scan() {
			f.Write(s.Bytes())
			f.Write(newline)
		}
		t.Close()
		if s.Err() != nil {
			return fmt.Errorf("error occurred while processing template '%v'", a)
		}
	}
	return nil
}
