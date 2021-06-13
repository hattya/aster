//
// aster :: language/vimscript.spec.js
//
//   Copyright (c) 2020-2021 Akinori Hattori <hattya@gmail.com>
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

os.whence = jest.fn();

const path = require('path');
const vimscript = require('../../lib/language/vimscript');

describe('language', () => {
  describe('vimscript', () => {
    describe('.covimerage()', () => {
      it('should notify "covimerage not found"', () => {
        os.whence.mockReturnValueOnce(false);

        expect(vimscript.covimerage('--version')).toBe(true);
        expect(os.whence).lastCalledWith('covimerage');
        expect(aster.notify).lastCalledWith('failure', 'aster: covimerage', 'covimerage not found!');
      });

      it('should execute `covimerage --version`', () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(vimscript.covimerage('--version')).toBe(false);
        expect(language.system).lastCalledWith({
          args: ['covimerage', '--version'],
          options: undefined,
          title: 'covimerage',
          success: '--version passed',
          failure: '--version failed',
        });
      });
    });

    describe('.covimerage', () => {
      describe.each([
        ['report'],
        ['xml'],
      ])('.%s()', (cmd) => {
        it(`should execute \`covimerage ${cmd}\``, () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(vimscript.covimerage[cmd]()).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['covimerage', cmd],
            options: undefined,
            title: 'covimerage',
            success: `${cmd} passed`,
            failure: `${cmd} failed`,
          });
        });
      });

      describe('.run()', () => {
        it("should execute `covimerage run vim -Nu test/vimrc -c 'Vader! test/**'`", () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(vimscript.covimerage.run('vim', '-Nu', 'test/vimrc', '-c', "'Vader! test/**'")).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['covimerage', 'run', 'vim', '-Nu', 'test/vimrc', '-c', "'Vader! test/**'"],
            options: undefined,
            title: 'covimerage',
            success: 'run passed',
            failure: 'run failed',
          });
        });
      });

      describe('.write_coverage()', () => {
        it('should execute `covimerage write_coverage profile.txt`', () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(vimscript.covimerage.write_coverage('profile.txt')).toBe(false);
          expect(language.system).lastCalledWith({
            args: ['covimerage', 'write_coverage', 'profile.txt'],
            options: undefined,
            title: 'covimerage',
            success: 'write_coverage passed',
            failure: 'write_coverage failed',
          });
        });
      });
    });

    describe('.themis()', () => {
      it('should notify "themis not found"', () => {
        os.whence.mockClear().mockReturnValue(false);

        expect(vimscript.themis()).toBe(true);
        expect(os.whence).nthCalledWith(1, 'themis');
        expect(os.whence).nthCalledWith(2, path.join('.', 'vim-themis', 'bin', 'themis'));
        expect(os.whence).nthCalledWith(3, path.join('..', 'vim-themis', 'bin', 'themis'));
        expect(aster.notify).lastCalledWith('failure', 'aster: themis', 'themis not found!');
      });

      it.each([
        ['themis'],
        [path.join('.', 'vim-themis', 'bin', 'themis')],
        [path.join('..', 'vim-themis', 'bin', 'themis')],
      ])('should execute `%s --reporter dot`', (script) => {
        os.whence.mockClear().mockImplementation((name) => name === script);
        language.system.mockReturnValueOnce(false);

        expect(vimscript.themis('--reporter', 'dot')).toBe(false);
        expect(os.whence).lastCalledWith(script);
        expect(language.system).lastCalledWith({
          args: [script, '--reporter', 'dot'],
          title: 'themis',
          success: 'passed',
          failure: 'failed',
        });
      });
    });
  });
});
