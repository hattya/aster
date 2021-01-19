//
// aster/internal/test :: test.go
//
//   Copyright (c) 2014-2021 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/hattya/aster"
	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/go.cli"
	"github.com/hattya/go.notify"
)

var ci bool

func init() {
	// Ubuntu 14.04 never reports Write
	ci = os.Getenv("CI") != "" && runtime.GOOS == "linux"
}

func Gen(src string) error {
	if ci {
		os.Remove("Asterfile")
	}
	return ioutil.WriteFile("Asterfile", []byte(src), 0o666)
}

func New() (*aster.Aster, error) {
	return aster.New(cli.NewCLI(), nil)
}

func Sandbox(test interface{}) error {
	dir, err := sh.Mkdtemp()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	popd, err := sh.Pushd(dir)
	if err != nil {
		return err
	}
	defer popd()

	switch t := test.(type) {
	case func():
		t()
		return nil
	case func() error:
		return t()
	default:
		return fmt.Errorf("unknown type: %T", test)
	}
}

type Notifier struct {
	calls map[string][]error
}

func NewNotifier() *Notifier {
	return &Notifier{calls: make(map[string][]error)}
}

func (p *Notifier) MockCall(method string, err error) {
	p.calls[method] = append(p.calls[method], err)
}

func (p *Notifier) Close() error {
	return p.mock("Close")
}

func (p *Notifier) Register(string, notify.Icon, map[string]interface{}) error {
	return p.mock("Register")
}

func (p *Notifier) Notify(string, string, string) error {
	return p.mock("Notify")
}

func (p *Notifier) mock(method string) error {
	if len(p.calls[method]) == 0 {
		return nil
	}
	err := p.calls[method][0]
	p.calls[method] = p.calls[method][1:]
	return err
}

func (p *Notifier) Sys() interface{} {
	return nil
}
