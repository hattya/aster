//
// aster :: otto.go
//
//   Copyright (c) 2014-2015 Akinori Hattori <hattya@gmail.com>
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
	"time"

	"github.com/hattya/go.binfmt"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
)

func newVM() *otto.Otto {
	vm := otto.New()
	// os object
	os, _ := vm.Object(`os = {}`)
	os.Set("getwd", os_getwd)
	os.Set("mkdir", os_mkdir)
	os.Set("remove", os_remove)
	os.Set("rename", os_rename)
	os.Set("stat", os_stat)
	os.Set("system", os_system)
	os.Set("whence", os_whence)

	vm.Run(fmt.Sprintf(`os.FileInfo = function(name, size, mode, mtime) {
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
  return this.mode & %v;
};
`, int64(1<<31), int64(0xff<<24), 0777))

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

func os_getwd(call otto.FunctionCall) otto.Value {
	wd, _ := os.Getwd()
	v, _ := call.Otto.ToValue(wd)
	return v
}

func os_mkdir(call otto.FunctionCall) otto.Value {
	if 1 <= len(call.ArgumentList) {
		path, _ := call.ArgumentList[0].ToString()
		var perm os.FileMode
		v := call.Argument(1)
		if v.IsNumber() {
			i, _ := v.ToInteger()
			perm = os.FileMode(i)
		}
		if perm == 0 {
			perm = os.FileMode(0777)
		}
		if os.MkdirAll(path, perm) != nil {
			return otto.TrueValue()
		}
	}
	return otto.UndefinedValue()
}

func os_remove(call otto.FunctionCall) otto.Value {
	if 1 <= len(call.ArgumentList) {
		path, _ := call.ArgumentList[0].ToString()
		os.RemoveAll(path)
	}
	return otto.UndefinedValue()
}

func os_rename(call otto.FunctionCall) otto.Value {
	if 2 <= len(call.ArgumentList) {
		src, _ := call.ArgumentList[0].ToString()
		dst, _ := call.ArgumentList[1].ToString()
		if os.Rename(src, dst) != nil {
			return otto.TrueValue()
		}
	}
	return otto.UndefinedValue()
}

func os_stat(call otto.FunctionCall) otto.Value {
	if 1 <= len(call.ArgumentList) {
		path, _ := call.ArgumentList[0].ToString()
		if fi, err := os.Stat(path); err == nil {
			mtime, _ := call.Otto.Call(`new Date`, nil, fi.ModTime().Unix()*1000+int64(fi.ModTime().Nanosecond())/int64(time.Millisecond))
			v, _ := call.Otto.Call(`new os.FileInfo`, nil, fi.Name(), fi.Size(), fi.Mode(), mtime)
			return v
		}
	}
	return otto.UndefinedValue()
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
				w = &Buffer{
					vm:  call.Otto,
					ary: v.Object(),
				}
			}
			return
		}
		// stdout
		switch w, err := redir(options, "stdout"); {
		case err != nil:
			return throw(call.Otto.ToValue(err.Error()))
		case w != nil:
			stdout = w
			defer w.Close()
		}
		// stderr
		switch w, err := redir(options, "stderr"); {
		case err != nil:
			return throw(call.Otto.ToValue(err.Error()))
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

	mu sync.Mutex
	b  bytes.Buffer
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
