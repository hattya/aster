//
// aster :: aster_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hattya/aster"
	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/aster/internal/test"
	"github.com/hattya/go.cli"
	"github.com/hattya/go.notify"
)

func TestIgnore(t *testing.T) {
	at := &asterTest{
		src: ``,
		test: func(d time.Duration, cancel context.CancelFunc) {
			sh.Touch("a.go")
			sh.Mkdir("doc")

			sh.Mkdir("build", "html")
			sh.Touch("build", "html", "index.html")
			time.Sleep(d)

			src := `aster.ignore.push(/^build$/);`
			if err := test.Gen(src); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)
		},
		after: func(_ *aster.Aster, w *aster.Watcher) {
			if g, e := w.Paths(), []string{".", "doc"}; !reflect.DeepEqual(g, e) {
				t.Fatalf("expected %v, got %v", e, g)
			}
			for _, n := range []string{"a.go", "doc", "build"} {
				n = "." + string(os.PathSeparator) + n
				if err := w.Add(n); err != nil {
					t.Error(err)
				}
				if err := w.Remove(n); err != nil {
					t.Error(err)
				}
			}
			if g, e := w.Paths(), []string{"."}; !reflect.DeepEqual(g, e) {
				t.Fatalf("expected %v, got %v", e, g)
			}
		},
	}
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	if g, e := stderr, ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestInterrupt(t *testing.T) {
	at := &asterTest{
		src: ``,
		before: func(_ *aster.Aster, cancel context.CancelFunc) {
			cancel()
		},
		test: func(d time.Duration, cancel context.CancelFunc) {
			time.Sleep(d)

			cancel()

			time.Sleep(d)
		},
		after: func(_ *aster.Aster, w *aster.Watcher) {
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}
		},
	}
	if _, err := at.Run(); err != context.Canceled {
		t.Fatal(err)
	}

	at.before = nil
	if _, err := at.Run(); err != context.Canceled {
		t.Fatal(err)
	}
}

func TestNotify(t *testing.T) {
	n := test.NewNotifier()

	at := &asterTest{
		src: cli.Dedent(`
			aster.watch(/.+\.go$/, function() {
				aster.notify("success", "title", "text");
			});
		`),
		notifier: n,
		test: func(d time.Duration, _ context.CancelFunc) {
			sh.Touch("a.go")
			time.Sleep(d)
		},
	}
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	if g, e := stderr, ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestNotifyError(t *testing.T) {
	n := test.NewNotifier()
	at := &asterTest{
		src: cli.Dedent(`
			aster.watch(/.+\.go$/, function() {
				aster.notify("success", "title", "text");
			});
		`),
		notifier: n,
		test: func(d time.Duration, _ context.CancelFunc) {
			sh.Touch("a.go")
			time.Sleep(d)
		},
	}

	e := fmt.Errorf("Notify error")
	n.MockCall("Notify", e)
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(stderr, "\n", 2)
	if g, e := lines[0], fmt.Sprintf("aster.test: %v", e); g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestReload(t *testing.T) {
	at := &asterTest{
		src: ``,
		test: func(d time.Duration, _ context.CancelFunc) {
			if err := sh.Mkdir("build"); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)

			sh.Touch("a.go")

			src := `aster.ignore.push(/^build$/);`
			if err := test.Gen(src); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)

			src = `++;`
			if err := test.Gen(src); err != nil {
				t.Fatal(err)
			}
			time.Sleep(d)
		},
	}
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(stderr, "\n", 2)
	if g, e := lines[0], "aster.test: failed to reload"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestWatch(t *testing.T) {
	at := &asterTest{
		src: cli.Dedent(`
			aster.watch(/.+\.go$/, function(files) {
				cycles.push(files);
			});
		`),
		before: func(a *aster.Aster, _ context.CancelFunc) {
			sh.Mkdir(".git")
			a.Eval(`var cycles = [];`)
		},
		test: func(d time.Duration, _ context.CancelFunc) {
			sh.Touch("a.go")
			os.Remove("a.go")

			sh.Touch("b.go")
			os.Rename("b.go", "b_.go")
			os.Remove("b_.go")

			sh.Touch("c.go")
			sh.Mkdir("dir.go")
			time.Sleep(d)

			os.Chmod("c.go", 0755)
			os.Rename("c.go", "c_.go")

			os.Rename("dir.go", "dir_.go")

			sh.Touch(".git", "git.go")

			sh.Mkdir(".hg")
			sh.Touch(".hg", "hg.go")
			time.Sleep(d)
		},
		after: func(a *aster.Aster, _ *aster.Watcher) {
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
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	if g, e := stderr, ""; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

func TestWatchError(t *testing.T) {
	at := &asterTest{
		src: cli.Dedent(`
			aster.watch(/.+\.go$/, function() {
				throw new Error();
			});
		`),
		test: func(d time.Duration, _ context.CancelFunc) {
			sh.Touch("a.go")
			time.Sleep(d)
		},
	}
	stderr, err := at.Run()
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.SplitN(stderr, "\n", 2)
	if g, e := lines[0], "aster.test: Error"; g != e {
		t.Errorf("expected %q, got %q", e, g)
	}
}

type asterTest struct {
	src      string
	notifier notify.Notifier
	before   func(*aster.Aster, context.CancelFunc)
	test     func(time.Duration, context.CancelFunc)
	after    func(*aster.Aster, *aster.Watcher)
}

func (t *asterTest) Run() (string, error) {
	s := 101 * time.Millisecond
	var b bytes.Buffer
	app := cli.NewCLI()
	app.Name = "aster.test"
	app.Stderr = &b
	app.Action = func(*cli.Context) error {
		return test.Sandbox(func() error {
			if err := test.Gen(t.src); err != nil {
				return err
			}
			a, err := aster.New(app, t.notifier)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if t.before != nil {
				t.before(a, cancel)
			}

			w, err := aster.NewWatcher(ctx, a)
			if err != nil {
				return err
			}
			defer w.Close()

			w.Squash = s
			go w.Watch()
			t.test(time.Duration(s.Nanoseconds()*2), cancel)

			if t.after != nil {
				t.after(a, w)
			}
			return ctx.Err()
		})
	}

	err := app.Run(nil)
	return b.String(), err
}
