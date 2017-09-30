//
// aster :: aster_test.go
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

package aster_test

import (
	"bytes"
	"context"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hattya/aster"
	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/aster/internal/test"
	"github.com/hattya/go.cli"
)

func TestWatch(t *testing.T) {
	tt := &asterTest{
		src: cli.Dedent(`
			aster.watch(/.+\.go$/, function(files) {
			  cycles.push(files);
			});
		`),
		before: func(a *aster.Aster) {
			sh.Mkdir(".git")
			a.Eval(`var cycles = [];`)
		},
		test: func(d time.Duration) {
			sh.Touch("a.go")
			os.Remove("a.go")

			sh.Touch("b.go")
			os.Rename("b.go", "b_.go")
			os.Remove("b_.go")

			sh.Touch("c.go")
			sh.Mkdir("dir.go")
			time.Sleep(d)

			os.Rename("c.go", "c_.go")
			os.Rename("dir.go", "dir_.go")

			sh.Touch(".git", "git.go")

			sh.Mkdir(".hg")
			sh.Touch(".hg", "hg.go")
			time.Sleep(d)
		},
		after: func(a *aster.Aster) {
			v, _ := a.Eval(`cycles;`)
			cycles := v.Object()
			// cycles.length
			v, _ = cycles.Get("length")
			n, _ := v.ToInteger()
			if g, e := n, int64(2); g != e {
				t.Fatalf("cycles.length = %v, expected %v", g, e)
			}
			// cycles[0].length
			v, _ = cycles.Get("0")
			files := v.Object()
			v, _ = files.Get("length")
			n, _ = v.ToInteger()
			if g, e := n, int64(1); g != e {
				t.Fatalf("cycles[0].length = %v, expected %v", g, e)
			}
			// cycles[0][0]
			v, _ = files.Get("0")
			s, _ := v.ToString()
			if g, e := s, "c.go"; g != e {
				t.Fatalf("cycles[0][0] = %q, expected %q", g, e)
			}
			// cycles[1].length
			v, _ = cycles.Get("1")
			files = v.Object()
			v, _ = files.Get("length")
			n, _ = v.ToInteger()
			if g, e := n, int64(1); g != e {
				t.Fatalf("cycles[1].length = %v, expected %v", g, e)
			}
			// cycles[1][0]
			v, _ = files.Get("0")
			s, _ = v.ToString()
			if g, e := s, "c_.go"; g != e {
				t.Fatalf("cycles[1][0] = %q, expected %q", g, e)
			}
		},
	}
	stderr, err := aster_(tt)
	if err != nil {
		t.Fatal(err)
	}
	if g, e := stderr, ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestReload(t *testing.T) {
	ci := os.Getenv("CI") != "" && runtime.GOOS == "linux"
	tt := &asterTest{
		src: ``,
		test: func(d time.Duration) {
			if err := sh.Mkdir("build"); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)

			if ci {
				os.Remove("Asterfile")
			}
			src := `aster.ignore.push(/^build$/);`
			if err := test.Gen(src); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)

			if ci {
				os.Remove("Asterfile")
			}
			src = `+;`
			if err := test.Gen(src); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)
		},
	}
	stderr, err := aster_(tt)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(stderr, "\n", 2)
	if g, e := lines[0], "aster: failed to reload"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

type asterTest struct {
	src    string
	before func(*aster.Aster)
	test   func(time.Duration)
	after  func(*aster.Aster)
}

func aster_(tt *asterTest) (string, error) {
	s := 101 * time.Millisecond
	var b bytes.Buffer
	app := cli.NewCLI()
	app.Stderr = &b
	app.Action = func(*cli.Context) error {
		return test.Sandbox(func() error {
			if err := test.Gen(tt.src); err != nil {
				return err
			}
			a, err := aster.New(app, "")
			if err != nil {
				return err
			}

			if tt.before != nil {
				tt.before(a)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			w, err := aster.NewWatcher(ctx, a)
			if err != nil {
				return err
			}
			w.Squash = s
			go w.Watch()
			tt.test(time.Duration(float64(s.Nanoseconds()) * 1.4))
			err = w.Close()

			if tt.after != nil {
				tt.after(a)
			}
			return err
		})
	}

	err := app.Run(nil)
	return b.String(), err
}
