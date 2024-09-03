import { parse } from '@babel/parser'
import * as t from '@babel/types'
import { createRequire } from 'module'
import { builtinModules as builtin } from 'node:module'
// Mimic import to preserve autocompletion
const require = createRequire(import.meta.url)
const traverse = require('@babel/traverse').default

export const stdlib = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

const newApi = (module, ident, line) => ({
  ident: `${module}.${ident}`,
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

    const contribs = []
    traverse(ast, {
      ImportDeclaration (path) {
        const { node: { source: { value: module }, specifiers } } = path
        if (!stdlib.includes(module)) {
          return
        }

        for (const { type } of specifiers) {
          switch (type) {
            // import fs from 'node:fs'
            case 'ImportDefaultSpecifier':
              resolveDefault(ast, module, contribs)
              break

            case 'ImportSpecifier':
              resolveSpecifier(ast, module, contribs)
              break
          }
        }
      }
    })

    return contribs
  } catch (error) {
    console.warn(error)
    return []
  }
}

const isModuleBinding = binding => {
  return (binding && binding.kind === 'module')
}

const isImport = path => {
  switch (path.parent.type) {
    case 'ImportDeclaration':
    case 'ImportDefaultSpecifier':
    case 'ImportNamespaceSpecifier':
    case 'ImportSpecifier':
      return true

    default:
      return false
  }
}

/**
 *
 * @param {t.Node} ast
 * @param {string} module
 * @param {string} ident
 * @param {Array<object>} apis
 */
const resolveDefault = (ast, module, apis) => {
  try {
    traverse(ast, {
      // Debug
      enter (path) {
        // console.debug(path)
      },

      Identifier (path) {
        /** @param {t.Identifier} node */
        const node = path.node

        if (isImport(path)) return

        const { scope } = path
        const binding = scope.getBinding(node.name)
        if (!isModuleBinding(binding)) return

        const { container } = path
        switch (container.type) {
          case 'CallExpression':
            break

          case 'MemberExpression': {
            const { property } = container
            switch (property.type) {
              case 'Identifier':
                apis.push(
                  newApi(module, property.name, node.loc.start.line)
                )
                break

              case 'StringLiteral':
                apis.push(
                  newApi(module, property.value, node.loc.start.line)
                )

                break
            }

            break
          }

          case 'ObjectProperty':
            switch (container.value.type) {
              case 'Identifier':
                apis.push(newApi(module, 'default', node.loc.start.line))
                break

              case 'MemberExpression':
                break
            }

            break

          case 'VariableDeclaration':
            break
        }
      }
    })
  } catch (error) {
    console.warn(error)
  }
}

/**
 *
 * @param {t.Node} ast
 * @param {string} module
 * @param {string} ident
 * @param {Array<object>} apis
 */
const resolveSpecifier = (ast, module, apis) => {
  try {
    traverse(ast, {
      // Debug
      enter (path) {
        // console.debug(path)
      },

      Identifier (path) {
        /** @param {t.Identifier} node */
        const node = path.node

        if (isImport(path)) return

        const { scope } = path
        const binding = scope.getBinding(node.name)
        if (!isModuleBinding(binding)) return

        const { container } = path
        switch (container.type) {
          case 'CallExpression': {
            const { callee: { name } } = container
            apis.push(newApi(module, name, node.loc.start.line))

            break
          }

          case 'MemberExpression': {
            break
          }

          case 'ObjectProperty':

            break

          case 'VariableDeclaration':
            break
        }
      }
    })
  } catch (error) {
    console.warn(error)
  }
}

export default extract
