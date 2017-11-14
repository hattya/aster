//
// aster :: asterfile.go
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

package aster

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/hattya/go.cli"
	"github.com/hattya/go.notify"
	"github.com/robertkrimen/otto"
)

var defaultIgnore string

func init() {
	var b bytes.Buffer
	b.WriteString(`^(?:`)
	for i, s := range []string{
		".bzr",
		".git",
		".hg",
		".svn",
		"CVS",
	} {
		if 0 < i {
			b.WriteRune('|')
		}
		b.WriteString(regexp.QuoteMeta(s))
	}
	b.WriteString(`)$`)
	defaultIgnore = b.String()
}

type Aster struct {
	ui *cli.CLI
	i  int32
	n  notify.Notifier

	mu      sync.Mutex
	vm      *otto.Otto
	watches []*watch
}

func New(ui *cli.CLI, n notify.Notifier) (*Aster, error) {
	if err := events(n); err != nil {
		return nil, err
	}
	a := &Aster{
		ui: ui,
		n:  n,
	}
	if err := a.eval(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Aster) eval() error {
	a.vm = newVM()
	a.watches = nil
	// aster object
	aster, _ := a.vm.Object(fmt.Sprintf(`
		aster = {
		  arch: %q,
		  ignore: [/%v/],
		  os: %q,
		}
	`, runtime.GOARCH, defaultIgnore, runtime.GOOS))
	aster.Set("notify", a.notify)
	aster.Set("title", a.title)
	aster.Set("watch", a.watch)
	// watch Asterfile
	rx, _ := a.vm.Call(`new RegExp`, nil, `^Asterfile$`)
	aster.Call("watch", rx, a.reload)
	// eval Asterfile
	script, err := a.vm.Compile("Asterfile", nil)
	if err != nil {
		return ottoError(err)
	}
	_, err = a.vm.Run(script)
	return ottoError(err)
}

func (a *Aster) watch(call otto.FunctionCall) otto.Value {
	rx := call.Argument(0)
	if rx.Class() == "RegExp" {
		fn := call.Argument(1)
		if fn.Class() == "Function" {
			a.watches = append(a.watches, &watch{
				rx: rx.Object(),
				fn: fn.Object(),
			})
		}
	}
	return otto.UndefinedValue()
}

func (a *Aster) notify(call otto.FunctionCall) otto.Value {
	if a.n == nil || len(call.ArgumentList) < 3 {
		return otto.UndefinedValue()
	}

	var args [3]string
	for i := range args {
		args[i], _ = call.ArgumentList[i].ToString()
	}
	if err := a.n.Notify(args[0], args[1], args[2]); err != nil {
		warn(a.ui, err)
	}
	return otto.UndefinedValue()
}

func (a *Aster) title(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	title, _ := call.ArgumentList[0].ToString()
	a.ui.Title(title)
	return otto.UndefinedValue()
}

func (a *Aster) reload(otto.FunctionCall) otto.Value {
	// create snapshot
	vm := a.vm
	watches := a.watches
	// eval
	var name, text string
	if err := a.eval(); err != nil {
		// report error
		warn(a.ui, "failed to reload\n", err)
		// rollback to snapshot
		a.vm = vm
		a.watches = watches

		name = "failure"
		text = "Error occurred while reloading Asterfile"
	} else {
		atomic.AddInt32(&a.i, 1)

		name = "success"
		text = "Asterfile has been reloaded"
	}
	title := "Aster reload"
	// call aster.notify
	v, _ := a.vm.Call(`aster.notify`, nil, name, title, text)
	return v
}

func (a *Aster) Eval(src interface{}) (otto.Value, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.vm.Run(src)
}

func (a *Aster) Ignore(name string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	ary, _ := a.vm.Object(`aster.ignore`)
	if ary.Class() == "Array" {
		// aster.ignore.length
		v, _ := ary.Get("length")
		n, _ := v.ToInteger()

		for i := int64(0); i < n; i++ {
			v, _ := ary.Get(strconv.FormatInt(i, 10))
			if v.Class() == "RegExp" {
				v, _ = v.Object().Call("test", name)
				if b, _ := v.ToBoolean(); b {
					return true
				}
			}
		}
	}
	return false
}

func (a *Aster) Reloaded() bool {
	return 0 < atomic.SwapInt32(&a.i, 0)
}

func (a *Aster) OnChange(ctx context.Context, files map[string]int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	i := atomic.LoadInt32(&a.i)
L:
	for _, w := range a.watches {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// call RegExp.test
		var cl []interface{}
		for n := range files {
			v, _ := w.rx.Call("test", n)
			if b, _ := v.ToBoolean(); b {
				cl = append(cl, n)
				delete(files, n)
			}
		}
		// call Function.call
		if 0 < len(cl) {
			ary, _ := a.vm.Call(`new Array`, nil, cl...)
			_, err := w.fn.Call("call", nil, ary)
			if err != nil {
				warn(a.ui, ottoError(err))
			}
		}

		if len(files) == 0 {
			break
		} else if v := atomic.LoadInt32(&a.i); i != v {
			// Asterfile has been reloaded
			i = v
			goto L
		}
	}
}

type watch struct {
	rx *otto.Object // RegExp
	fn *otto.Object // Function
}
