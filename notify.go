//
// aster :: notify.go
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

package aster

import (
	"strconv"
	"strings"

	"github.com/mattn/go-gntp"
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

func (g *GNTPValue) Get() interface{} { return string(*g) }
func (g *GNTPValue) String() string   { return string(*g) }
func (g *GNTPValue) IsBoolFlag() bool { return true }

func newNotifier(server string) (c *gntp.Client, err error) {
	if server == "" {
		return
	}

	c = gntp.NewClient()
	c.Server = server
	c.AppName = "Aster"
	err = c.Register([]gntp.Notification{
		{
			Event:   "success",
			Enabled: true,
		},
		{
			Event:   "failure",
			Enabled: true,
		},
	})
	return
}

func notify(c *gntp.Client, name, title, text string) error {
	if c == nil {
		return nil
	}

	return c.Notify(&gntp.Message{
		Event: name,
		Title: title,
		Text:  text,
	})
}
