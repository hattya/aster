//
// aster :: language.spec.js
//
//   Copyright (c) 2020-2024 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//
//

global.aster = {
  notify: jest.fn(),
};

jest.mock('os');

const os = require('os');
const path = require('path');
const process = require('process');

os.getwd = jest.fn().mockReturnValue(path.join(process.cwd(), 'aster'));
os.system = jest.fn();

const language = require('../lib/language');

describe('language', () => {
  describe('.prefix', () => {
    it('should be "aster: "', () => {
      expect(language.prefix).toBe('aster: ');
    });
  });

  describe('.prompt', () => {
    it('is "> "', () => {
      expect(language.prompt).toBe('> ');
    });
  });

  describe('.system()', () => {
    const obj = {
      args: ['jest', '--coverage'],
      options: {},
      title: 'test',
      success: 'passed',
      failure: 'failed',
    };

    it('should notify success', () => {
      const spy = jest.spyOn(console, 'log').mockImplementation(() => {});
      os.system.mockReturnValueOnce(false);

      language.system(obj);
      expect(spy).lastCalledWith(`> ${obj.args.join(' ')}`);
      expect(os.system).lastCalledWith(obj.args, obj.options);
      expect(aster.notify).lastCalledWith('success', `aster: ${obj.title}`, obj.success);

      spy.mockRestore();
    });

    it('should notify failure', () => {
      const spy = jest.spyOn(console, 'log').mockImplementation((x) => x);
      os.system.mockReturnValueOnce(true);

      language.system({
        args: ['jest', '--coverage'],
        options: {},
        title: 'test',
        success: 'passed',
        failure: 'failed',
      });
      expect(spy).lastCalledWith(`> ${obj.args.join(' ')}`);
      expect(os.system).lastCalledWith(obj.args, obj.options);
      expect(aster.notify).lastCalledWith('failure', `aster: ${obj.title}`, obj.failure);

      spy.mockRestore();
    });
  });
});
