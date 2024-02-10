import test, { describe } from 'node:test'
import { strict as assert } from 'node:assert'
import { readFileSync } from 'node:fs'
import { dirname } from 'path'
import { fileURLToPath } from 'url'
import extract from './extract.mjs'

const __dirname = dirname(fileURLToPath(import.meta.url))

const openTest = name => {
  try {
    const file = readFileSync(`${__dirname}/testfiles/${name}.txt`, 'utf-8')
    return file.toString()
  } catch (error) {
    throw new Error(error)
  }
}

describe('extract', undefined, () => {
  test('vuejs/core/scripts/build.js', () => {
    assert.deepEqual(extract(openTest('build.js')), [
      {
        ident: 'node:fs/promises.default',
        line: 19
      },
      {
        ident: 'node:fs.existsSync',
        line: 20
      },
      {
        ident: 'node:fs.readFileSync',
        line: 20
      },
      {
        ident: 'node:path.default',
        line: 21
      },
      {
        ident: 'node:zlib.gzipSync',
        line: 23
      },
      {
        ident: 'node:zlib.brotliCompressSync',
        line: 23
      },
      {
        ident: 'node:os.cpus',
        line: 26
      },
      {
        ident: 'node:module.createRequire',
        line: 27
      }
    ])
  })

  test('Automattic/wp-calypso/bin/did-calypso-app-change.mjs', () => {
    assert.deepEqual(extract(openTest('did-calypso-app-change.mjs')), [
      {
        ident: 'node:child_process.exec',
        line: 1
      },
      {
        ident: 'node:fs.createWriteStream',
        line: 2
      },
      {
        ident: 'node:stream.Readable',
        line: 3
      },
      {
        ident: 'node:stream/promises.finished',
        line: 4
      },
      {
        ident: 'node:util.default',
        line: 5
      },
    ])
  })

  test('sveltejs/svelte/sites/svelte.dev/vite.config.js', () => {
    assert.deepEqual(extract(openTest('vite.config.js')), [
      {
        ident: 'node:fs/promises.readFile',
        line: 3
      },
    ])
  })

  test('zx/src/goods.ts', () => {
    assert.deepEqual(extract(openTest('goods.ts')), [
      {
        ident: 'node:assert.default',
        line: 15
      },
      {
        ident: 'node:readline.createInterface',
        line: 19
      }
    ])
  })
})
