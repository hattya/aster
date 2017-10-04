//
// aster :: asterfile_test.go
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
