import assert from 'node:assert'

const o = {
  a: assert,
  b: assert.log,
  c: true ? assert.equal : assert.notEqual,
  d: assert['equal']
}
