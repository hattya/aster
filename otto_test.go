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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

func TestOSSystem(t *testing.T) {
	// stdout
	stdout, err := mktemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(stdout.Name())
	defer stdout.Close()
	// stderr
	stderr, err := mktemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(stderr.Name())
	defer stderr.Close()

	vm := newVM()
	n1 := filepath.ToSlash(stdout.Name())
	n2 := filepath.ToSlash(stderr.Name())
	tmpl := `os.system(['go', 'run', 'otto_test_cmd.go', '-code', '%d'], {'stdout': '%s', 'stderr': '%s'})`

	src := fmt.Sprintf(tmpl, 0, n1, n2)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
	stdout.Seek(0, os.SEEK_SET)
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(string(data), "\n", 2)
	if g, e := lines[0], "stdout"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	src = fmt.Sprintf(tmpl, 1, n1, n2)
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true, got %v", b)
	}
	stderr.Seek(0, os.SEEK_SET)
	data, err = ioutil.ReadAll(stderr)
	if err != nil {
		t.Fatal(err)
	}
	lines = strings.SplitN(string(data), "\n", 2)
	if g, e := lines[0], "stderr"; g != e {
		t.Errorf("expected %v, got %v", e, g)
	}

	// invalid args
	src = fmt.Sprintf(tmpl, 1, ".", n2)
	if _, err := vm.Run(src); err == nil {
		t.Error("expected error")
	}

	// invalid args
	src = fmt.Sprintf(tmpl, 1, n1, ".")
	if _, err := vm.Run(src); err == nil {
		t.Error("expected error")
	}

	// invalid args
	src = fmt.Sprintf(`os.system(['1'])`)
	if _, err := vm.Run(src); err == nil {
		t.Errorf("expected error")
	}

	// invalid args
	src = fmt.Sprintf(`os.system('1')`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// invalid args
	src = fmt.Sprintf(`os.system([1])`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
}

func TestOSWhence(t *testing.T) {
	vm := newVM()

	src := `os.whence()`
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = `os.whence('go')`
	if _, err := testString(vm, src); err != nil {
		t.Error(err)
	}
}

func testUndefined(vm *otto.Otto, src string) error {
	v, err := vm.Run(src)
	if err == nil && !v.IsUndefined() {
		err = fmt.Errorf("expected undefined, got %v", v)
	}
	return err
}

func testBoolean(vm *otto.Otto, src string) (bool, error) {
	v, err := vm.Run(src)
	if err == nil {
		if !v.IsBoolean() {
			err = fmt.Errorf("expected boolean, got %v", v)
		} else {
			return v.ToBoolean()
		}
	}
	return false, err
}

func testString(vm *otto.Otto, src string) (string, error) {
	v, err := vm.Run(src)
	if err == nil {
		if !v.IsString() {
			err = fmt.Errorf("expected string, got %v", v)
		} else {
			return v.ToString()
		}
	}
	return "", err
}
