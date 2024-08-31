import assert from 'node:assert'

function a (b) {
  assert.log()
  function c (assert) {
    assert.log()
  }
}
