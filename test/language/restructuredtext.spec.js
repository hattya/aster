//
// aster :: language/restructuredtext.spec.js
//
//   Copyright (c) 2020-2025 Akinori Hattori <hattya@gmail.com>
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

const restructuredtext = require('../../lib/language/restructuredtext');

describe('language', () => {
  describe('restructuredtext', () => {
    const options = ['--strict'];
    const src = 'README.rst';
    const dst = 'README.html';

    describe('.rst2html()', () => {
      it('should notify "rst2html not found"', () => {
        os.whence.mockClear().mockReturnValue(false);

        expect(restructuredtext.rst2html()).toBe(true);
        expect(os.whence).toHaveBeenNthCalledWith(1, 'rst2html5.py');
        expect(os.whence).toHaveBeenNthCalledWith(2, 'rst2html5');
        expect(os.whence).toHaveBeenNthCalledWith(3, 'rst2html.py');
        expect(os.whence).toHaveBeenNthCalledWith(4, 'rst2html');
        expect(aster.notify).toHaveBeenLastCalledWith('failure', 'aster: rst2html', 'rst2html not found!');
      });

      it.each([
        ['rst2html5.py'],
        ['rst2html5'],
        ['rst2html.py'],
        ['rst2html'],
      ])(`should execute \`%s ${src} ${dst}\``, (script) => {
        os.whence.mockClear().mockImplementation((name) => name === script);
        language.system.mockReturnValueOnce(false);

        expect(restructuredtext.rst2html({ src })).toBe(false);
        expect(os.whence).toHaveBeenLastCalledWith(script);
        expect(language.system).toHaveBeenLastCalledWith({
          args: [script, src, dst],
          title: 'rst2html',
          success: `${src}`,
          failure: `${src} failed`,
        });
      });

      it.each([
        ['rst2html5.py'],
        ['rst2html5'],
        ['rst2html.py'],
        ['rst2html'],
      ])(`should execute \`%s ${options.join(' ')} ${src} ${dst}\``, (script) => {
        os.whence.mockClear().mockImplementation((name) => name === script);
        language.system.mockReturnValueOnce(false);

        expect(restructuredtext.rst2html({ options, src, dst })).toBe(false);
        expect(os.whence).toHaveBeenLastCalledWith(script);
        expect(language.system).toHaveBeenLastCalledWith({
          args: [script].concat(options, [src, dst]),
          title: 'rst2html',
          success: `${src}`,
          failure: `${src} failed`,
        });
      });
    });
  });
});
