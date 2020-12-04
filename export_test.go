//
// aster :: export_test.go
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import "sort"

var (
	NewBuffer = newBuffer
	NewVM     = newVM
)

func (a *Aster) NumWatches() int {
	return len(a.watches)
}

func (w *Watcher) Paths() []string {
	w.mu.Lock()
	defer w.mu.Unlock()

	paths := make(sort.StringSlice, len(w.paths))
	i := 0
	for k := range w.paths {
		paths[i] = k
		i++
	}
	paths.Sort()
	return paths
}
