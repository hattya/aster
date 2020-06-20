//
// aster/cmd/aster :: aster_unix.go
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

// +build !plan9,!windows

package main

import (
	"github.com/hattya/go.notify"
	"github.com/hattya/go.notify/freedesktop"
)

var impls = map[string]interface{}{
	"fdo":         "freedesktop",
	"freedesktop": "freedesktop",
	"gntp":        "gntp",
}

func newNotifier(name, impl string) (n notify.Notifier) {
	switch impl {
	case "freedesktop":
		n, _ = freedesktop.NewNotifier(name)
		if n != nil {
			// discard signals
			go func() {
				c := n.Sys().(*freedesktop.Client)
				for {
					select {
					case <-c.ActionInvoked:
					case <-c.NotificationClosed:
					}
				}
			}()
		}
	}
	return
}
