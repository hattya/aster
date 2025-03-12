//
// aster :: notify.go
//
//   Copyright (c) 2014-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import (
	"strconv"
	"strings"
)

type GNTPValue string

func (g *GNTPValue) Set(s string) error {
	if v, err := strconv.ParseBool(s); err == nil || s == "" {
		if v {
			*g = "localhost:23053"
		} else {
			*g = ""
		}
	} else {
		if !strings.Contains(s, ":") {
			s += ":23053"
		}
		*g = GNTPValue(s)
	}
	return nil
}

func (g *GNTPValue) Get() any         { return string(*g) }
func (g *GNTPValue) String() string   { return string(*g) }
func (g *GNTPValue) IsBoolFlag() bool { return true }
