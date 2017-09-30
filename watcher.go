//
// aster :: watcher.go
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

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

func init() {
	app.Flags.Duration("s", 727*time.Millisecond, "squash events during <duration> (default: 727ms)")
	app.Flags.MetaVar("s", " <duration>")
}

type Watcher struct {
	*fsnotify.Watcher

	a    *Aster
	quit chan struct{}
}

func newWatcher(a *Aster) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		Watcher: fsw,
		a:       a,
		quit:    make(chan struct{}),
	}
	if err := w.Update("."); err != nil {
		w.Close()
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Close() error {
	w.quit <- struct{}{}
	<-w.quit
	return w.Watcher.Close()
}

func (w *Watcher) Update(root string) error {
	fi, err := os.Lstat(root)
	if err != nil || !fi.IsDir() {
		return err
	}
	return w.update(root, fi, false)
}

func (w *Watcher) update(path string, fi os.FileInfo, ignore bool) (err error) {
	if !ignore {
		ignore = w.a.Ignore(path)
	}
	if ignore {
		w.Remove(path)
	} else if err = w.Add(path); err != nil {
		return
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, fi := range list {
		if fi.IsDir() {
			if err = w.update(filepath.Join(path, fi.Name()), fi, ignore); err != nil {
				return
			}
		}
	}
	return
}

func (w *Watcher) Watch() {
	var mu sync.Mutex
	files := make(map[string]int)
	fire := make(chan struct{}, 1)
	done := make(chan struct{}, 1)
	var retry int32

	timer := time.AfterFunc(0, func() {
		mu.Lock()
		defer mu.Unlock()

		if 0 < len(files) {
			select {
			case fire <- struct{}{}:
			default:
			}
		}
	})

	done <- struct{}{}
	for {
		select {
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
				case fi.IsDir():
					if err := w.Update(ev.Name); err != nil {
						warn(err)
					}
					continue
				}
			}
			if ev.Op == fsnotify.Chmod {
				continue
			}

			mu.Lock()
			n := len(files)
			if ev.Op&fsnotify.Remove != 0 || ev.Op&fsnotify.Rename != 0 {
				w.Remove(ev.Name)
				delete(files, ev.Name)
			} else {
				files[ev.Name]++
			}
			mu.Unlock()
			// new cycle has begun
			if n == 0 {
				timer.Reset(app.Flags.Get("s").(time.Duration))
			}
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
				w.a.OnChange(ss)
				if w.a.Reloaded() {
					if err := w.Update("."); err != nil {
						warn(err)
					}
				}
				// retry
				if 0 < atomic.SwapInt32(&retry, 0) {
					select {
					case fire <- struct{}{}:
					default:
					}
				}
				done <- struct{}{}
			}()
		case err := <-w.Errors:
			warn(err)
		case <-w.quit:
			timer.Stop()
			atomic.SwapInt32(&retry, 0)
			<-done
			close(w.quit)
			return
		}
	}
}
