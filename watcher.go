//
// aster :: watcher.go
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
	"sync/atomic"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

func init() {
	app.Flags.Duration("s", 727*time.Millisecond, "squash events during <duration> (default: 727ms)")
	app.Flags.MetaVar("s", " <duration>")
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
		if af.ignore.Match(path) {
			return filepath.SkipDir
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
	var retry int32

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
			// remove "./" prefix
			if 2 < len(ev.Name) && ev.Name[0] == '.' && os.IsPathSeparator(ev.Name[1]) {
				ev.Name = ev.Name[2:]
			}
			// filter
			if ev.Op&fsnotify.Create != 0 {
				switch fi, err := os.Lstat(ev.Name); {
				case err != nil:
					// removed immediately?
					continue
				case fi.IsDir() && !w.af.ignore.Match(ev.Name):
					w.Add(ev.Name)
					continue
				}
			}
			if ev.Op&fsnotify.Remove != 0 || ev.Op&fsnotify.Rename != 0 {
				w.Remove(ev.Name)
			}
			if ev.Op == fsnotify.Chmod {
				continue
			}

			mu.Lock()
			n := len(files)
			if ev.Op&fsnotify.Remove != 0 || ev.Op&fsnotify.Rename != 0 {
				delete(files, ev.Name)
			} else {
				files[ev.Name]++
			}
			mu.Unlock()
			// new cycle has begun
			if n == 0 {
				timer.Reset(app.Flags.Get("s").(time.Duration))
			}
		case err := <-w.Errors:
			warn(err)
		case <-fire:
			go func() {
				select {
				case <-done:
				default:
					// retry later
					atomic.AddInt32(&retry, 1)
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
				// retry
				if 0 < atomic.SwapInt32(&retry, 0) {
					select {
					case fire <- true:
					default:
					}
				}
				done <- true
			}()
		}
	}
}
