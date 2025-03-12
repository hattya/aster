//
// aster/cmd/aster :: aster_windows.go
//
//   Copyright (c) 2017-2025 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package main

import (
	"github.com/hattya/go.notify"
	"github.com/hattya/go.notify/freedesktop"
	"github.com/hattya/go.notify/windows"
)

var impls = map[string]any{
	"fdo":         "freedesktop",
	"freedesktop": "freedesktop",
	"gntp":        "gntp",
	"windows":     "windows",
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
	case "windows":
		n, _ = windows.NewNotifier(name, nil)
		if n != nil {
			icon["success"] = windows.IconInfo
			icon["failure"] = windows.IconError
			// discard events
			go func() {
				ni := n.Sys().(*windows.NotifyIcon)
				for {
					select {
					case <-ni.Balloon:
					case <-ni.Menu:
					}
				}
			}()
		}
	}
	return
}
