//
// aster/cmd/aster :: init_windows.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package main

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrEnv = errors.New("%APPDATA% is not set")

func init() {
	configDir = func() (string, error) {
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", ErrEnv
		}
		return filepath.Join(appdata, "Aster"), nil
	}
}
