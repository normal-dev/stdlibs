import babelParser from '@babel/parser'
import _babelTraverse from '@babel/traverse'
// TODO: Import "@babel/types" for "ImportDefaultSpecifier"
import { builtinModules as builtin } from 'node:module'

const babelTraverse = _babelTraverse.default

export const publicBuiltinModules = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

const extract = src => {
  try {
    const ast = babelParser.parse(src, {
      allowAwaitOutsideFunction: true,
      allowImportExportEverywhere: true,
      allowNewTargetOutsideFunction: true,
      allowReturnOutsideFunction: true,
      allowSuperOutsideMethod: true,
      allowUndeclaredExports: true,
      errorRecovery: false,
      strictMode: false,
      sourceType: 'unambiguous',
      plugins: ['typescript']
    })

    const apis = []
    babelTraverse(ast, {
      // import fs from 'node:fs'
      ImportDeclaration: path => {
        const node = path.node
        const module = node.source.value
        if (!publicBuiltinModules.includes(module)) {
          return
        }

        for (const specifier of node.specifiers) {
          switch (specifier.type) {
            // TODO: import { strict as assert } from 'node:assert'
            // TODO: Merge "import * as path from 'node:path'"

            // import { readFile } from 'node:fs'
            case 'ImportSpecifier':
              switch (specifier.imported.type) {
                case 'Identifier':
                  apis.push({
                    ident: `${module}.${specifier.imported.name}`,
                    line: node.loc.start.line,
                  })

                  break
              }

              break

            // import fs from 'node:fs'
            case 'ImportDefaultSpecifier':
              apis.push({
                ident: `${module}.default`,
                line: node.loc.start.line,
              })

              break

            // import * as path from 'node:fs'
            case 'ImportNamespaceSpecifier':
              apis.push({
                ident: module,
                line: node.loc.start.line,
              })

              break
          }
        }
      }
    })

    return apis
  } catch (error) {
    // This can fail
    return []
  }
}

export default extract