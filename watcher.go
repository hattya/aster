//
// aster :: watch.go
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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var asterS time.Duration

func init() {
	flag.DurationVar(&asterS, "s", 727*time.Millisecond, "")
}

type Watcher struct {
	*fsnotify.Watcher

	af   *Aster
	quit chan bool
}

func newWatcher(af *Aster) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return err
		}
		return watcher.Add(path)
	})
	if err != nil {
		watcher.Close()
		return nil, err
	}
	w := &Watcher{
		Watcher: watcher,
		af:      af,
		quit:    make(chan bool),
	}
	return w, nil
}

func (w *Watcher) Close() error {
	w.quit <- true
	return w.Watcher.Close()
}

func (w *Watcher) Watch() {
	var mu sync.Mutex
	files := make(map[string]int)
	fire := make(chan bool, 1)
	done := make(chan bool, 1)

	timer := time.AfterFunc(0, func() {
		mu.Lock()
		defer mu.Unlock()
		if 0 < len(files) {
			select {
			case fire <- true:
			default:
			}
		}
	})
	done <- true

	for {
		select {
		case <-w.quit:
			return
		case ev := <-w.Events:
			// filter
			switch ev.Op {
			case fsnotify.Create:
				w.Add(ev.Name)
			case fsnotify.Remove:
				w.Remove(ev.Name)
			case fsnotify.Rename, fsnotify.Chmod:
				continue
			}
			name := ev.Name
			// remove "./" prefix
			if 2 < len(name) && name[0] == '.' && os.IsPathSeparator(name[1]) {
				name = name[2:]
			}

			mu.Lock()
			n := len(files)
			switch ev.Op {
			case fsnotify.Remove:
				delete(files, name)
			default:
				files[name]++
			}
			mu.Unlock()
			// new cycle has begun
			if n == 0 {
				timer.Reset(asterS)
			}
		case err := <-w.Errors:
			warn(err)
		case <-fire:
			go func() {
				select {
				case <-done:
				case fire <- true:
					// retry
					return
				default:
					return
				}

				// create snapshot & clear
				mu.Lock()
				ss := make(map[string]int)
				for n, c := range files {
					ss[n] = c
					delete(files, n)
				}
				mu.Unlock()
				// process
				w.af.OnChange(ss)
				done <- true
			}()
		}
	}
}
