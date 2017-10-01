//
// aster/internal/sh :: sh.go
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

package sh

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Mkdir(s ...string) error {
	return os.MkdirAll(filepath.Join(s...), 0777)
}

func Mkdtemp() (string, error) {
	dir, err := ioutil.TempDir("", "aster.test")
	return filepath.ToSlash(dir), err
}

func Pushd(s ...string) (func() error, error) {
	wd, err := os.Getwd()
	popd := func() error {
		if err == nil {
			return os.Chdir(wd)
		}
		return err
	}
	return popd, os.Chdir(filepath.Join(s...))
}

func Touch(s ...string) error {
	return ioutil.WriteFile(filepath.Join(s...), []byte{}, 0666)
}
