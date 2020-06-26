//
// aster :: language/javascript.spec.js
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

const javascript = require('../../lib/language/javascript');

describe('language', () => {
  describe('javascript', () => {
    describe('.npm()', () => {
      it('should notify "npm not found"', () => {
        os.whence.mockReturnValueOnce(false);

        expect(javascript.npm('prefix')).toBe(true);
        expect(os.whence).lastCalledWith('npm');
        expect(aster.notify).lastCalledWith('failure', 'aster: npm', 'npm not found!');
      });

      it('should execute `npm prefix`', () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(javascript.npm('prefix')).toBe(false);
        expect(os.whence).lastCalledWith('npm');
        expect(language.system).lastCalledWith({
          args: ['npm', 'prefix'],
          options: undefined,
          title: 'npm',
          success: 'prefix passed',
          failure: 'prefix failed',
        });
      });
    });

    describe('.npm', () => {
      describe.each([
        ['install'],
        ['test'],
      ])('.%s()', (cmd) => {
        it(`should execute \`npm ${cmd}\``, () => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(javascript.npm[cmd]()).toBe(false);
          expect(os.whence).lastCalledWith('npm');
          expect(language.system).lastCalledWith({
            args: ['npm', cmd],
            options: undefined,
            title: 'npm',
            success: `${cmd} passed`,
            failure: `${cmd} failed`,
          });
        });
      });

      describe('.run()', () => {
        it.each([
          ['cover'],
          ['lint'],
        ])('should execute `npm run %s`', (script) => {
          os.whence.mockReturnValueOnce(true);
          language.system.mockReturnValueOnce(false);

          expect(javascript.npm.run(script)).toBe(false);
          expect(os.whence).lastCalledWith('npm');
          expect(language.system).lastCalledWith({
            args: ['npm', 'run', script],
            options: undefined,
            title: 'npm',
            success: `${script} script passed`,
            failure: `${script} script failed`,
          });
        });
      });
    });
  });
});
