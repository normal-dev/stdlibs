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

    const locus = []
    traverse(ast, {
      ImportDeclaration (path) {
        const { node: { source: { value: module } } } = path
        if (!stdlib.includes(module)) {
          return
        }

        resolveModule(ast, module, locus)
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
const isModuleImport = (path, module) => {
  const { scope, node } = path
  const binding = scope.getBinding(node.name)
  if ((!binding || binding.kind !== 'module')) return false
  if (binding.path.parent.source.value !== module) return false
  return true
}

/**
 * @param {NodePath} path
 * @returns {string}
 */
const resolveCanonicalName = path => {
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
 * @param {NodePath} path
 * @returns {boolean}
 */
const hasLocation = path => {
  return !!path.node.loc
}

/**
 * @param {Node} ast
 * @param {string} module
 * @param {Array<object>} apis
 */
const resolveModule = (ast, module, apis) => {
  try {
    traverse(ast, {
      Identifier (path) {
        if (isImported(path)) return
        if (!isModuleImport(path, module)) return
        if (!hasLocation(path)) return

        const { node, container } = path
        const { loc: { start: { line } }, name } = node

        if (Array.isArray(container)) {
          for (const node of container) {
            if (node.name !== name) continue

            const canonical = resolveCanonicalName(path)
            apis.push(newApi(module, canonical, line))
          }

          return
        }

        const { type, value } = container
        switch (type) {
          case 'CallExpression': {
            break
          }

          case 'MemberExpression': {
            const { property: { type } } = container
            const { node: { name }, scope } = path
            const binding = scope.getBinding(name)

            switch (binding.path.type) {
              case 'ImportSpecifier':
                apis.push(
                  newApi(module, name, line)
                )

                return
            }

            switch (type) {
              case 'Identifier': {
                const { property: { name } } = container
                apis.push(
                  newApi(module, name, line)
                )

                return
              }

              case 'MemberExpression': {
                const { property: { name } } = container
                apis.push(
                  newApi(module, name, line)
                )

                return
              }

              case 'StringLiteral': {
                const { property: { value } } = container
                apis.push(
                  newApi(module, value, line)
                )

                return
              }
            }

            break
          }

          case 'ObjectProperty':
            switch (value.type) {
              case 'Identifier':
                break
            }

            break
        }

        const canonical = resolveCanonicalName(path)
        apis.push(newApi(module, canonical, line))
      }

    })
  } catch (error) {
    console.warn(error)
  }
}

export default extract
