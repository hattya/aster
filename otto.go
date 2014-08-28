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
	"os"
	"os/exec"
	"strings"

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
	stdout := os.Stdout
	stderr := os.Stderr
	// args
	v, _ := call.Argument(0).Export()
	ary, ok := v.([]interface{})
	if !ok {
		return otto.UndefinedValue()
	}
	args := make([]string, len(ary))
	for i := range ary {
		s, ok := ary[i].(string)
		if !ok {
			return otto.UndefinedValue()
		}
		args[i] = s
	}
	// options
	v, _ = call.Argument(1).Export()
	options, ok := v.(map[string]interface{})
	if ok {
		var err error
		// stdout
		if v, ok := options["stdout"]; ok {
			if s, ok := v.(string); ok {
				stdout, err = os.Create(s)
				if err != nil {
					return throw(call.Otto.ToValue(err.Error()))
				}
				defer stdout.Close()
			}
		}
		// stderr
		if v, ok := options["stderr"]; ok {
			if s, ok := v.(string); ok {
				stderr, err = os.Create(s)
				if err != nil {
					return throw(call.Otto.ToValue(err.Error()))
				}
				defer stderr.Close()
			}
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
