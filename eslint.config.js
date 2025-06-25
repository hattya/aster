const { defineConfig } = require('eslint/config');
const js = require('@eslint/js');
const globals = require('globals');
const { FlatCompat } = require('@eslint/eslintrc');

const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

module.exports = defineConfig([{
  extends: compat.extends('eslint:recommended', 'airbnb-base'),
  languageOptions: {
    globals: {
      aster: 'readonly',
    },
    sourceType: 'script',
  },
}, {
  files: ['lib/**/*.js'],
  rules: {
    // Possible Errors
    'no-console': 'off',

    // Best Practices
    'no-else-return': 'off',
    'no-param-reassign': 'off',
    'vars-on-top': 'off',
    'yoda': 'off',

    // Strict Mode
    'strict': ['error', 'global'],

    // Variables
    'no-unused-vars': ['error', {
      'caughtErrors': 'none',
      'ignoreRestSiblings': true,
    }],

    // Stylistic Issues
    'func-names': 'off',
    'no-bitwise': 'off',
    'no-multi-assign': 'off',
    'no-nested-ternary': 'off',
    'no-plusplus': 'off',
    'no-underscore-dangle': 'off',
    'space-before-function-paren': ['error', 'never'],

    // ECMAScript 6
    'no-var': 'off',
    'object-shorthand': 'off',
    'prefer-arrow-callback': 'off',
    'prefer-destructuring': 'off',
    'prefer-rest-params': 'off',
    'prefer-spread': 'off',
    'prefer-template': 'off',

    // plugin: import
    'import/no-unresolved': 'off',
    'import/no-extraneous-dependencies': 'off',
  },
}, {
  files: ['test/**/*.js'],
  languageOptions: {
    globals: { ...globals.jest },
  },
  settings: {
    "import/core-modules": [
      "language",
    ],
  },
  rules: {
    // Stylistic Issues
    'no-plusplus': 'off',
  },
}]);
