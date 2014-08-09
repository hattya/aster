//
// aster :: otto_test.go
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
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/robertkrimen/otto"
)

func TestOSSystem(t *testing.T) {
	prog := filepath.ToSlash(os.Args[0])
	vm := newVM()

	src := fmt.Sprintf(`os.system("%s", "-test.run=TestProcess", "0")`, prog)
	switch v, err := vm.Run(src); {
	case err != nil:
		t.Error(err)
	case !v.IsUndefined():
		t.Errorf("expected undefined, got %v", v)
	}

	src = fmt.Sprintf(`os.system("%s", "-test.run=TestProcess", "1")`, prog)
	if v, err := vm.Run(src); err != nil {
		t.Error(err)
	} else if err = isTrue(v); err != nil {
		t.Error(err)
	}

	src = fmt.Sprintf(`os.system(1)`)
	if _, err := vm.Run(src); err == nil {
		t.Errorf("expected error")
	}
}

func isTrue(v otto.Value) error {
	if !v.IsBoolean() {
		return fmt.Errorf("expected boolean, got %v", v)
	}
	switch b, err := v.ToBoolean(); {
	case err != nil:
		return err
	case !b:
		return fmt.Errorf("expected true, got %v", b)
	}
	return nil
}

func TestProcess(*testing.T) {
	if len(os.Args) != 3 || os.Args[1] != "-test.run=TestProcess" {
		return
	}
	n, _ := strconv.Atoi(os.Args[2])
	os.Exit(n)
}
