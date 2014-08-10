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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

type Watcher struct {
	af      *Aster
	watcher *fsnotify.Watcher
	qc      chan string
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
		qc:      make(chan string, 1),
		af:      af,
		watcher: watcher,
	}
	return w, nil
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func (w *Watcher) Watch() {
	go w.WaitEvent()
	w.Loop()
}

func (w *Watcher) WaitEvent() {
	for {
		select {
		case ev := <-w.watcher.Events:
			switch ev.Op {
			case fsnotify.Create:
				w.watcher.Add(ev.Name)
			case fsnotify.Remove:
				w.watcher.Remove(ev.Name)
				continue
			case fsnotify.Chmod, fsnotify.Rename:
				continue
			}
			w.qc <- ev.Name
		case err := <-w.watcher.Errors:
			if err != nil {
				warn(err)
			}
		}
	}
}

func (w *Watcher) Loop() {
	var mu sync.Mutex
	files := make(map[string]int)
	fire := make(chan bool, 1)
	done := make(chan bool, 1)

	done <- true
	for {
		select {
		case name := <-w.qc:
			mu.Lock()
			files[name]++
			mu.Unlock()
		squash:
			for {
				select {
				case <-time.After(1 * time.Second):
					break squash
				case name := <-w.qc:
					mu.Lock()
					files[name]++
					mu.Unlock()
				}
			}
			select {
			case fire <- true:
			default:
			}
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
