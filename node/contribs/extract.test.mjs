import { strict as assert } from 'node:assert'
import { readFileSync } from 'node:fs'
import test, { describe } from 'node:test'
import { dirname } from 'path'
import { fileURLToPath } from 'url'
import extract from './extract.mjs'

const __dirname = dirname(fileURLToPath(import.meta.url))

describe('ImportDefaultSpecifier', undefined, () => {
  describe('local', undefined, () => {
    describe('shadow', undefined, () => {
      test('VariableDeclaration', () => {
        assert.deepEqual(extract(openTest('default/local/shadow/VariableDeclaration.js')), [])
      })
    })
  })

  // describe('global', undefined, () => {
  //   test('CallExpression', () => {
  //     assert.deepEqual(extract(openTest('default/global/CallExpression.js')), [
  //       {
  //         ident: 'node:assert.default',
  //         line: 3
  //       }
  //     ])
  //   })

  //   test('ObjectProperty', () => {
  //     assert.deepEqual(extract(openTest('default/global/ObjectProperty.js')), [
  //       {
  //         ident: 'node:assert.default',
  //         line: 4
  //       }
  //     ])
  //   })

  //   test('VariableDeclaration', () => {
  //     assert.deepEqual(extract(openTest('default/global/VariableDeclaration.js')), [
  //       {
  //         ident: 'node:assert.default',
  //         line: 3
  //       }
  //     ])
  //   })
  // })
})

const openTest = name => {
  try {
    const file = readFileSync(`${__dirname}/testfiles/${name}`, 'utf-8')
    return file.toString()
  } catch (error) {
    throw new Error(error)
  }
}
