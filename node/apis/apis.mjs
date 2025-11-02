import { builtinModules as builtin } from 'node:module'

const stdlib = builtin
  .filter(module => !module.startsWith('_'))
  .map(module => `node:${module}`)

// Modules, which can't be imported, are added to the Docker image
let i = stdlib.length
while (i--) {
  try {
    await import(stdlib[i])
  } catch (error) {
    console.warn(error)
    stdlib.splice(i, 1)
  }
}

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
