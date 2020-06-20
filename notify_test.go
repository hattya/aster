//
// aster :: notify_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
