//
// aster :: otto.go
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
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/hattya/go.binfmt"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
)

func newVM() *otto.Otto {
	vm := otto.New()
	// os object
	os, _ := vm.Object(`os = {}`)
	os.Set("system", os_system)
	os.Set("whence", os_whence)

	return vm
}

func throw(value otto.Value, _ error) otto.Value {
	panic(value)
}

func ottoError(err error) error {
	switch e := err.(type) {
	case *otto.Error:
		return fmt.Errorf(strings.TrimSpace(e.String()))
	case parser.ErrorList:
		var b bytes.Buffer
		for i, pe := range e {
			if 0 < i {
				fmt.Fprintln(&b)
			}
			fmt.Fprintf(&b, "Asterfile:%d:%d: %s", pe.Position.Line, pe.Position.Column, pe.Message)
			if pe.Message == "Unexpected end of input" {
				break
			}
		}
		return fmt.Errorf(b.String())
	}
	return err
}

func os_system(call otto.FunctionCall) otto.Value {
	// defaults
	var stdout io.WriteCloser = os.Stdout
	var stderr io.WriteCloser = os.Stderr
	// args
	v := call.Argument(0)
	if v.Class() != "Array" {
		return otto.UndefinedValue()
	}
	ary := v.Object()
	v, _ = ary.Get("length")
	n, _ := v.ToInteger()
	args := make([]string, n)
	for i := int64(0); i < n; i++ {
		v, _ := ary.Get(strconv.FormatInt(i, 10))
		if !v.IsString() {
			return otto.UndefinedValue()
		}
		args[i], _ = v.ToString()
	}
	// options
	v = call.Argument(1)
	if v.Class() == "Object" {
		options := v.Object()
		redir := func(o *otto.Object, k string) (io.WriteCloser, error) {
			switch v, _ = o.Get(k); {
			case v.IsString():
				s, _ := v.ToString()
				return os.Create(s)
			case v.IsNull():
				return discard, nil
			case v.Class() == "Array":
				return newBuffer(call.Otto, v.Object())
			}
			return nil, nil
		}
		// stdout
		switch wc, err := redir(options, "stdout"); {
		case err != nil:
			return throw(call.Otto.ToValue(err.Error()))
		case wc != nil:
			stdout = wc
			defer wc.Close()
		}
		// stderr
		switch wc, err := redir(options, "stderr"); {
		case err != nil:
			return throw(call.Otto.ToValue(err.Error()))
		case wc != nil:
			stderr = wc
			defer wc.Close()
		}
	}

	cmd := binfmt.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return otto.TrueValue()
		}
		return throw(call.Otto.ToValue(err.Error()))
	}
	return otto.UndefinedValue()
}

func os_whence(call otto.FunctionCall) otto.Value {
	if 1 <= len(call.ArgumentList) {
		s, _ := call.ArgumentList[0].ToString()
		path, err := exec.LookPath(s)
		if err == nil {
			v, _ := call.Otto.ToValue(path)
			return v
		}
	}
	return otto.UndefinedValue()
}

type Buffer struct {
	vm  *otto.Otto
	ary *otto.Object
	mu  sync.Mutex
	b   bytes.Buffer
}

func newBuffer(vm *otto.Otto, o *otto.Object) (b io.WriteCloser, err error) {
	if o.Class() == "Array" {
		b = &Buffer{
			vm:  vm,
			ary: o,
		}
	} else {
		err = fmt.Errorf("instance is %q", o.Class())
	}
	return
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.b.Write(p)
	for {
		if s, err := b.b.ReadString('\n'); err == nil {
			v, _ := b.vm.ToValue(s[:len(s)-1])
			b.ary.Call("push", v)
		} else {
			b.b.WriteString(s)
			break
		}
	}
	return len(p), nil
}

func (b *Buffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if 0 < b.b.Len() {
		v, _ := b.vm.ToValue(b.b.String())
		b.ary.Call("push", v)
	}
	return nil
}
