//
// aster :: aster_test.go
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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hattya/go.cli"
)

func TestWatch(t *testing.T) {
	tt := &asterTest{
		js: `aster.watch(/.+\.go$/, function(files) {
  cycles.push(files);
})`,
		before: func(af *Aster) {
			mkdir(".git")
			af.vm.Run(`var cycles = []`)
		},
		test: func(d time.Duration) {
			touch("a.go")
			os.Remove("a.go")

			touch("b.go")
			os.Rename("b.go", "b_.go")
			os.Remove("b_.go")

			touch("c.go")
			mkdir("dir.go")
			time.Sleep(d)

			os.Rename("c.go", "c_.go")
			os.Rename("dir.go", "dir_.go")

			touch(filepath.Join(".git", "git.go"))

			mkdir(".hg")
			touch(filepath.Join(".hg", "hg.go"))
			time.Sleep(d)
		},
		after: func(af *Aster) {
			v, _ := af.vm.Get("cycles")
			cycles := v.Object()
			// cycles.length
			v, _ = cycles.Get("length")
			n, _ := v.ToInteger()
			if g, e := n, int64(2); g != e {
				t.Fatalf("expected %v, got %v", e, g)
			}
			// cycles[0].length
			v, _ = cycles.Get("0")
			files := v.Object()
			v, _ = files.Get("length")
			n, _ = v.ToInteger()
			if g, e := n, int64(1); g != e {
				t.Fatalf("expected %v, got %v", e, g)
			}
			// cycles[0]
			v, _ = files.Get("0")
			s, _ := v.ToString()
			if g, e := s, "c.go"; g != e {
				t.Fatalf("expected %v, got %v", e, g)
			}
			// cyles[1].length
			v, _ = cycles.Get("1")
			files = v.Object()
			v, _ = files.Get("length")
			n, _ = v.ToInteger()
			if g, e := n, int64(1); g != e {
				t.Fatalf("expected %v, got %v", e, g)
			}
			// cycles[1]
			v, _ = files.Get("0")
			s, _ = v.ToString()
			if g, e := s, "c_.go"; g != e {
				t.Fatalf("expected %v, got %v", e, g)
			}
		},
	}
	stderr, err := aster(tt)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := stderr, ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestReload(t *testing.T) {
	tt := &asterTest{
		js: ``,
		test: func(d time.Duration) {
			js := ``
			if err := genAsterfile(js); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)

			js = `+`
			if err := genAsterfile(js); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)
		},
	}
	stderr, err := aster(tt)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(stderr, "\n", 2)
	if g, e := lines[0], "aster: failed to reload"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

type asterTest struct {
	js     string
	before func(*Aster)
	test   func(time.Duration)
	after  func(*Aster)
}

func aster(tt *asterTest) (string, error) {
	var b bytes.Buffer
	app.Stderr = &b
	app.Action = func(*cli.Context) error {
		return sandbox(func() error {
			if err := genAsterfile(tt.js); err != nil {
				return err
			}
			af, err := newAsterfile()
			if err != nil {
				return err
			}

			if tt.before != nil {
				tt.before(af)
			}

			watcher, err := newWatcher(af)
			if err != nil {
				return err
			}
			defer watcher.Close()

			go watcher.Watch()
			tt.test(149 * time.Millisecond)

			if tt.after != nil {
				tt.after(af)
			}
			return nil
		})
	}

	err := app.Run([]string{"-s", "101ms"})
	// reset
	app.Stderr = nil
	app.Flags.Set("s", app.Flags.Lookup("s").Default)

	return b.String(), err
}

func genAsterfile(js string) error {
	return ioutil.WriteFile("Asterfile", []byte(js), 0666)
}

func sandbox(test interface{}) error {
	dir, err := tempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	popd, err := pushd(dir)
	if err != nil {
		return err
	}
	defer popd()

	switch t := test.(type) {
	case func():
		t()
		return nil
	case func() error:
		return t()
	default:
		return fmt.Errorf("unknown type: %T", test)
	}
}

func pushd(path string) (func() error, error) {
	wd, err := os.Getwd()
	popd := func() error {
		if err == nil {
			return os.Chdir(wd)
		}
		return err
	}
	return popd, os.Chdir(path)
}

func mkdir(name string) error {
	return os.MkdirAll(name, 0777)
}

func touch(name string) error {
	return ioutil.WriteFile(name, []byte{}, 0666)
}

func tempDir() (string, error) {
	dir, err := ioutil.TempDir("", "aster.test")
	return filepath.ToSlash(dir), err
}
