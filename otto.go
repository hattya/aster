//
// aster :: otto.go
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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hattya/go.binfmt"
	"github.com/hattya/go.cli"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
)

func newVM() *otto.Otto {
	vm := otto.New()
	// os object
	os_, _ := vm.Object(`os = {}`)
	os_.Set("getwd", os_getwd)
	os_.Set("mkdir", os_mkdir)
	os_.Set("open", os_open)
	os_.Set("remove", os_remove)
	os_.Set("rename", os_rename)
	os_.Set("stat", os_stat)
	os_.Set("system", os_system)
	os_.Set("whence", os_whence)

	script, _ := vm.Compile("os", fmt.Sprintf(cli.Dedent(`
		os.File = function(impl) {
		  this._impl = impl;
		};

		os.File.prototype.close = function() {
		  return this._impl.Close.apply(this, arguments);
		};

		os.File.prototype.name = function() {
		  return this._impl.Name.apply(this, arguments);
		};

		os.File.prototype.read = function() {
		  return this._impl.Read.apply(this, arguments);
		};

		os.File.prototype.readLine = function() {
		  return this._impl.ReadLine.apply(this, arguments);
		};

		os.File.prototype.write = function() {
		  return this._impl.Write.apply(this, arguments);
		};

		os.FileInfo = function(name, size, mode, mtime) {
		  this.name = name;
		  this.size = size;
		  this.mode = mode;
		  this.mtime = mtime;
		};

		os.FileInfo.prototype.isDir = function() {
		  return (this.mode & %v) !== 0;
		};

		os.FileInfo.prototype.isRegular = function() {
		  return (this.mode & %v) === 0;
		};

		os.FileInfo.prototype.perm = function() {
		  return this.mode & 0777;
		};
	`), uint(os.ModeDir), uint(os.ModeType)))
	vm.Run(script)

	return vm
}

func throw(vm *otto.Otto, err error) otto.Value {
	if _, ok := err.(*otto.Error); ok {
		panic(err)
	}
	panic(vm.MakeCustomError("Error", err.Error()))
}

func ottoError(err error) error {
	switch e := err.(type) {
	case *otto.Error:
		err = errors.New(strings.TrimSpace(e.String()))
	case parser.ErrorList:
		err = fmt.Errorf("%v:%v:%v: %v", e[0].Position.Filename, e[0].Position.Line, e[0].Position.Column, e[0].Message)
	}
	return err
}

func os_getwd(call otto.FunctionCall) otto.Value {
	wd, _ := os.Getwd()
	v, _ := call.Otto.ToValue(wd)
	return v
}

func os_mkdir(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	path, _ := call.ArgumentList[0].ToString()
	var perm os.FileMode
	if v := call.Argument(1); v.IsNumber() {
		i, _ := v.ToInteger()
		perm = os.FileMode(i)
	}
	if perm == 0 {
		perm = os.FileMode(0777)
	}
	if os.MkdirAll(path, perm) != nil {
		return otto.TrueValue()
	}
	return otto.UndefinedValue()
}

func os_open(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	name, _ := call.ArgumentList[0].ToString()
	flag := os.O_RDONLY
	switch mode, _ := call.Argument(1).ToString(); mode {
	case "r":
		flag = os.O_RDONLY
	case "r+":
		flag = os.O_RDWR
	case "w":
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "w+":
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	case "a":
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "a+":
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}

	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return throw(call.Otto, err)
	}
	v, _ := call.Otto.Call(`new os.File`, nil, newFile(call.Otto, f))
	return v
}

func os_remove(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	path, _ := call.ArgumentList[0].ToString()
	os.RemoveAll(path)
	return otto.UndefinedValue()
}

func os_rename(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 2 {
		return otto.UndefinedValue()
	}

	src, _ := call.ArgumentList[0].ToString()
	dst, _ := call.ArgumentList[1].ToString()
	if os.Rename(src, dst) != nil {
		return otto.TrueValue()
	}
	return otto.UndefinedValue()
}

