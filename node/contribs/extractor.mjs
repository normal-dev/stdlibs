import { parse } from '@babel/parser'
// eslint-disable-next-line no-unused-vars
import { NodePath } from '@babel/traverse'
import babelTypes from '@babel/types'
import { createRequire } from 'module'
import { builtinModules as builtin } from 'node:module'

// eslint-disable-next-line no-unused-vars
const { Node, Identifier } = babelTypes

const require = createRequire(import.meta.url)
const traverse = require('@babel/traverse').default

export const stdlib = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

const newLocus = (module, ident, line, _) => ({
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

    const locus = []
    traverse(ast, {
      ImportDeclaration (path) {
        const { node: { source: { value: module } } } = path
        if (!stdlib.includes(module)) {
          return
        }

        findLocus(ast, module, locus)
      }
    })

    return locus
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
const isModuleIdent = (path, module) => {
  const { scope, node: { name } } = path
  const binding = scope.getBinding(name)
  if ((!binding || binding.kind !== 'module')) {
    return false
  }
  if (binding.path.parent.source.value !== module) {
    return false
  }

  return true
}

/**
 * @param {NodePath} path
 * @returns {boolean}
 */
const isImportNode = path => {
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
 * @param {NodePath} path
 * @returns {boolean}
 */
const hasLocation = path => {
  return !!path.node.loc
}

/**
 * @param {NodePath} path
 * @param {string} module
 * @returns {boolean}
 */
const isResolveable = (path, module) => {
  if (isImportNode(path)) {
    return false
  }
  if (!isModuleIdent(path, module)) {
    return false
  }
  if (!hasLocation(path)) {
    return false
  }

  return true
}

/**
 * @param {NodePath} path
 * @returns {string}
 */
const findCanonical = path => {
  const { node: { name }, scope } = path
  const binding = scope.getBinding(name)
  switch (binding.path.node.type) {
    case 'ImportNamespaceSpecifier':
    case 'ImportDefaultSpecifier':
      return 'default'

    case 'ImportSpecifier':
      break
  }

  return binding.path.node.imported.name
}

/**
 * @param {Node} ast
 * @param {string} module
 * @param {Array<object>} locus
 */
const findLocus = (ast, module, locus) => {
  try {
    traverse(ast, {
      Identifier (path) {
        if (!isResolveable(path, module)) {
          return
        }

        const { node, container } = path
        const { loc: { start: { column, line } } } = node

        if (Array.isArray(container)) {
          const canonical = findCanonical(path)
          locus.push(newLocus(module, canonical, line, column))

          return
        }

        const { type, value } = container
        switch (type) {
          case 'CallExpression':
          case 'ConditionalExpression':
            break

          case 'MemberExpression': {
            const { property: { type } } = container
            const { node: { name }, scope } = path
            const binding = scope.getBinding(name)

            switch (binding.path.type) {
              case 'ImportSpecifier':
                locus.push(
                  newLocus(module, name, line, column)
                )

                return
            }

            switch (type) {
              case 'Identifier':
              case 'MemberExpression': {
                const { property: { name } } = container
                locus.push(
                  newLocus(module, name, line, column)
                )

                return
              }

              case 'StringLiteral': {
                const { property: { value } } = container
                locus.push(
                  newLocus(module, value, line, column)
                )

                return
              }
            }

            break
          }

          case 'ObjectProperty': {
            switch (path.key) {
              case 'key':
                return
            }

            switch (value.type) {
              case 'Identifier':
                break
            }
          }
        }

        const canonical = findCanonical(path)
        locus.push(newLocus(module, canonical, line, column))
      }
    })
  } catch (error) {
    console.warn(error)
  }
}

export default extract
