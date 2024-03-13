import { builtinModules as builtin } from 'node:module'
import { version } from 'node:process'
import mongoClient from './db.mjs'

// MongoDB Id
const CAT_ID = '_cat'

const mongoCollection = mongoClient.db('apis').collection('node')
const publicBuiltinModules = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

const nodeJsApiDocumentationCurry = async () => {
  console.debug('caching Node.js Apis documentation...')

  const response = await fetch('https://nodejs.org/docs/latest-v20.x/api/all.json')
  const { modules } = await response.json()

  return api => {
    for (const module of modules) {
      if (!publicBuiltinModules.includes(`node:${module.name}`)) {
        continue
      }

      // Default imports
      if (api === `node:${module.name}`) {
        return module.desc
      }
      // Classes
      for (const { name, desc } of module.classes ?? []) {
        if (`node:${name}` === api) {
          return desc
        }
      }
      // Methods
      for (const { name, desc } of module.methods ?? []) {
        if (`node:${module.name}.${name}` === api) {
          return desc
        }
      }
    }

    return null
  }
}

const isTypeOfClass = value => {
  return typeof value === 'function' && /^\s*class\s+/.test(value.toString())
}

try {
  console.debug('using Node.js version %s', version)

  console.debug('cleaning...')
  await mongoCollection.deleteMany({})

  const getApiDocumentation = await nodeJsApiDocumentationCurry()

  let amountApis = 0
  for (const module of publicBuiltinModules) {
    console.debug('processing module %s...', module)

    // Default imports: "import fs from 'node:fs'"
    const documentation = getApiDocumentation(module)
    await mongoCollection.insertOne({
      _id: `${module}.default`,
      doc: documentation,
      name: 'default',
      ns: module,
      type: 'module',
    })

    for (const [name, value] of Object.entries(await import(module))) {
      if (!name) {
        continue
      }
      // Ignore private and default modules
      if (name.startsWith('_') || name === 'default') {
        continue
      }

      const apiDocumentation = getApiDocumentation(`${module}.${name}`)
      // Normal imports: import { readFileSync } from 'fs'
      await mongoCollection.insertOne({
        _id: `${module}.${name}`,
        doc: apiDocumentation,
        name,
        ns: module,
        type: isTypeOfClass(value) ? 'class' : typeof value,
      })

      // TODO: Include nested modules

      amountApis++
    }
  }

  console.debug('saving catalogue...')
  await mongoCollection.insertOne({
    _id: CAT_ID,
    n_apis: amountApis,
    n_ns: publicBuiltinModules.length,
    ns: publicBuiltinModules,
    version: version.substring(1),
    vids: {
      'node:assert': '4Vnn8JUyotw'
    }
  })

  process.exit(0)
} catch (error) {
  console.error(error)
  process.exit(1)
}
