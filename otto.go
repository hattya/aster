//
// aster :: otto.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/hattya/go.binfmt"
	"github.com/hattya/otto.module"
	"github.com/robertkrimen/otto"
)

func newVM() *module.Otto {
	vm, err := module.New()
	if err != nil {
		panic(err)
	}
	vm.Register(new(stdLoader))

	file := new(module.FileLoader)
	folder := &module.FolderLoader{File: file}
	vm.Register(file)
	vm.Register(folder)
	vm.Register(&module.NodeModulesLoader{
		File:   file,
		Folder: folder,
	})
	// os binding
	vm.Bind("os", func(o *otto.Object) error {
		o.Set("MODE_DIR", os.ModeDir)
		o.Set("MODE_TYPE", os.ModeType)
		o.Set("MODE_PERM", os.ModePerm)

		m := new(os_)
		o.Set("getwd", m.getwd)
		o.Set("mkdir", m.mkdir)
		o.Set("open", m.open)
		o.Set("remove", m.remove)
		o.Set("rename", m.rename)
		o.Set("stat", m.stat)
		o.Set("system", m.system)
		o.Set("whence", m.whence)
		return nil
	})
	// for backward compatibility
	vm.Run(`var os = require('os');`)
	return vm
}

type os_ struct {
}

func (*os_) getwd(call otto.FunctionCall) otto.Value {
	wd, _ := os.Getwd()
	v, _ := call.Otto.ToValue(wd)
	return v
}

func (*os_) mkdir(call otto.FunctionCall) otto.Value {
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

func (*os_) open(call otto.FunctionCall) otto.Value {
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
		return module.Throw(call.Otto, err)
	}
	this := call.This.Object()
	this.Set("_impl", newFile(call.Otto, f))
	return call.This
}

func (*os_) remove(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	path, _ := call.ArgumentList[0].ToString()
	os.RemoveAll(path)
	return otto.UndefinedValue()
}

func (*os_) rename(call otto.FunctionCall) otto.Value {
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

func (*os_) stat(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) < 1 {
		return otto.UndefinedValue()
	}

	path, _ := call.ArgumentList[0].ToString()
	fi, err := os.Stat(path)
	if err != nil {
		return otto.UndefinedValue()
	}
	this := call.This.Object()
	this.Set("name", fi.Name())
	this.Set("size", fi.Size())
	this.Set("mode", fi.Mode())
	mtime, _ := call.Otto.Call(`new Date`, nil, fi.ModTime().Unix()*1000+int64(fi.ModTime().Nanosecond())/int64(time.Millisecond))
	this.Set("mtime", mtime)
	return call.This
}

func (*os_) system(call otto.FunctionCall) otto.Value {
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
			return module.Throw(call.Otto, err)
		case w != nil:
			stdout = w
			defer w.Close()
		}
		// stderr
		switch w, err := redir(options, "stderr"); {
		case err != nil:
			return module.Throw(call.Otto, err)
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
		return module.Throw(call.Otto, err)
	}
	return otto.UndefinedValue()
}

func (*os_) whence(call otto.FunctionCall) otto.Value {
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
		return module.Throw(f.vm, err)
	}
	rv, _ := f.vm.Object(`({})`)
	rv.Set("eof", err == io.EOF)
	rv.Set("buffer", string(p[:n]))
	return rv.Value()
}

func (f *file) ReadLine(call otto.FunctionCall) otto.Value {
	s, err := f.br.ReadString('\n')

	if err != nil && err != io.EOF {
		return module.Throw(f.vm, err)
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
		return module.Throw(f.vm, err)
	}
	return otto.UndefinedValue()
}
