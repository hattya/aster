//
// aster :: otto_test.go
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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/robertkrimen/otto"
)

func TestOS_Getwd(t *testing.T) {
	vm := newVM()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	src := `os.getwd();`
	s, err := testString(vm, src)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := s, wd; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestOS_Mkdir(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	vm := newVM()

	src := fmt.Sprintf(`os.mkdir('%v');`, path.Join(dir, "a"))
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = fmt.Sprintf(`os.mkdir('%v', 0777);`, path.Join(dir, "b"))
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	touch(path.Join(dir, "c"))
	src = fmt.Sprintf(`os.mkdir('%v', 0777);`, path.Join(dir, "c"))
	switch v, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !v:
		t.Errorf("expected true, got %v", v)
	}
}

func TestOS_Remove(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	vm := newVM()

	touch(path.Join(dir, "a"))
	src := fmt.Sprintf(`os.remove('%v');`, path.Join(dir, "a"))
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = fmt.Sprintf(`os.remove('%v');`, path.Join(dir, "b"))
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
}

func TestOS_Rename(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	vm := newVM()

	touch(path.Join(dir, "a"))
	src := fmt.Sprintf(`os.rename('%v', '%v');`, path.Join(dir, "a"), path.Join(dir, "b"))
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = fmt.Sprintf(`os.rename('%v', '%v');`, path.Join(dir, "a"), path.Join(dir, "c"))
	switch v, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !v:
		t.Errorf("expected true, got %v", v)
	}
}

var os_StatTests = []struct {
	path string
	dir  bool
}{
	{".", true},
	{"Asterfile", false},
}

func TestOS_Stat(t *testing.T) {
	vm := newVM()

	for _, tt := range os_StatTests {
		src := fmt.Sprintf(`os.stat('%v');`, tt.path)
		v, err := vm.Run(src)
		if err != nil {
			t.Fatal(err)
		}
		fi := v.Object()
		switch v, err := fi.Get("name"); {
		case err != nil:
			t.Error(err)
		case !v.IsString():
			t.Errorf("expected String, got %v", v)
		}
		switch v, err = fi.Get("size"); {
		case err != nil:
			t.Error(err)
		case !v.IsNumber():
			t.Errorf("expected Number, got %v", v)
		}
		switch v, err := fi.Get("mode"); {
		case err != nil:
			t.Error(err)
		case !v.IsNumber():
			t.Errorf("expected Number, got %v", v)
		}
		switch v, err := fi.Get("mtime"); {
		case err != nil:
			t.Error(err)
		case !v.IsObject() || v.Class() != "Date":
			t.Errorf("expected Date, got %v", v)
		}
		switch v, err := fi.Call("isDir"); {
		case err != nil:
			t.Error(err)
		case !v.IsBoolean():
			t.Errorf("expected Boolean, got %v", v)
		default:
			if g, _ := v.ToBoolean(); g != tt.dir {
				t.Errorf("expected %v, got %v", tt.dir, g)
			}
		}
		switch v, err := fi.Call("isRegular"); {
		case err != nil:
			t.Error(err)
		case !v.IsBoolean():
			t.Errorf("expected Boolean, got %v", v)
		default:
			if g, _ := v.ToBoolean(); g != !tt.dir {
				t.Errorf("expected %v, got %v", !tt.dir, g)
			}
		}
		switch v, err := fi.Call("perm"); {
		case err != nil:
			t.Error(err)
		case !v.IsNumber():
			t.Errorf("expected Number, got %v", v)
		}
	}

	src := `os.stat('__aster__');`
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
}

