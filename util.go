//
// aster :: util.go
//
//   Copyright (c) 2014-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import (
	"io"

	"github.com/hattya/go.cli"
)

func warn(ui *cli.CLI, a ...any) {
	ui.Errorf("%v: ", ui.Name)
	ui.Errorln(a...)
}

func trim(s string) string {
	switch {
	case len(s) > 1 && s[len(s)-2:] == "\r\n":
		s = s[:len(s)-2]
	case len(s) > 0 && s[len(s)-1:] == "\n":
		s = s[:len(s)-1]
	}
	return s
}

var discard io.WriteCloser = devNull(0)

type devNull int

func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func (devNull) Close() error {
	return nil
}
