//
// aster :: aster_test.go
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
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestWatch(t *testing.T) {
	err := sandbox(func() {
		js := `aster.watch(/.+\.go$/, function(files) {
  //console.log(files);
})`
		if err := genAsterfile(js); err != nil {
			t.Fatal(err)
		}
		af, err := newAsterfile()
		if err != nil {
			t.Fatal(err)
		}

		watcher, err := newWatcher(af)
		if err != nil {
			exit(err)
		}
		defer watcher.Close()

		go watcher.Watch()
		touch("a.go")
		touch("b.go")
		os.Rename("b.go", "_b.go")
		os.Remove("_b.go")
		touch("c.go")
		time.Sleep(1201 * time.Millisecond)
	})
	if err != nil {
		t.Error(err)
	}
}

func genAsterfile(js string) error {
	return ioutil.WriteFile("Asterfile", []byte(js), 0666)
}

func sandbox(test func()) error {
	dir, err := ioutil.TempDir("", "aster.test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	popd, err := pushd(dir)
	if err != nil {
		return err
	}
	defer popd()

	test()
	return nil
}

func pushd(path string) (func() error, error) {
	wd, err := os.Getwd()
	popd := func() error {
		if err == nil {
			return os.Chdir(wd)
		}
		return err
	}
	return popd, os.Chdir(path)
}

func touch(path string) error {
	return ioutil.WriteFile(path, []byte{}, 0666)
}
