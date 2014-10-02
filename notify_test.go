//
// aster :: notify_test.go
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

import "testing"

func TestGNTPValue(t *testing.T) {
	var v GNTPValue
	if g, e := v.IsBoolFlag(), true; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}
	if g, e := v.String(), ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}

	for _, tt := range [][]string{
		{"", ""},
		{"true", "localhost:23053"},
		{"false", ""},
		{"localhost", "localhost:23053"},
		{"localhost:23053", "localhost:23053"},
	} {
		var v GNTPValue
		v.Set(tt[0])
		if g, e := v.String(), tt[1]; g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
		if g, e := v.Get(), tt[1]; g != e {
			t.Errorf("expected %q, got %q", e, g)
		}
	}
}
