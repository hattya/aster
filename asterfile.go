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

package main

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/mattn/go-gntp"
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
	vm      *otto.Otto
	n       int32
	watches []*watchExpr
	ignore  AsterIgnore
	gntp    *gntp.Client
}

func newAsterfile() (*Aster, error) {
	a := &Aster{
		gntp: newNotifier(),
	}
	if err := a.Eval(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Aster) Eval() error {
	a.vm = newVM()
	a.watches = nil
	// aster object
	aster, _ := a.vm.Object(fmt.Sprintf(`
		aster = {
		  arch:   %q,
		  os:     %q,
		  ignore: [/%v/],
		}
	`, runtime.GOARCH, runtime.GOOS, defaultIgnore))
	aster.Set("watch", a.Watch)
	aster.Set("notify", a.Notify)
	aster.Set("title", a.Title)
	// watch Asterfile
	rx, _ := a.vm.Call(`new RegExp`, nil, `^Asterfile$`)
	fn, _ := a.vm.ToValue(a.Reload)
	aster.Call("watch", rx, fn)
	// eval Asterfile
	script, err := a.vm.Compile("Asterfile", nil)
	if err != nil {
		return ottoError(err)
	}
	_, err = a.vm.Run(script)
	if err != nil {
		return ottoError(err)
	}
	// update ignore
	ary, _ := a.vm.Object(`aster.ignore`)
	if ary.Class() == "Array" {
		// get aster.ignore.length
		v, _ := ary.Get("length")
		n, _ := v.ToInteger()

		a.ignore = make(AsterIgnore, n)
		var i int64
		for j := int64(0); j < n; j++ {
			v, _ := ary.Get(strconv.FormatInt(j, 10))
			if v.Class() == "RegExp" {
				a.ignore[i] = v.Object()
				i++
			}
		}
		a.ignore = a.ignore[:i]
	}
	return nil
}

func (a *Aster) Watch(call otto.FunctionCall) otto.Value {
	rx := call.Argument(0)
	if rx.Class() == "RegExp" {
		fn := call.Argument(1)
		if fn.Class() == "Function" {
			a.watches = append(a.watches, &watchExpr{
				rx: rx.Object(),
				fn: fn.Object(),
			})
		}
	}
	return otto.UndefinedValue()
}

func (a *Aster) Notify(call otto.FunctionCall) otto.Value {
	if 3 <= len(call.ArgumentList) {
		var args [3]string
		for i := range args {
			args[i], _ = call.ArgumentList[i].ToString()
		}
		notify(a.gntp, args[0], args[1], args[2])
	}
	return otto.UndefinedValue()
}

func (a *Aster) Title(call otto.FunctionCall) otto.Value {
	if 1 <= len(call.ArgumentList) {
		title, _ := call.ArgumentList[0].ToString()
		app.Title(title)
	}
	return otto.UndefinedValue()
}

func (a *Aster) Reload(otto.FunctionCall) otto.Value {
	// create snapshot
	vm := a.vm
	watches := a.watches
	ignore := a.ignore
	// eval
	var name, text otto.Value
	if err := a.Eval(); err != nil {
		// report error
		warn("failed to reload\n", err)
		// rollback to snapshot
		a.vm = vm
		a.watches = watches
		a.ignore = ignore

		name, _ = a.vm.ToValue("failure")
		text, _ = a.vm.ToValue("Error occurred while reloading Asterfile")
	} else {
		atomic.AddInt32(&a.n, 1)

		name, _ = a.vm.ToValue("success")
		text, _ = a.vm.ToValue("Asterfile has been reloaded")
	}
	title, _ := a.vm.ToValue("Aster reload")
	// call aster.notify
	aster, _ := a.vm.Object(`aster`)
	v, _ := aster.Call("notify", name, title, text)
	return v
}

func (a *Aster) Ignore(name string) bool {
	return a.ignore.Match(name)
}

func (a *Aster) Reloaded() bool {
	return 0 < atomic.SwapInt32(&a.n, 0)
}

func (a *Aster) OnChange(files map[string]int) {
	for _, w := range a.watches {
		// call RegExp.test
		var cl []interface{}
		for n := range files {
			v, _ := w.rx.Call("test", n)
			if b, _ := v.ToBoolean(); b {
				cl = append(cl, n)
			}
		}
		// call callback
		if 0 < len(cl) {
			ary, _ := a.vm.Call(`new Array`, nil, cl...)
			_, err := w.fn.Call("call", nil, ary)
			if err != nil {
				warn(err)
			}
			return
		}
	}
}

type watchExpr struct {
	rx *otto.Object // RegExp
	fn *otto.Object // Function
}

type AsterIgnore []*otto.Object // []RegExp

func (ig AsterIgnore) Match(s string) bool {
	for _, rx := range []*otto.Object(ig) {
		v, _ := rx.Call("test", s)
		test, _ := v.ToBoolean()
		if test {
			return true
		}
	}
	return false
}
