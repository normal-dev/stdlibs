import { builtinModules as builtin } from 'node:module'
import { version } from 'node:process'
import mongoClient from '../mongo/client.mjs'
import apis from './apis.mjs'

// MongoDB Id
const CAT_ID = '_cat'

const mongoColl = mongoClient.db('apis').collection('node')

console.debug('version: %s', version)

await mongoColl.deleteMany({})

const defaults = new Set()
let apisn = 0
for (const [api, type] of apis) {
  console.debug('api: %s', api)

  const module = api.split('.')[0]
  const ident = api.split('.')[1]
  console.debug('namespace: %s', module)
  console.debug('name: %s', ident)
  console.debug('type: %s', type)

  if (!defaults.has(module)) {
    await mongoColl.insertOne({
      _id: `${module}.default`,
      doc: '',
      name: 'default',
      ns: module,
      type: 'module'
    })

    defaults.add(module)
  }

  await mongoColl.insertOne({
    _id: `${module}.${ident}`,
    doc: '',
    name: ident,
    ns: module,
    type
  })

  apisn++
}

await mongoColl.insertOne({
  _id: CAT_ID,
  n_apis: apisn,
  n_ns: apis.length,
  ns: apis,
  version: version.substring(1)
})

process.exit(0)
