//
// aster/cmd/aster :: init_test.go
//
//   Copyright (c) 2015-2017 Akinori Hattori <hattya@gmail.com>
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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/aster/internal/test"
)

var initTests = []struct {
	src, newline string
}{
	{"  ", "\n"},
	{"\n", "\n"},
	{"\r", "\r"},
	{"\r\n", "\r\n"},
}

func TestInit(t *testing.T) {
	save := configDir
	defer func() { configDir = save }()

	configDir = os.Getwd
	tmpl := "template go"
	err := test.Sandbox(func() error {
		if err := sh.Mkdir("template"); err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join("template", "go"), []byte(tmpl+"\n"), 0666); err != nil {
			return err
		}
		for _, tt := range initTests {
			if err := test.Gen(tt.src); err != nil {
				return err
			}
			if err := app.Run([]string{"init", "go"}); err != nil {
				return err
			}
			data, err := ioutil.ReadFile("Asterfile")
			if err != nil {
				return err
			}
			if g, e := string(data), tt.src+tt.newline+tmpl+tt.newline; g != e {
				t.Fatalf("expected %q, got %q", e, g)
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestInitError(t *testing.T) {
	save := configDir
	defer func() { configDir = save }()

	app.Stderr = ioutil.Discard
	defer func() { app.Stderr = nil }()

	args := []string{"init", "go"}
	err := test.Sandbox(func() error {
		// error
		configDir = func() (string, error) { return "", fmt.Errorf("error") }
		if err := app.Run(args); err == nil {
			t.Fatal("expected error")
		}

		// template directory not found
		configDir = os.Getwd
		if err := app.Run(args); err != nil {
			return err
		}

		// cannot open Asterfile
		if err := sh.Mkdir("template"); err != nil {
			return err
		}
		if err := sh.Mkdir("Asterfile"); err != nil {
			return err
		}
		if err := app.Run(args); err == nil {
			t.Fatal("expected error")
		}

		// template not found
		if err := os.Remove("Asterfile"); err != nil {
			return err
		}
		if err := app.Run(args); err == nil {
			t.Fatal("expected error")
		}

		// cannot read template
		if err := sh.Mkdir("template", "go"); err != nil {
			return err
		}
		if err := app.Run(args); err == nil {
			t.Fatal("expected error")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestConfigDir(t *testing.T) {
	environ := map[string]string{
		"HOME":    os.Getenv("HOME"),
		"APPDATA": os.Getenv("APPDATA"),
	}
	defer func() {
		for k, v := range environ {
			os.Setenv(k, v)
		}
	}()

	if _, err := configDir(); err != nil {
		t.Fatal(err)
	}

	// error
	for k := range environ {
		os.Setenv(k, "")
	}
	if _, err := configDir(); err != ErrEnv {
		t.Fatalf("expected ErrEnv, got %#v", err)
	}
}
