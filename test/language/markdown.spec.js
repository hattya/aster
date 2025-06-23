//
// aster :: language/markdown.spec.js
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

const markdown = require('../../lib/language/markdown');

describe('language', () => {
  describe('markdown', () => {
    const options = ['-m'];
    const src = 'README.md';
    const dst = 'README.html';

    describe('.md2html()', () => {
      it('should notify "md2html not found"', () => {
        expect(markdown.md2html()).toBe(true);
        expect(os.whence).toHaveBeenLastCalledWith('md2html');
        expect(aster.notify).toHaveBeenLastCalledWith('failure', 'aster: md2html', 'md2html not found!');
      });

      it(`should execute \`md2html ${src} ${dst}\``, () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(markdown.md2html({ src })).toBe(false);
        expect(language.system).toHaveBeenLastCalledWith({
          args: ['md2html', src, dst],
          title: 'md2html',
          success: `${src}`,
          failure: `${src} failed`,
        });
      });

      it(`should execute \`md2html ${options.join(' ')} ${src} ${dst}\``, () => {
        os.whence.mockReturnValueOnce(true);
        language.system.mockReturnValueOnce(false);

        expect(markdown.md2html({ options, src, dst })).toBe(false);
        expect(language.system).toHaveBeenLastCalledWith({
          args: ['md2html', '-m', src, dst],
          title: 'md2html',
          success: `${src}`,
          failure: `${src} failed`,
        });
      });
    });
  });
});
