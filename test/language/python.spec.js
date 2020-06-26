//
// aster :: language/python.spec.js
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

os.whence = jest.fn();

const python = require('../../lib/language/python');

describe('language', () => {
  describe('python', () => {
    describe('.coverage()', () => {
      it('should notify "coverage not found"', () => {
        os.whence.mockReturnValueOnce(false);

        expect(python.coverage('help')).toBe(true);
        expect(os.whence).lastCalledWith('coverage');
        expect(aster.notify).lastCalledWith('failure', 'aster: coverage', 'coverage not found!');
      });

      it('should execute `coverage help`', () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(python.coverage('help')).toBe(false);
        expect(os.whence).lastCalledWith('coverage');
        expect(language.system).lastCalledWith({
          args: ['coverage', 'help'],
          options: undefined,
          title: 'coverage',
          success: 'help passed',
          failure: 'help failed',
        });
      });
    });

    describe('.coverage', () => {
      describe.each([
        ['annotate'],
        ['combine'],
        ['erase'],
        ['html'],
        ['report'],
        ['xml'],
      ])('.%s()', (cmd) => {
        it(`should execute \`coverage ${cmd}\``, () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(python.coverage[cmd]()).toBe(false);
          expect(os.whence).lastCalledWith('coverage');
          expect(language.system).lastCalledWith({
            args: ['coverage', cmd],
            options: undefined,
            title: 'coverage',
            success: `${cmd} passed`,
            failure: `${cmd} failed`,
          });
        });
      });

      describe('.run()', () => {
        it.each([
          ['spam.py'],
          ['-m spam'],
        ])('should execute `coverage run %s`', (a) => {
          const args = a.split(' ');
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(python.coverage.run.apply(null, args)).toBe(false);
          expect(os.whence).lastCalledWith('coverage');
          expect(language.system).lastCalledWith({
            args: ['coverage', 'run'].concat(args),
            options: undefined,
            title: 'coverage',
            success: 'run passed',
            failure: 'run failed',
          });
        });
      });
    });
  });
});
