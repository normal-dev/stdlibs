module.exports = {
  root: true,
  env: {
    node: true
  },

  extends: [
    'plugin:vue/recommended',
    '@vue/standard'
  ],

  rules: {
    camelcase: 'off',
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-empty-pattern': 0,
    'new-cap': 0,

    'vue/html-closing-bracket-newline': ['error', {
      singleline: 'never',
      multiline: 'never'
    }],
    'vue/attributes-order': 'error',
    'vue/max-attributes-per-line': ['error', {
      singleline: 1,
      multiline: 1
    }],
    'vue/html-indent': ['error', 2],
    'vue/order-in-components': ['error'],
    'vue/require-default-prop': 0,
    'vue/no-multiple-template-root': 'off',
    'vue/no-v-model-argument': 'off'
  }
}
