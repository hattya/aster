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

package aster

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	*fsnotify.Watcher
	Squash time.Duration

	ctx  context.Context
	a    *Aster
	quit chan struct{}

	mu   sync.Mutex
	done chan struct{}
}

func NewWatcher(ctx context.Context, a *Aster) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		Watcher: fsw,
		ctx:     ctx,
		a:       a,
		done:    make(chan struct{}),
		quit:    make(chan struct{}, 1),
	}
	if err := w.Update("."); err != nil {
		fsw.Close()
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	select {
	case <-w.done:
		return nil
	default:
	}

	w.quit <- struct{}{}
	<-w.done
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
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		default:
		}

		if fi.IsDir() {
			if err = w.update(filepath.Join(path, fi.Name()), fi, ignore); err != nil {
				return
			}
		}
	}
	return
}

func (w *Watcher) Watch() error {
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
					go func() {
						if err := w.Update(ev.Name); err != nil {
							warn(w.a.ui, err)
						}
					}()
					continue
				}
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
				timer.Reset(w.Squash)
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
				w.a.OnChange(w.ctx, ss)
				if w.a.Reloaded() {
					if err := w.Update("."); err != nil {
						warn(w.a.ui, err)
					}
				}
				done <- struct{}{}
				// retry
				if 0 < atomic.SwapInt32(&retry, 0) {
					select {
					case fire <- struct{}{}:
					default:
					}
				}
			}()
		case err := <-w.Errors:
			if err != nil {
				warn(w.a.ui, err)
			}
		case <-w.done:
		case <-w.quit:
			<-done
			timer.Stop()
			atomic.SwapInt32(&retry, 0)
			close(w.done)
			return w.ctx.Err()
		case <-w.ctx.Done():
			select {
			case w.quit <- struct{}{}:
			default:
			}
		}
	}
}
