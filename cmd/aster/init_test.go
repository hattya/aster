//
// aster/cmd/aster :: init_test.go
//
//   Copyright (c) 2015-2021 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
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
		if err := ioutil.WriteFile(filepath.Join("template", "go"), []byte(tmpl+"\n"), 0o666); err != nil {
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
	if _, err := configDir(); err == nil {
		t.Fatal("expected error")
	}
}
