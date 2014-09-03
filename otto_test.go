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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

func TestOSSystem(t *testing.T) {
	dir, err := mkdtemp()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	// stdout
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	// stderr
	stderr, err := os.Create(filepath.Join(dir, "stderr"))
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	// exe
	exe := filepath.ToSlash(filepath.Join(dir, "otto_cmd.exe"))
	out, err := exec.Command("go", "build", "-o", exe, "otto_test_cmd.go").CombinedOutput()
	if err != nil {
		t.Fatalf("build failed\n%s", out)
	}

	vm := newVM()
	n1 := "'" + filepath.ToSlash(stdout.Name()) + "'"
	n2 := "'" + filepath.ToSlash(stderr.Name()) + "'"
	tmpl := `os.system(['%s', '-code', '%d'], {'stdout': %s, 'stderr': %s})`

	// String: stdout
	src := fmt.Sprintf(tmpl, exe, 0, n1, n2)
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
		t.Errorf("expected %q, got %q", e, g)
	}

	// String: stderr
	src = fmt.Sprintf(tmpl, exe, 1, n1, n2)
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true")
	}
	stderr.Seek(0, os.SEEK_SET)
	data, err = ioutil.ReadAll(stderr)
	if err != nil {
		t.Fatal(err)
	}
	lines = strings.SplitN(string(data), "\n", 2)
	if g, e := lines[0], "stderr"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}

	// null: stdout
	src = fmt.Sprintf(tmpl, exe, 0, "null", "null")
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// null: stderr
	src = fmt.Sprintf(tmpl, exe, 1, "null", "null")
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true")
	}

	// Array: stdout
	ary, _ := vm.Object(`b = []`)
	src = fmt.Sprintf(tmpl, exe, 0, "b", "b")
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
	v, _ := ary.Get("length")
	n, _ := v.ToInteger()
	if g, e := n, int64(1); g != e {
		t.Errorf("expected %d, got %d", e, g)
	}
	v, _ = ary.Get(strconv.FormatInt(0, 10))
	s, _ := v.ToString()
	if g, e := s, "stdout"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}

	// Array: stderr
	ary, _ = vm.Object(`b = []`)
	src = fmt.Sprintf(tmpl, exe, 1, "b", "b")
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true")
	}
	v, _ = ary.Get("length")
	n, _ = v.ToInteger()
	if g, e := n, int64(1); g != e {
		t.Errorf("expected %d, got %d", e, g)
	}
	v, _ = ary.Get(strconv.FormatInt(0, 10))
	s, _ = v.ToString()
	if g, e := s, "stderr"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}

	// invalid args
	src = fmt.Sprintf(tmpl, exe, 1, ".", n2)
	if _, err := vm.Run(src); err == nil {
		t.Error("expected error")
	}

	// invalid args
	src = fmt.Sprintf(tmpl, exe, 1, n1, ".")
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
			err = fmt.Errorf("expected Boolean, got %v", v)
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
			err = fmt.Errorf("expected String, got %v", v)
		} else {
			return v.ToString()
		}
	}
	return "", err
}
