import babelParser from '@babel/parser'
import _babelTraverse from '@babel/traverse'
import { builtinModules as builtin } from 'node:module'

const babelTraverse = _babelTraverse.default

const eq = (a, b) => a === b

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
      ImportDeclaration(path) {
        const node = path.node
        const module = node.source.value // Save as [fs node:fs]
        // Exlcude non-built-in modules
        if (!publicBuiltinModules.includes(module)) {
          return
        }

        // For each import
        for (const specifier of node.specifiers) {
          // Some prevalidation
          switch (specifier.type) {
            // TODO: Support aliases: import { strict as assert } from 'node:assert'

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
              // apis.push({
              //   ident: `${module}.default`,
              //   line: node.loc.start.line,
              // })
              resolveImportDefaultSpec(ast, module, apis)

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

const resolveImportSpec = (node, apis) => {
}

const resolveImportDefaultSpec = (ast, module, apis) => {
  try {
    babelTraverse(ast, {
      // Class: console.Console

      // Object: path.posix

      // Literal: path.delim

      // Reference: assert.default()

      // Function: assert()
      CallExpression(path) {
        const node = path.node
        switch (node.callee.type) {
          case 'Identifier':
            // Check if callee is imported module name ("fs")
            // TODO: The imported identifier in a named import needs to be checked if it equals the callee
            if (eq(`node:${node.callee.name}`, module)) {
              if (!hasModuleBinding(node.path.scope)) {
                apis.push({
                  ident: `${module}.default`,
                  line: node.loc.start.line,
                })
              }
            }
        }
      }
    })
  } catch (error) {
    // This can fail
    console.warn(error)
  }
}

const hasModuleBinding = (scope, module) => {
  for (const [name, _] of Object.entries(scope.bindings)) {
    if (`node:${name}` === module) {
      return false
    }
  }
  return true
}


const resolveImportNamespaceSpec = (node, apis) => {
}

export default extract