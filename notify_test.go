//
// aster :: notify_test.go
//
//   Copyright (c) 2014-2018 Akinori Hattori <hattya@gmail.com>
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

package aster_test

import (
	"testing"

	"github.com/hattya/aster"
)

var gntpValueTests = []struct {
	in, out string
}{
	{"", ""},
	{"true", "localhost:23053"},
	{"false", ""},
	{"localhost", "localhost:23053"},
	{"localhost:23053", "localhost:23053"},
}

func TestGNTPValue(t *testing.T) {
	var v aster.GNTPValue
	if !v.IsBoolFlag() {
		t.Error("GNTPValue.IsBoolFlag() = false, expected true")
	}
	if g, e := v.String(), ""; g != e {
		t.Errorf("GNTPValue.String() = %q, expected %q", g, e)
	}

	for _, tt := range gntpValueTests {
		var v aster.GNTPValue
		v.Set(tt.in)
		if g, e := v.String(), tt.out; g != e {
			t.Errorf("GNTPValue.String() = %q, expected %q", g, e)
		}
		if g, e := v.Get(), tt.out; g != e {
			t.Errorf("GNTPValue.Get() = %q, expected %q", g, e)
		}
	}
}
