import * as fs from 'node:fs'
import * as http from 'node:http'

fs.existsSync('/dev/null')
const agent = new http.Agent()
