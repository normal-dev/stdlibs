import assert from 'node:assert'
import { createReadStream } from 'node:fs'

createReadStream('/dev/null')

assert.log()
