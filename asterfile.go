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
	"fmt"

	"github.com/mattn/go-gntp"
	"github.com/robertkrimen/otto"
)

type Aster struct {
	vm    *otto.Otto
	gntp  *gntp.Client
	watch []*AsterWatch
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
	aster, _ := a.vm.Object(`aster = {}`)
	aster.Set("watch", a.Watch)
	aster.Set("notify", a.Notify)
	// watch Asterfile
	re, _ := a.vm.Call(`new RegExp`, nil, `^Asterfile$`)
	cb, _ := a.vm.ToValue(a.Reload)
	aster.Call("watch", re, cb)
	// eval Asterfile
	script, err := a.vm.Compile("Asterfile", nil)
	if err != nil {
		return ottoError(err)
	}
	_, err = a.vm.Run(script)
	if err != nil {
		return ottoError(err)
	}
	return nil
}

func (a *Aster) Reload(otto.FunctionCall) otto.Value {
	// create snapshot
	vm := a.vm
	watch := a.watch
	// eval
	var name string
	if err := a.Eval(); err != nil {
		name = "failure"
		// report error
		warn("failed to reload")
		fmt.Fprintln(stderr, err)
		// rollback to snapshot
		a.vm = vm
		a.watch = watch
	} else {
		name = "success"
	}
	a.vm.Run(fmt.Sprintf(`aster.notify("%s", "Aster reload", "Asterfile has been reloaded")`, name))
	return otto.UndefinedValue()
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
		var args []string
		for _, a := range call.ArgumentList {
			s, _ := a.ToString()
			args = append(args, s)
		}
		notify(a.gntp, args[0], args[1], args[2])
	}
	return otto.UndefinedValue()
}

func (a *Aster) OnChange(files map[string]int) {
	for _, w := range a.watch {
		// call RegExp.test
		var cl []interface{}
		for n := range files {
			v, _ := w.re.Call("test", n)
			test, _ := v.ToBoolean()
			if test {
				cl = append(cl, n)
			}
		}
		// call callback
		if 0 < len(cl) {
			ary, _ := a.vm.Call(`new Array`, nil, cl...)
			_, err := w.cb.Call("call", nil, ary)
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
