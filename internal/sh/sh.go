//
// aster/internal/sh :: sh.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
