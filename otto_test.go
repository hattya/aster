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
	"strconv"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

func TestOSSystem(t *testing.T) {
	prog := filepath.ToSlash(os.Args[0])
	vm := newVM()
	tmpl := `os.system(["%s", "-test.run=TestProcess", "%d"])`

	src := fmt.Sprintf(tmpl, prog, 0)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = fmt.Sprintf(tmpl, prog, 1)
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true, got %v", b)
	}

	// invalid args
	src = fmt.Sprintf(`os.system(["1"])`)
	if _, err := vm.Run(src); err == nil {
		t.Errorf("expected error")
	}

	// invalid args
	src = fmt.Sprintf(`os.system("1")`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// invalid args
	src = fmt.Sprintf(`os.system([1])`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
}

func TestOSSystemRedir(t *testing.T) {
	f, err := mktemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	prog := filepath.ToSlash(os.Args[0])
	name := filepath.ToSlash(f.Name())
	vm := newVM()
	tmpl := `os.system(["%s", "-test.run=TestRedirection", "%[2]s"], {"%[2]s": "%s"})`

	for _, redir := range []string{"stdout", "stderr"} {
		src := fmt.Sprintf(tmpl, prog, redir, name)
		if err := testUndefined(vm, src); err != nil {
			t.Error(err)
		}
		f.Seek(0, os.SEEK_SET)
		data, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		list := strings.SplitN(string(data), "\n", 2)
		if g, e := list[0], redir; g != e {
			t.Errorf("expected %v, got %v", e, g)
		}

		// invalid args
		src = fmt.Sprintf(tmpl, prog, redir, ".")
		if _, err := vm.Run(src); err == nil {
			t.Error("expected error")
		}
	}
}

func TestOSWhence(t *testing.T) {
	vm := newVM()

	src := `os.whence()`
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = `os.whence("go")`
	if _, err := testString(vm, src); err != nil {
		t.Error(err)
	}
}

func TestProcess(*testing.T) {
	if len(os.Args) != 3 || os.Args[1] != "-test.run=TestProcess" {
		return
	}
	n, _ := strconv.Atoi(os.Args[2])
	os.Exit(n)
}

func TestRedirection(*testing.T) {
	if len(os.Args) != 3 || os.Args[1] != "-test.run=TestRedirection" {
		return
	}
	switch os.Args[2] {
	case "stdout":
		fmt.Fprintln(os.Stdout, "stdout")
	case "stderr":
		fmt.Fprintln(os.Stderr, "stderr")
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
			err = fmt.Errorf("expected string, got %v")
		} else {
			return v.ToString()
		}
	}
	return "", err
}