func TestOS_System(t *testing.T) {
	dir, err := tempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	// stdout
	stdout, err := os.Create(path.Join(dir, "stdout"))
	if err != nil {
		t.Fatal(err)
	}
	defer stdout.Close()
	// stderr
	stderr, err := os.Create(path.Join(dir, "stderr"))
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()
	// exe
	exe := path.Join(dir, "otto_cmd.exe")
	out, err := exec.Command("go", "build", "-o", exe, "otto_test_cmd.go").CombinedOutput()
	if err != nil {
		t.Fatalf("build failed\n%s", out)
	}

	vm := newVM()
	n1 := "'" + stdout.Name() + "'"
	n2 := "'" + stderr.Name() + "'"
	tmpl := `os.system(['%v', '-code', '%v'], {'stdout': %v, 'stderr': %v});`

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
		t.Errorf("stdout = %q, expected %q", g, e)
	}

	// String: stderr
	src = fmt.Sprintf(tmpl, exe, 1, n1, n2)
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
		t.Errorf("stderr = %q, expected %q", g, e)
	}

	// null: stdout
	src = fmt.Sprintf(tmpl, exe, 0, `null`, `null`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// null: stderr
	src = fmt.Sprintf(tmpl, exe, 1, `null`, `null`)
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true, got %v", b)
	}

	// Array: stdout
	ary, _ := vm.Object(`b = []`)
	src = fmt.Sprintf(tmpl, exe, 0, `b`, `b`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
	v, _ := ary.Get("length")
	n, _ := v.ToInteger()
	if g, e := n, int64(1); g != e {
		t.Errorf("stdout.length = %v, expected %v", g, e)
	}
	v, _ = ary.Get("0")
	s, _ := v.ToString()
	if g, e := s, "stdout"; g != e {
		t.Errorf("stdout = %q, expected %q", g, e)
	}

	// Array: stderr
	ary, _ = vm.Object(`b = []`)
	src = fmt.Sprintf(tmpl, exe, 1, `b`, `b`)
	switch b, err := testBoolean(vm, src); {
	case err != nil:
		t.Error(err)
	case !b:
		t.Errorf("expected true")
	}
	v, _ = ary.Get("length")
	n, _ = v.ToInteger()
	if g, e := n, int64(1); g != e {
		t.Errorf("stderr.length = %v, expected %v", g, e)
	}
	v, _ = ary.Get("0")
	s, _ = v.ToString()
	if g, e := s, "stderr"; g != e {
		t.Errorf("stderr = %q, expected %q", g, e)
	}

	// Object
	src = fmt.Sprintf(tmpl, exe, 0, `null`, `{}`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// dir
	src = fmt.Sprintf(`os.system(['./%v'], {'dir': '%v', 'stdout': null});`, path.Base(exe), dir)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// invalid args
	src = fmt.Sprintf(tmpl, exe, 1, `'.'`, n2)
	if _, err := vm.Run(src); err == nil {
		t.Error("expected error")
	}

	// invalid args
	src = fmt.Sprintf(tmpl, exe, 1, n1, `'.'`)
	if _, err := vm.Run(src); err == nil {
		t.Error("expected error")
	}

	// invalid args
	src = fmt.Sprintf(`os.system(['1']);`)
	if _, err := vm.Run(src); err == nil {
		t.Errorf("expected error")
	}

	// invalid args
	src = fmt.Sprintf(`os.system('1');`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	// invalid args
	src = fmt.Sprintf(`os.system([1]);`)
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}
}

func TestOS_Whence(t *testing.T) {
	vm := newVM()

	src := `os.whence();`
	if err := testUndefined(vm, src); err != nil {
		t.Error(err)
	}

	src = `os.whence('go');`
	if _, err := testString(vm, src); err != nil {
		t.Error(err)
	}
}

func TestBuffer(t *testing.T) {
	vm := otto.New()
	ary, _ := vm.Object(`ary = []`)
	b := &Buffer{
		vm:  vm,
		ary: ary,
	}
	b.Write([]byte("0\n1\n2"))
	b.Close()

	v, _ := ary.Get("length")
	n, _ := v.ToInteger()
	if g, e := n, int64(3); g != e {
		t.Errorf("ary.length = %v, expected %v", g, e)
	}
	for i := int64(0); i < n; i++ {
		e := strconv.FormatInt(i, 10)
		v, _ = ary.Get(e)
		g, _ := v.ToString()
		if g != e {
			t.Errorf("ary[%v] = %v, expected %v", i, g, e)
		}
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
