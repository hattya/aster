//
// aster :: language/go.spec.js
//
//   Copyright (c) 2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

global.aster = {
  notify: jest.fn(),
};

jest.mock('os');
jest.mock('language', () => ({ prefix: 'aster: ', system: jest.fn() }), { virtual: true });

const os = require('os');
const language = require('language');

os.getwd = jest.fn();
os.open = jest.fn();
os.system = jest.fn();
os.whence = jest.fn();

const path = require('path');
const process = require('process');
const go = require('../../lib/language/go');

describe('language', () => {
  describe('go', () => {
    describe('.dep()', () => {
      it('should notify "dep not found"', () => {
        os.whence.mockReturnValueOnce(false);

        expect(go.dep('version')).toBe(true);
        expect(os.whence).lastCalledWith('dep');
        expect(aster.notify).lastCalledWith('failure', 'aster: dep', 'dep not found!');
      });

      it('should execute `dep version`', () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(go.dep('version')).toBe(false);
        expect(language.system).lastCalledWith({
          args: ['dep', 'version'],
          options: undefined,
          title: 'dep',
          success: 'version passed',
          failure: 'version failed',
        });
      });
    });

    describe('.dep', () => {
      describe.each([
        ['ensure'],
        ['prune'],
      ])('.%s()', (cmd) => {
        it(`should execute \`dep ${cmd}\``, () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(go.dep[cmd]()).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['dep', cmd],
            options: undefined,
            title: 'dep',
            success: `${cmd} passed`,
            failure: `${cmd} failed`,
          });
        });
      });
    });

    describe('.go()', () => {
      it.each([
        ['version', 0],
        ['tool pprof cpu.prof', 1],
      ])('should execute `go %s`', (a, i) => {
        const args = a.split(' ');
        language.system.mockReturnValueOnce(false);

        expect(go.go.apply(null, args)).toBe(false);
        expect(language.system).lastCalledWith({
          args: ['go'].concat(args),
          options: undefined,
          title: 'go',
          success: `${args[i]} passed`,
          failure: `${args[i]} failed`,
        });
      });
    });

    describe('.go', () => {
      describe.each([
        ['build'],
        ['fix'],
        ['fmt'],
        ['generate'],
        ['get'],
        ['install'],
        ['run'],
        ['vet'],
      ])('.%s()', (cmd) => {
        it(`should execute \`go ${cmd} ./...\``, () => {
          language.system.mockReturnValueOnce(false);

          expect(go.go[cmd]('./...')).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['go', cmd, './...'],
            options: undefined,
            title: 'go',
            success: `${cmd} passed`,
            failure: `${cmd} failed`,
          });
        });
      });

      describe('.env()', () => {
        it('should execute `go env`, and return an array', () => {
          expect(go.go.env()).toStrictEqual([]);
          expect(language.system).lastCalledWith({
            args: ['go', 'env'],
            options: { stdout: [] },
            title: 'go',
            success: 'env passed',
            failure: 'env failed',
          });
        });
      });

      describe('.list()', () => {
        it('should execute `go list -f {{.Dir}}`, and return an array', () => {
          [
            [],
            ['-f', '{{.ImportPath}}'],
            ['-f={{.ImportPath}}'],
            ['-json'],
            ['-json=true'],
          ].forEach((args) => {
            expect(go.go.list.apply(null, args)).toStrictEqual([]);
            expect(language.system).lastCalledWith({
              args: ['go', 'list', '-f', '{{.Dir}}'],
              options: { stdout: [] },
              title: 'go',
              success: 'list passed',
              failure: 'list failed',
            });
          });
        });

        it('should execute `go list -f {{.Dir}} ./...`, and return an array', () => {
          expect(go.go.list('./...')).toStrictEqual([]);
          expect(language.system).lastCalledWith({
            args: ['go', 'list', '-f', '{{.Dir}}', './...'],
            options: { stdout: [] },
            title: 'go',
            success: 'list passed',
            failure: 'list failed',
          });
        });
      });

      describe('.mod', () => {
        describe.each([
          ['download'],
          ['tidy'],
          ['vendor'],
        ])('.%s()', (cmd) => {
          it(`should execute \`go mod ${cmd}\``, () => {
            language.system.mockReturnValueOnce(false);

            expect(go.go.mod[cmd]()).toBe(false);
            expect(language.system).lastCalledWith({
              args: ['go', 'mod', cmd],
              options: undefined,
              title: 'go',
              success: `mod ${cmd} passed`,
              failure: `mod ${cmd} failed`,
            });
          });
        });
      });

      describe('.test()', () => {
        it('should execute `go test -race ./...`', () => {
          aster.arch = 'amd64';
          language.system.mockReturnValueOnce(false);

          expect(go.go.test('-race', './...')).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['go', 'test', '-race', './...'],
            options: undefined,
            title: 'go',
            success: 'test passed',
            failure: 'test failed',
          });

          delete aster.arch;
        });

        it('should execute `go test ./...`', () => {
          aster.arch = '386';
          language.system.mockReturnValueOnce(false);

          expect(go.go.test('-race', './...')).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['go', 'test', './...'],
            options: undefined,
            title: 'go',
            success: 'test passed',
            failure: 'test failed',
          });

          delete aster.arch;
        });
      });

      describe('.tool', () => {
        describe('.cover()', () => {
          it.each([
            ['-func cover.out'],
            ['-html cover.out -o coverage.html'],
          ])('should execute `go tool cover %s`', (a) => {
            const args = a.split(' ');
            language.system.mockReturnValueOnce(false);

            expect(go.go.tool.cover.apply(null, args)).toBe(false);
            expect(language.system).lastCalledWith({
              args: ['go', 'tool', 'cover'].concat(args),
              options: undefined,
              title: 'go',
              success: `cover ${args[0]} passed`,
              failure: `cover ${args[0]} failed`,
            });
          });
        });
      });
    });

    describe('.combine()', () => {
      it('should combine the specified coverage profiles', () => {
        const b = [];
        os.open.mockImplementation((name) => {
          let i = 0;
          return {
            close: () => {},
            readLine: () => ({ buffer: name, eof: i++ > 1 }),
            write: (s) => b.push(s),
          };
        });
        const spy = jest.spyOn(go.go, 'list').mockImplementation(() => ([
          process.cwd(),
          path.join(process.cwd(), 'cmd', 'aster'),
        ]));

        expect(go.combine({ out: 'cover.all.out', profile: 'cover.out', packages: ['./...'] })).toBe('cover.all.out');
        expect(b).toStrictEqual([
          'mode: atomic\n',
          `${path.join(process.cwd(), 'cover.out')}\n`,
          `${path.join(process.cwd(), 'cmd', 'aster', 'cover.out')}\n`,
        ]);

        os.open.mockClear();
        spy.mockRestore();
      });
    });

    describe('.packagesOf()', () => {
      it('should return an empty array', () => {
        os.system.mockReturnValueOnce(true);

        expect(go.packagesOf(['aster.go'])).toHaveLength(0);
        expect(os.system).toBeCalled();
      });

      it('should return an array of packages', () => {
        os.getwd.mockReturnValueOnce(process.cwd());
        os.system.mockImplementationOnce((_, obj) => {
          obj.stdout.push([
            process.cwd(),
            'github.com/hattya/aster',
            [
              'github.com/hattya/aster',
              'github.com/hattya/aster/internal/sh',
              'github.com/hattya/aster/internal/test',
            ].join(','),
          ].join('\t'));
          obj.stdout.push([
            path.join(process.cwd(), 'cmd', 'aster'),
            'github.com/hattya/aster/cmd/aster',
            [
              'github.com/hattya/aster',
            ].join(','),
          ].join('\t'));
          obj.stdout.push([
            path.join(process.cwd(), 'internal', 'sh'),
            'github.com/hattya/aster/internal/sh',
            [
            ].join(','),
          ].join('\t'));
          obj.stdout.push([
            path.join(process.cwd(), 'internal', 'test'),
            'github.com/hattya/aster/internal/test',
            [
              'github.com/hattya/aster',
              'github.com/hattya/aster/internal/sh',
            ].join(','),
          ].join('\t'));
          return false;
        });

        expect(go.packagesOf([
          path.join('internal', 'test', 'test.go'),
          path.join('internal', 'sh', 'sh.go'),
          path.join('cmd', 'aster', 'aster.go'),
          path.join('_', '_.go'),
          'aster.go',
        ])).toHaveLength(4);
        expect(os.getwd).toBeCalled();
        expect(os.system).toBeCalled();
      });
    });
  });
});
