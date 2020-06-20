//
// aster :: asterfile_test.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster_test

import (
	"testing"

	"github.com/hattya/aster/internal/test"
)

func TestNoAsterfile(t *testing.T) {
	err := test.Sandbox(func() {
		if _, err := test.New(); err == nil {
			t.Error("expected error")
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	err := test.Sandbox(func() {
		src := `++;`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		if _, err := test.New(); err == nil {
			t.Error("expected error")
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestWatchArgs(t *testing.T) {
	err := test.Sandbox(func() {
		src := `aster.watch(/.+\.go/, function() { });`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		a, err := test.New()
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if g, e := a.NumWatches(), 2; g != e {
			t.Errorf("len(Aster.watches) = %v, expected %v", g, e)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestWatchInvalidArgs(t *testing.T) {
	err := test.Sandbox(func() {
		// too few args
		src := `aster.watch();`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		a, err := test.New()
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if g, e := a.NumWatches(), 1; g != e {
			t.Errorf("len(Aster.watches) = %v, expected %v", g, e)
		}

		// invalid args
		src = `aster.watch('', 1);`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		a, err = test.New()
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if g, e := a.NumWatches(), 1; g != e {
			t.Errorf("len(Aster.watches) = %v, expected %v", g, e)
		}

		// invalid args
		src = `aster.watch(/.+/, 1);`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		a, err = test.New()
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if g, e := a.NumWatches(), 1; g != e {
			t.Errorf("len(Aster.watches) = %v, expected %v", g, e)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNotifyArgs(t *testing.T) {
	err := test.Sandbox(func() {
		src := `aster.notify('name', 'title', 'text');`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		if _, err := test.New(); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNotifyInvalidArgs(t *testing.T) {
	err := test.Sandbox(func() {
		// too few args
		src := `aster.notify();`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		if _, err := test.New(); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTitleArgs(t *testing.T) {
	err := test.Sandbox(func() {
		src := `aster.title('aster.test');`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		if _, err := test.New(); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTitleInvalidArgs(t *testing.T) {
	err := test.Sandbox(func() {
		// too few args
		src := `aster.title();`
		if err := test.Gen(src); err != nil {
			t.Fatal(err)
		}
		if _, err := test.New(); err != nil {
			t.Fatal("unexpected error:", err)
		}
	})
	if err != nil {
		t.Error(err)
	}
}
