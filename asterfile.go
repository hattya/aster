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
	"io/ioutil"

	"github.com/mattn/go-gntp"
	"github.com/robertkrimen/otto"
)

type Aster struct {
	vm    *otto.Otto
	gntp  *gntp.Client
	watch map[*otto.Object]*otto.Object
}

func newAsterfile() (*Aster, error) {
	a := &Aster{
		vm:   newVM(),
		gntp: newNotifier(),
	}
	// aster object
	aster, _ := a.vm.Object(`aster = {}`)
	aster.Set("watch", a.Watch)
	aster.Set("notify", a.Notify)
	// eval Asterfile
	src, err := ioutil.ReadFile("Asterfile")
	if err != nil {
		return nil, err
	}
	script, err := a.vm.Compile("Asterfile", src)
	if err != nil {
		return nil, ottoError(err)
	}
	_, err = a.vm.Run(script)
	if err != nil {
		return nil, ottoError(err)
	}
	return a, nil
}

func (a *Aster) Watch(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 2 {
		re := call.Argument(0)
		if re.Class() == "RegExp" {
			co := call.Argument(1)
			if co.Class() == "Function" {
				if a.watch == nil {
					a.watch = make(map[*otto.Object]*otto.Object)
				}
				a.watch[re.Object()] = co.Object()
			}
		}
	}
	return otto.UndefinedValue()
}

func (a *Aster) Notify(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 3 {
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
	for re, co := range a.watch {
		var changes []interface{}
		for n := range files {
			v, _ := re.Call("test", n)
			test, _ := v.ToBoolean()
			if test {
				changes = append(changes, n)
			}
		}
		if 0 < len(changes) {
			ary, _ := a.vm.Call(`new Array`, nil, changes...)
			_, err := co.Call("call", nil, ary)
			if err != nil {
				warn(err)
			}
			return
		}
	}
}
