import { parse } from '@babel/parser'
import * as t from '@babel/types'
import { createRequire } from 'module'
import { builtinModules as builtin } from 'node:module'

const require = createRequire(import.meta.url)
const traverse = require('@babel/traverse').default

export const stdlib = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

const newApi = (ident, line, _default = false) => ({
  ident: `${ident}${_default ? '.default' : ''}`,
  line
})

const extract = src => {
  try {
    const ast = parse(src, {
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
    traverse(ast, {
      ImportDeclaration (path) {
        const { node: { source: { value: module }, specifiers } } = path
        if (!stdlib.includes(module)) {
          return
        }

        for (const { type, local: { name: ident } } of specifiers) {
          switch (type) {
            // import fs from 'node:fs'
            case 'ImportDefaultSpecifier':
              resolveDefault(ast, module, ident, apis)

              break
          }
        }
      }
    })

    return apis
  } catch (error) {
    console.warn(error)
    return []
  }
}

/**
 *
 * @param {t.Node} ast
 * @param {string} module
 * @param {string} ident
 * @param {Array<object>} apis
 */
const resolveDefault = (ast, module, ident, apis) => {
  try {
    traverse(ast, {
      // Debug
      enter (path) {
        // console.debug(path)
      },

      // assert.log()
      CallExpression (path) {
        /** @param {t.CallExpression} node */
        const node = path.node

        if (isGlobalScope(path)) {
          if (node.callee?.object.name === ident) {
            apis.push(
              newApi(
                module,
                node.loc.start.line,
                true
              )
            )
          }
        } else {
          const obj = node.callee.object.name
          const binding = path.scope.bindings
          if (binding[obj] === ident && binding[obj].kind !== 'param') {
            console.log(obj)
          }
        }
      },

      // const obj = {
      //  a: assert
      // }
      ObjectProperty (path) {
        /** @param {t.ObjectProperty} node */
        const node = path.node

        if (isGlobalScope(path)) {
          if (node?.value?.name === ident) {
            apis.push(
              newApi(
                module,
                path.node.key.loc.start.line,
                true
              )
            )
          }
        }
      },

      // const a = assert
      VariableDeclaration (path) {
        /** @param {t.VariableDeclaration} node */
        const node = path.node

        if (isGlobalScope(path)) {
          // Global

          for (const decl of node.declarations) {
            if (decl?.init?.name === ident) {
              apis.push(
                newApi(
                  module,
                  decl.loc.start.line,
                  true
                )
              )
            }
          }
        }
      }
    })
  } catch (error) {
    console.warn(error)
  }
}

const isGlobalScope = path => {
  const { scope: { parent } } = path
  return parent === undefined
}

export default extract
