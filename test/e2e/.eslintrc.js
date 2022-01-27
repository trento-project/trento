/* eslint-disable no-undef */
module.exports = {
  env: {
    browser: true,
    es2021: true,
    jquery: true,
    node: true,
    'cypress/globals': true,
  },
  extends: ['eslint:recommended', 'plugin:cypress/recommended'],
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 12,
    sourceType: 'module',
  },
  plugins: ['cypress'],
  rules: {
    'no-console': 'error',
    'no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
  },
};
