//
// aster :: watcher.go
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

package aster

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	Squash time.Duration

	ctx  context.Context
	a    *Aster
	w    *fsnotify.Watcher
	quit chan struct{}

	mu    sync.Mutex
	paths map[string]struct{}
	done  chan struct{}
}

func NewWatcher(ctx context.Context, a *Aster) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w := &Watcher{
		ctx:   ctx,
		a:     a,
		w:     fsw,
		quit:  make(chan struct{}, 1),
		paths: make(map[string]struct{}),
		done:  make(chan struct{}),
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
	return w.w.Close()
}

func (w *Watcher) Add(name string) error {
	return w.walk(name, func(path string) error {
		if w.a.Ignore(path) {
			return filepath.SkipDir
		}
		return w.add(path)
	})
}

func (w *Watcher) Remove(name string) error {
	return w.walk(name, func(path string) error {
		if err := w.remove(path); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func (w *Watcher) Update(name string) error {
	return w.walk(name, func(path string) error {
		if w.a.Ignore(path) {
			if err := w.remove(path); err != nil {
				return err
			}
			return filepath.SkipDir
		}
		return w.add(path)
	})
}

func (w *Watcher) add(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.paths[name] = struct{}{}
	return w.w.Add(name)
}

func (w *Watcher) remove(name string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.paths[name]; ok {
		delete(w.paths, name)

		name += string(os.PathSeparator)
		for k := range w.paths {
			if strings.HasPrefix(k, name) {
				delete(w.paths, k)
				if err := w.w.Remove(k); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *Watcher) walk(root string, fn func(string) error) error {
	return filepath.Walk(filepath.Clean(root), func(path string, fi os.FileInfo, err error) error {
		select {
		case <-w.ctx.Done():
			return w.ctx.Err()
		default:
		}

		if err != nil || !fi.IsDir() {
			return err
		}
		return fn(path)
	})
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
		case ev := <-w.w.Events:
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
		case err := <-w.w.Errors:
			if err != nil {
				warn(w.a.ui, err)
			}
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
