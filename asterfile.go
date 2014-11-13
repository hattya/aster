//
// aster :: asterfile.go
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
	"regexp"
	"strconv"

	"github.com/mattn/go-gntp"
	"github.com/robertkrimen/otto"
)

var defaultIgnore string

func init() {
	var b bytes.Buffer
	b.WriteString(`^(?:`)
	for i, s := range []string{
		".hg",
		".git",
		".svn",
		".bzr",
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
	vm     *otto.Otto
	watch  []*AsterWatch
	gntp   *gntp.Client
	ignore AsterIgnore
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
	a.watch = nil
	// aster object
	aster, _ := a.vm.Object(fmt.Sprintf(`aster = {
	  'ignore': [
	    /%v/,
	  ],
	}`, defaultIgnore))
	aster.Set("watch", a.Watch)
	aster.Set("notify", a.Notify)
	aster.Set("title", a.Title)
	// watch Asterfile
	re, _ := a.vm.Call(`new RegExp`, nil, `^Asterfile$`)
	cb, _ := a.vm.ToValue(a.Reload)
	aster.Call(`watch`, re, cb)
	// eval Asterfile
	script, err := a.vm.Compile("Asterfile", nil)
	if err != nil {
		return ottoError(err)
	}
	_, err = a.vm.Run(script)
	if err != nil {
		return ottoError(err)
	}
	a.Ignore()
	return nil
}

func (a *Aster) Ignore() {
	ary, _ := a.vm.Object(`aster.ignore`)
	if ary.Class() == "Array" {
		// get aster.ignore.length
		v, _ := ary.Get("length")
		n, _ := v.ToInteger()

		a.ignore = make(AsterIgnore, n)
		i := 0
		for j := int64(0); j < n; j++ {
			v, _ := ary.Get(strconv.FormatInt(j, 10))
			if v.Class() == "RegExp" {
				a.ignore[i] = v.Object()
				i++
			}
		}
		a.ignore = a.ignore[:i]
	}
}

func (a *Aster) Reload(otto.FunctionCall) otto.Value {
	// create snapshot
	vm := a.vm
	watch := a.watch
	ignore := a.ignore
	// eval
	var name, text otto.Value
	if err := a.Eval(); err != nil {
		// report error
		warn("failed to reload\n", err)
		// rollback to snapshot
		a.vm = vm
		a.watch = watch
		a.ignore = ignore

		name, _ = a.vm.ToValue("failure")
		text, _ = a.vm.ToValue("Error occurred while reloading Asterfile")
	} else {
		name, _ = a.vm.ToValue("success")
		text, _ = a.vm.ToValue("Asterfile has been reloaded")
	}
	title, _ := a.vm.ToValue("Aster reload")
	// call aster.notify
	aster, _ := a.vm.Object(`aster`)
	v, _ := aster.Call(`notify`, name, title, text)
	return v
}

func (a *Aster) Watch(call otto.FunctionCall) otto.Value {
	re := call.Argument(0)
	if re.Class() == "RegExp" {
		cb := call.Argument(1)
		if cb.Class() == "Function" {
			a.watch = append(a.watch, &AsterWatch{re.Object(), cb.Object()})
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

func (a *Aster) OnChange(files map[string]int) {
	for _, w := range a.watch {
		// call RegExp.test
		var cl []interface{}
		for n := range files {
			v, _ := w.re.Call(`test`, n)
			test, _ := v.ToBoolean()
			if test {
				cl = append(cl, n)
			}
		}
		// call callback
		if 0 < len(cl) {
			ary, _ := a.vm.Call(`new Array`, nil, cl...)
			_, err := w.cb.Call(`call`, nil, ary)
			if err != nil {
				warn(err)
			}
			return
		}
	}
}

type AsterWatch struct {
	re *otto.Object // RegExp
	cb *otto.Object // Function
}

type AsterIgnore []*otto.Object // []RegExp

func (ig AsterIgnore) Match(s string) bool {
	for _, re := range []*otto.Object(ig) {
		v, _ := re.Call(`test`, s)
		test, _ := v.ToBoolean()
		if test {
			return true
		}
	}
	return false
}
