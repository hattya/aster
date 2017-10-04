//
// aster :: util.go
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
