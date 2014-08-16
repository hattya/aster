//
// aster :: notify.go
//
//   Copyright (c) 2014 Akinori Hattori <hattya@gmail.com>
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
	"flag"

	"github.com/mattn/go-gntp"
)

type growlValue string

func (g *growlValue) Set(value string) error {
	if value == "true" {
		*g = "localhost:23053"
	} else {
		*g = growlValue(value)
	}
	return nil
}

func (g *growlValue) String() string {
	if *g == "" {
		return "false"
	}
	return string(*g)
}

func (g *growlValue) IsBoolFlag() bool {
	return true
}

var asterG growlValue

func init() {
	flag.Var(&asterG, "g", "")
}

func newNotifier() *gntp.Client {
	if asterG == "" {
		return nil
	}
	c := gntp.NewClient()
	c.Server = asterG.String()
	c.AppName = "Aster"
	err := c.Register([]gntp.Notification{
		{
			Event:   "success",
			Enabled: true,
		},
		{
			Event:   "failure",
			Enabled: true,
		},
	})
	if err != nil {
		warn(err)
		return nil
	}
	return c
}

func notify(c *gntp.Client, name, title, text string) {
	if c == nil {
		return
	}
	err := c.Notify(&gntp.Message{
		Event: name,
		Title: title,
		Text:  text,
	})
	if err != nil {
		warn(err)
	}
}
