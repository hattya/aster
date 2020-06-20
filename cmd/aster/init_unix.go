//
// aster/cmd/aster :: init_unix.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package main

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrEnv = errors.New("$HOME is not set")

func init() {
	configDir = func() (string, error) {
		config := os.Getenv("XDG_CONFIG_HOME")
		if config == "" {
			home := os.Getenv("HOME")
			if home == "" {
				return "", ErrEnv
			}
			config = filepath.Join(home, ".config")
		}
		return filepath.Join(config, "aster"), nil
	}
}
