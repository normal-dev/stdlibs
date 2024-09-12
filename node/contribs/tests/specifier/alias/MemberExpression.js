import { exec as _exec, fork as f } from 'node:child_process'
import util from 'node:util'

const exec = util.promisify(_exec, console.log)

f('/dev/null')
