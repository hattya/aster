//
// aster :: util.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import (
	"io"

	"github.com/hattya/go.cli"
)

func warn(ui *cli.CLI, a ...interface{}) {
	ui.Errorf("%v: ", ui.Name)
	ui.Errorln(a...)
}

func trimNewline(s string) string {
	switch {
	case 1 < len(s) && s[len(s)-2:] == "\r\n":
		s = s[:len(s)-2]
	case 0 < len(s) && s[len(s)-1:] == "\n":
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
