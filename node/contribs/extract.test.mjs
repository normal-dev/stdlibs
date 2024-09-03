import { strict as assert } from 'node:assert'
import { readFileSync } from 'node:fs'
import test, { describe } from 'node:test'
import { dirname } from 'path'
import { fileURLToPath } from 'url'
import extract from './extract.mjs'

const __dirname = dirname(fileURLToPath(import.meta.url))

// describe('ImportDefaultSpecifier', undefined, () => {
//   describe('global', undefined, () => {
//     test('MemberExpression', () => {
//       assert.deepEqual(extract(openTest('default/global/MemberExpression.js')), [
//         {
//           ident: 'node:assert.log',
//           line: 6
//         },
//         {
//           ident: 'node:fs.createReadStream',
//           line: 4
//         }
//       ])
//     })

//     test('ObjectProperty', () => {
//       assert.deepEqual(extract(openTest('default/global/ObjectProperty.js')), [
//         {
//           ident: 'node:assert.default',
//           line: 4
//         },
//         {
//           ident: 'node:assert.log',
//           line: 5
//         },
//         {
//           ident: 'node:assert.equal',
//           line: 6
//         },
//         {
//           ident: 'node:assert.notEqual',
//           line: 6
//         },
//         {
//           ident: 'node:assert.equal',
//           line: 7
//         }
//       ])
//     })

//     test('VariableDeclaration', () => {
//       assert.deepEqual(extract(openTest('default/global/VariableDeclaration.js')), [
//         {
//           ident: 'node:assert.default',
//           line: 3
//         }
//       ])
//     })
//   })

//   describe('local', undefined, () => {
//     test('ObjectProperty', () => {
//       assert.deepEqual(extract(openTest('default/local/ObjectProperty.js')), [{
//         ident: 'node:assert.default',
//         line: 5
//       }])
//     })

//     test('VariableDeclaration', () => {
//       assert.deepEqual(extract(openTest('default/local/VariableDeclaration.js')), [{
//         ident: 'node:assert.default',
//         line: 4
//       }])
//     })
//   })
// })

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

    // test('MemberExpression', () => {
    //   assert.deepEqual(extract(openTest('specifier/global/MemberExpression.js')), [
    //     {
    //       ident: 'node:crypto.constants',
    //       line: 3
    //     }
    //   ])
    // })

    test('ObjectProperty', () => {
      assert.deepEqual(extract(openTest('specifier/global/ObjectProperty.js')), [
        {
          ident: 'node:assert.equal',
          line: 5
        }
      ])
    })
  })
})

// describe('proprietary', undefined, () => {
//   test('build.js', () => {
//     assert.equal(extract(openTest('proprietary/build.js')), [
//       {
//         ident: 'node:fs/promises.mkdir',
//         line: 50
//       },
//       {
//         ident: 'node:fs/promises.rm',
//         line: 111
//       },
//       {
//         ident: 'node:fs/promises.readFile',
//         line: 160
//       },
//       {
//         ident: 'node:fs/promises.writeFile',
//         line: 175
//       },
//       {
//         ident: 'node:fs.existsSync',
//         line: 110
//       },
//       {
//         ident: 'node:fs.existsSync',
//         line: 157
//       },
//       {
//         ident: 'node:path.resolve',
//         line: 45
//       },
//       {
//         ident: 'node:path.resolve',
//         line: 101
//       },
//       {
//         ident: 'node:path.resolve',
//         line: 149
//       },
//       {
//         ident: 'node:path.basename',
//         line: 161
//       },
//       {
//         ident: 'node:path.resolve',
//         line: 176
//       },
//       {
//         ident: 'node:zlib.gzipSync',
//         line: 163
//       },
//       {
//         ident: 'node:zlib.brotliCompressSync',
//         line: 164
//       },
//       {
//         ident: 'node:cpus.default',
//         line: 79
//       },
//       {
//         ident: 'node:module.createRequire',
//         line: 32
//       }
//     ])
//   })
// })

const openTest = name => {
  const file = readFileSync(`${__dirname}/testfiles/${name}`, 'utf-8')
  return file.toString()
}
