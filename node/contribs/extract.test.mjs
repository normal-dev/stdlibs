import { strict as assert } from 'node:assert'
import { readFileSync } from 'node:fs'
import test, { describe } from 'node:test'
import { dirname } from 'path'
import { fileURLToPath } from 'url'
import extract from './extract.mjs'

const __dirname = dirname(fileURLToPath(import.meta.url))

describe('ImportDefaultSpecifier', undefined, () => {
  describe('global', undefined, () => {
    // test('MemberExpression', () => {
    //   assert.deepEqual(extract(openTest('default/global/MemberExpression.js')), [
    //     {
    //       ident: 'node:assert.log',
    //       line: 3
    //     }
    //   ])
    // })

    test('ObjectProperty', () => {
      assert.deepEqual(extract(openTest('default/global/ObjectProperty.js')), [
        {
          ident: 'node:assert.default',
          line: 4
        },
        {
          ident: 'node:assert.log',
          line: 5
        },
        {
          ident: 'node:assert.equal',
          line: 6
        },
        {
          ident: 'node:assert.notEqual',
          line: 6
        }
      ])
    })

    // test('VariableDeclaration', () => {
    //   assert.deepEqual(extract(openTest('default/global/VariableDeclaration.js')), [
    //     {
    //       ident: 'node:assert.default',
    //       line: 3
    //     }
    //   ])
    // })
    // })

    // describe('local', undefined, () => {
    //   describe('param', undefined, () => {
    //     test('CallExpression', () => {
    //       assert.deepEqual(extract(openTest('default/local/param/CallExpression.js')), [{
    //         ident: 'node:assert.equal',
    //         line: 4
    //       }])
    //     })
    //   })

    //   test('ObjectProperty', () => {
    //     assert.deepEqual(extract(openTest('default/local/ObjectProperty.js')), [{
    //       ident: 'node:assert.default',
    //       line: 5
    //     }])
    //   })

    //   test('VariableDeclaration', () => {
    //     assert.deepEqual(extract(openTest('default/local/VariableDeclaration.js')), [{
    //       ident: 'node:assert.default',
    //       line: 4
    //     }])
    //   })
  })
})

describe('ImportSpecifier', undefined, () => {
  describe('global', undefined, () => {
    // test('CallExpression', () => {
    //   assert.deepEqual(extract(openTest('specifier/global/CallExpression.js')), [
    //     {
    //       ident: 'node:assert.equal',
    //       line: 3
    //     }
    //   ])
    // })
  })
})

const openTest = name => {
  const file = readFileSync(`${__dirname}/testfiles/${name}`, 'utf-8')
  return file.toString()
}