func os_stat(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	path, _ := call.ArgumentList[0].ToString()
	fi, err := os.Stat(path)
	if err != nil {
		return otto.UndefinedValue()
	}
	mtime, _ := call.Otto.Call(`new Date`, nil, fi.ModTime().Unix()*1000+int64(fi.ModTime().Nanosecond())/int64(time.Millisecond))
	v, _ := call.Otto.Call(`new os.FileInfo`, nil, fi.Name(), fi.Size(), fi.Mode(), mtime)
	return v
}

func os_system(call otto.FunctionCall) otto.Value {
	// defaults
	var dir string
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
		// dir
		v, _ = options.Get("dir")
		if v.IsString() {
			dir, _ = v.ToString()
		}

		redir := func(o *otto.Object, k string) (w io.WriteCloser, err error) {
			switch v, _ = o.Get(k); {
			case v.IsString():
				s, _ := v.ToString()
				w, err = os.Create(s)
			case v.IsNull():
				w = discard
			case v.Class() == "Array":
				w = newBuffer(call.Otto, v.Object())
			}
			return
		}
		// stdout
		switch w, err := redir(options, "stdout"); {
		case err != nil:
			return throw(call.Otto, err)
		case w != nil:
			stdout = w
			defer w.Close()
		}
		// stderr
		switch w, err := redir(options, "stderr"); {
		case err != nil:
			return throw(call.Otto, err)
		case w != nil:
			stderr = w
			defer w.Close()
		}
	}

	cmd := binfmt.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return otto.TrueValue()
		}
		return throw(call.Otto, err)
	}
	return otto.UndefinedValue()
}

func os_whence(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	s, _ := call.ArgumentList[0].ToString()
	name, err := exec.LookPath(s)
	if err != nil {
		return otto.UndefinedValue()
	}
	v, _ := call.Otto.ToValue(name)
	return v
}

type buffer struct {
	vm  *otto.Otto
	ary *otto.Object

	mu sync.Mutex
	b  bytes.Buffer
}

func newBuffer(vm *otto.Otto, ary *otto.Object) *buffer {
	return &buffer{
		vm:  vm,
		ary: ary,
	}
}

func (b *buffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.b.Write(p)
	for {
		if s, err := b.b.ReadString('\n'); err == nil {
			b.ary.Call("push", trimNewline(s))
		} else {
			b.b.WriteString(s)
			break
		}
	}
	return len(p), nil
}

func (b *buffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if 0 < b.b.Len() {
		b.ary.Call("push", b.b.String())
	}
	return nil
}

type file struct {
	vm *otto.Otto
	f  *os.File

	br *bufio.Reader
}

func newFile(vm *otto.Otto, f *os.File) *file {
	return &file{
		vm: vm,
		f:  f,
		br: bufio.NewReader(f),
	}
}

func (f *file) Close(call otto.FunctionCall) otto.Value {
	err := f.f.Close()

	if err != nil {
		return otto.TrueValue()
	}
	return otto.UndefinedValue()
}

func (f *file) Name(call otto.FunctionCall) otto.Value {
	v, _ := call.Otto.ToValue(f.f.Name())
	return v
}

func (f *file) Read(call otto.FunctionCall) otto.Value {
	v, _ := call.Argument(0).ToInteger()
	n := int(v)
	if n < 0 {
		n = 0
	}
	p := make([]byte, n)
	n, err := f.br.Read(p)

	if err != nil && err != io.EOF {
		return throw(f.vm, err)
	}
	rv, _ := f.vm.Object(`({})`)
	rv.Set("eof", err == io.EOF)
	rv.Set("buffer", string(p[:n]))
	return rv.Value()
}

func (f *file) ReadLine(call otto.FunctionCall) otto.Value {
	s, err := f.br.ReadString('\n')

	if err != nil && err != io.EOF {
		return throw(f.vm, err)
	}
	rv, _ := f.vm.Object(`({})`)
	rv.Set("eof", err == io.EOF)
	rv.Set("buffer", trimNewline(s))
	return rv.Value()
}

func (f *file) Write(call otto.FunctionCall) otto.Value {
	v, _ := call.Argument(0).ToString()
	_, err := f.f.WriteString(v)

	if err != nil {
		return throw(f.vm, err)
	}
	return otto.UndefinedValue()
}
