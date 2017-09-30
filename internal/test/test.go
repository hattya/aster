//
// test :: test.go
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

package test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hattya/aster"
	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/go.cli"
)

func Gen(src string) error {
	return ioutil.WriteFile("Asterfile", []byte(src), 0666)
}

func New() (*aster.Aster, error) {
	return aster.New(cli.NewCLI(), "")
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
