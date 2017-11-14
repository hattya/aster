//
// aster/cmd/aster :: aster_windows.go
//
//   Copyright (c) 2017 Akinori Hattori <hattya@gmail.com>
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

package main

import (
	"github.com/hattya/go.notify"
	"github.com/hattya/go.notify/freedesktop"
	"github.com/hattya/go.notify/windows"
)

var impls = map[string]interface{}{
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
