//
// aster/internal/test :: test.go
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

package test

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/hattya/aster"
	"github.com/hattya/aster/internal/sh"
	"github.com/hattya/go.cli"
)

var ci bool

func init() {
	// Ubuntu 14.04 never reports Write
	ci = os.Getenv("CI") != "" && runtime.GOOS == "linux"
}

func Gen(src string) error {
	if ci {
		os.Remove("Asterfile")
	}
	return ioutil.WriteFile("Asterfile", []byte(src), 0666)
}

func New() (*aster.Aster, error) {
	return aster.New(cli.NewCLI(), "")
}

func Sandbox(test interface{}) error {
	dir, err := sh.Mkdtemp()
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	popd, err := sh.Pushd(dir)
	if err != nil {
		return err
	}
	defer popd()

	switch t := test.(type) {
	case func():
		t()
		return nil
	case func() error:
		return t()
	default:
		return fmt.Errorf("unknown type: %T", test)
	}
}

var GNTPError = errors.New("INTERNAL_SERVER_ERROR")

type GNTPServer struct {
	Server string

	l    net.Listener
	done chan struct{}
	wg   sync.WaitGroup

	mu     sync.Mutex
	reject map[string]struct{}
}

func NewGNTPServer() *GNTPServer {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(fmt.Sprintf("test: cannot listen: %v", err))
	}
	return &GNTPServer{
		l:      l,
		done:   make(chan struct{}),
		reject: make(map[string]struct{}),
	}
}

func (s *GNTPServer) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reject = make(map[string]struct{})
}

func (s *GNTPServer) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-s.done:
		return
	default:
	}

	close(s.done)
	s.l.Close()
	s.wg.Wait()
}

func (s *GNTPServer) Reject(msgtype string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reject[strings.ToUpper(msgtype)] = struct{}{}
}

func (s *GNTPServer) Start() {
	s.Server = s.l.Addr().String()

	s.wg.Add(1)
	go s.serve()
}

func (s *GNTPServer) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.l.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
			}
			panic(fmt.Sprintf("test: cannot accept: %v", err))
		}

		r := bufio.NewReader(conn)
		l, err := r.ReadString('\n')
		if err != nil {
			panic(fmt.Sprintf("test: cannot read: %v", err))
		}
		v := strings.Split(l, " ")
		if _, ok := s.reject[v[1]]; ok {
			conn.Write([]byte("GNTP/1.0 -ERROR\r\n"))
			conn.Write([]byte("Error-Code: 500\r\n"))
			fmt.Fprintf(conn, "Error-Description: %v\r\n", GNTPError)
		}
		conn.Close()

		select {
		case <-s.done:
			return
		default:
		}
	}
}
