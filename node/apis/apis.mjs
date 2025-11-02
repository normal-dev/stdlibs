import { builtinModules as builtin } from 'node:module'

const stdlib = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)
  // Modules, which can't be imported, are added to the Docker image
  .filter(async module => {
    try {
      await import(module)
      return true
    } catch (error) {
      console.warn(error)
      return false
    }
  })

const isTypeOfClass = value => {
  return typeof value === 'function' && /^\s*class\s+/.test(value.toString())
}

const apis = new Map()
for (const module of stdlib) {
  for (const [name, value] of Object.entries(await import(module))) {
    if (!name) {
      continue
    }
    // Ignore private and default modules
    if (name.startsWith('_') || name === 'default') {
      continue
    }

    apis.set(`${module}.${name}`, isTypeOfClass(value) ? 'class' : typeof value)
  }
}

export default apis
