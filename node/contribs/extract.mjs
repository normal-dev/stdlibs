import { parse } from '@babel/parser'
import { NodePath } from '@babel/traverse'
import babelTypes from '@babel/types'
import { createRequire } from 'module'
import { builtinModules as builtin } from 'node:module'
const { Node, Identifier } = babelTypes

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
        const { node: { source: { value: module } } } = path
        if (!stdlib.includes(module)) {
          return
        }

        resolveModule(ast, module, contribs)
      }
    })

    return contribs
  } catch (error) {
    console.warn(error)
    return []
  }
}

/**
 * @param {NodePath} path
 * @param {string} module
 * @returns {boolean}
 */
const isModuleBinding = (path, module) => {
  const { scope, node } = path
  const binding = scope.getBinding(node.name)
  if ((!binding || binding.kind !== 'module')) return false
  if (binding.path.parent.source.value !== module) return false
  return true
}

/**
 * @param {NodePath} path
 * @returns {boolean}
 */
const isImported = path => {
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
 * @param {NodePath} path
 * @returns {boolean}
 */
const hasLocation = path => {
  return !!path.node.loc
}

/**
 *
 * @param {Node} ast
 * @param {string} module
 * @param {Array<object>} apis
 */
const resolveModule = (ast, module, apis) => {
  try {
    traverse(ast, {
      Identifier (path) {
        if (isImported(path)) return
        if (!isModuleBinding(path, module)) return
        if (!hasLocation(path)) return

        const { scope, node } = path
        const binding = scope.getBinding(node.name)
        switch (binding.path.type) {
          case 'ImportDefaultSpecifier':
            resolveDefault(path, module, apis)
            break

          case 'ImportNamespaceSpecifier':
            resolveDefault(path, module, apis)
            break

          case 'ImportSpecifier':
            resolveSpecifier(path, module, apis)
            break
        }
      }

    })
  } catch (error) {
    console.warn(error)
  }
}

/**
 * @param {NodePath} path
 * @param {string} module
 * @param {Array<object>} apis
 */
const resolveDefault = (path, module, apis) => {
  try {
    /** @param {Identifier} node */
    const { node } = path
    const { loc: { start: { line } } } = node

    const { container } = path
    switch (container.type) {
      case 'CallExpression':
        break

      case 'MemberExpression': {
        const { property } = container
        switch (property.type) {
          case 'Identifier':
            apis.push(
              newApi(module, property.name, line)
            )
            break

          case 'StringLiteral':
            apis.push(
              newApi(module, property.value, line)
            )

            break
        }

        break
      }

      case 'ObjectProperty':
        switch (container.value.type) {
          case 'Identifier':
            apis.push(newApi(module, 'default', line))
            break
        }

        break

      case 'VariableDeclarator':
        apis.push(newApi(module, 'default', line))

        break
    }
  } catch (error) {
    console.warn(error)
  }
}

/**
 * @param {NodePath} path
 * @param {string} module
 * @param {Array<object>} apis
 */
const resolveSpecifier = (path, module, apis) => {
  try {
    /** @param {Identifier} node */
    const { node } = path
    const { loc: { start: { line } }, name } = node

    const { container } = path
    switch (container.type) {
      case 'CallExpression': {
        const { callee: { name } } = container
        apis.push(newApi(module, name, line))

        break
      }

      case 'MemberExpression': {
        switch (container.type) {
          case 'Identifier':
            break

          case 'StringLiteral':
            break

          case 'MemberExpression': {
            if (path.key === 'object') {
              apis.push(newApi(module, name, line))
            }
          }
        }

        break
      }

      case 'ObjectProperty':
        apis.push(newApi(module, name, line))
        break

      case 'VariableDeclaration':
        break
    }
  } catch (error) {
    console.warn(error)
  }
}

export default extract
