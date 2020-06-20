//
// aster :: export_test.go
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

var (
	NewBuffer = newBuffer
	NewVM     = newVM
)

func (a *Aster) NumWatches() int {
	return len(a.watches)
}
