import { execSync } from 'node:child_process'
import { readdirSync, readFileSync, rmSync, statSync } from 'node:fs'
import path from 'node:path'
import { Octokit } from 'octokit'
import mongoClient from './db.mjs'
import extract from './extract.mjs'

// MongoDB Ids
const CAT_ID = '_cat'
const LICENSES_ID = '_licenses'

const mongoCollection = mongoClient.db('contribs').collection('node')

const findNodeJsFiles = async (dir, files) => {
  for (const filename of readdirSync(dir)) {
    try {
      const stat = path.join(dir, filename)

      const stats = statSync(stat)
      if (stats.isDirectory()) {
        const dir = stat
        files = await findNodeJsFiles(dir, files)
        continue
      }

      const file = stat
      switch (path.extname(filename)) {
        case '.js':
        case '.mjs':
        case '.cjs':
        case '.ts':
          files.push(file)
      }
    } catch (error) {
      // Symbolic links result into an error and can be ignored
      console.warn(error)
    }
  }

  return files
}

const insertLicenses = async () => {
  await mongoCollection.deleteOne({ _id: LICENSES_ID })
  const insertOneResult = await mongoCollection.insertOne({
    _id: LICENSES_ID,
    repos: [
      {
        author: 'Guillermo Rauch and Socket.IO contributors',
        repo: ['socketio', 'socket.io'],
        type: 'MIT license'
      },
      {
        author: 'OpenJS Foundation',
        repo: ['jquery', 'jquery'],
        type: 'MIT license'
      },
      {
        author: 'Transloadit',
        repo: ['transloadit', 'uppy'],
        type: 'MIT license'
      },
      {
        author: 'Sequelize Authors',
        repo: ['sequelize', 'sequelize'],
        type: 'MIT license'
      },
      {
        author: 'OpenJS Foundation',
        repo: ['appium', 'appium'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Google, Inc.',
        repo: ['puppeteer', 'puppeteer'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Sindre Sorhus',
        repo: ['avajs', 'ava'],
        type: 'MIT license'
      },
      {
        author: 'PlayCanvas Ltd.',
        repo: ['playcanvas', 'engine'],
        type: 'MIT license'
      },
      {
        author: 'Tutao GmbH',
        repo: ['tutao', 'tutanota'],
        type: 'GNU General Public License v3.0'
      },
      {
        author: 'Highlight Inc.',
        repo: ['highlight', 'highlight'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Supabase Inc',
        repo: ['supabase', 'supabase'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Herbie Project Modified work, Google Inc.',
        repo: ['herbie-fp', 'herbie'],
        type: 'Copyright notice'
      },
      {
        author: 'Google LLC.',
        repo: ['angular', 'angular'],
        type: 'MIT license'
      },
      {
        author: '650 Industries, Inc.',
        repo: ['expo', 'expo'],
        type: 'MIT license'
      },
      {
        author: 'Deno authors',
        repo: ['denoland', 'deno'],
        type: 'MIT license'
      },
      {
        author: 'Yuxi (Evan) You and Vite contributors',
        repo: ['vitejs', 'vite'],
        type: 'MIT license'
      },
      {
        author: 'Noel Buechler, Amish Shah',
        repo: ['discordjs', 'discord.js'],
        type: 'MIT license'
      },
      {
        author: 'Printio (Juriy Zaytsev, Maxim Chernyak)',
        repo: ['fabricjs', 'fabric.js'],
        type: 'Copyright notice'
      },
      {
        author: 'Yuxi (Evan) You',
        repo: ['vuejs', 'core'],
        type: 'MIT license'
      },
      {
        author: 'MarkedJS, Christopher Jeffrey',
        repo: ['markedjs', 'marked'],
        type: 'MIT license'
      },
      {
        author: 'Apify Technologies s.r.o.',
        repo: ['apify', 'crawlee'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Electron contributors, GitHub Inc.',
        repo: ['electron', 'electron'],
        type: 'MIT license'
      },
      {
        author: 'Zach Leatherman',
        repo: ['11ty', 'eleventy'],
        type: 'MIT license'
      },
      {
        author: 'Dagger, Inc.',
        repo: ['dagger', 'dagger'],
        type: 'Apache License 2.0'
      },
      {
        author: 'The Cheerio contributors',
        repo: ['cheeriojs', 'cheerio'],
        type: 'MIT license'
      },
      {
        author: 'Automattic Inc.',
        repo: ['Automattic', 'wp-calypso'],
        type: 'GNU General Public License v2.0'
      },
      {
        author: 'Temporal Technologies Inc.',
        repo: ['temporalio', 'sdk-typescript'],
        type: 'MIT license'
      },
      {
        author: 'Andrey Okonetchnikov',
        repo: ['lint-staged', 'lint-staged'],
        type: 'MIT license'
      },
      {
        author: 'Metarhia contributors',
        repo: ['metarhia', 'impress'],
        type: 'MIT license'
      },
      {
        author: 'Bun',
        repo: ['oven-sh', 'bun'],
        type: 'MIT license'
      },
      {
        author: 'LongYinan',
        repo: ['napi-rs', 'napi-rs'],
        type: 'MIT license'
      },
      {
        author: 'PingCAP',
        repo: ['pingcap', 'ossinsight'],
        type: 'Apache License 2.0'
      },
      {
        author: 'Sindre Sorhus',
        repo: ['sindresorhus', 'got'],
        type: 'MIT license'
      },
      {
        author: 'The Bootstrap Authors',
        repo: ['twbs', 'bootstrap'],
        type: 'MIT license'
      },
      {
        author: 'Svelte contributors',
        repo: ['sveltejs', 'svelte'],
        type: 'MIT license'
      },
      {
        author: 'midwayjs',
        repo: ['midwayjs', 'midway'],
        type: 'MIT license'
      },
      {
        author: 'The Fastify Team',
        repo: ['fastify', 'fastify'],
        type: 'MIT license'
      },
      {
        author: 'GitHub Inc',
        repo: ['hubotio', 'hubot'],
        type: 'MIT license'
      },
      {
        author: 'James Long and contributors',
        repo: ['prettier', 'prettier'],
        type: 'MIT license'
      },
      {
        author: 'Vercel, Inc.',
        repo: ['vercel', 'next.js'],
        type: 'MIT license'
      },
      {
        author: 'Biome Developers and Contributors',
        repo: ['biomejs', 'biome'],
        type: 'MIT license'
      },
      {
        author: 'Microsoft Corporation',
        repo: ['microsoft', 'vscode'],
        type: 'MIT license'
      },
      {
        author: 'THE ETHERPAD FOUNDATION',
        repo: ['ether', 'etherpad-lite'],
        type: 'Apache-2.0 license'
      },
      {
        author: 'Monospace, Inc.',
        repo: ['directus', 'directus'],
        type: 'GNU General Public License v3.0'
      },
      {
        author: 'Spacedrive Technology Inc.',
        repo: ['spacedriveapp', 'spacedrive'],
        type: 'AGPL-3.0 license'
      },
      {
        author: 'The Backstage Authors',
        repo: ['backstage', 'backstage'],
        type: 'Apache-2.0 license'
      },
      {
        author: 'Peter Hedenskog',
        repo: ['sitespeedio', 'sitespeed.io'],
        type: 'MIT license'
      },
      {
        author: 'OpenJS Foundation and other contributors',
        repo: ['webdriverio', 'webdriverio'],
        type: 'MIT license'
      },
      {
        author: 'Streetwriters (Private) Ltd.',
        repo: ['streetwriters', 'notesnook'],
        type: 'GPL-3.0 license'
      },
      {
        author: 'Nomic, Inc.',
        repo: ['nomic-ai', 'gpt4all'],
        type: 'MIT license'
      },
      {
        author: 'Luciano Mammino, will Farrell and the Middy team',
        repo: ['middyjs', 'middy'],
        type: 'MIT license'
      },
      {
        author: 'Salesforce.com, Inc.',
        repo: ['salesforce', 'lwc'],
        type: 'MIT license'
      },
      {
        author: 'Botpress Technologies, Inc.',
        repo: ['botpress', 'botpress'],
        type: 'MIT license'
      }
    ]
  })
  return insertOneResult.acknowledged
}

const insertContribs = async contributions => {
  const insertManyResult = await mongoCollection.insertMany(contributions)
  return insertManyResult.insertedCount
}

const insertCatalogue = async (contribsn, reposn) => {
  await mongoCollection.deleteOne({ _id: CAT_ID })
  const insertOneResult = await mongoCollection.insertOne({
    _id: CAT_ID,
    n_contribs: contribsn,
    n_repos: reposn
  })
  return insertOneResult.acknowledged
}

const getRepos = async client => {
  const repos = []
  for (const repo of [
    ['socketio', 'socket.io'],
    ['jquery', 'jquery'],
    ['transloadit', 'uppy'],
    ['sequelize', 'sequelize'],
    ['appium', 'appium'],
    ['puppeteer', 'puppeteer'],
    ['avajs', 'ava'],
    ['playcanvas', 'engine'],
    ['tutao', 'tutanota'],
    ['highlight', 'highlight'],
    ['supabase', 'supabase'],
    ['herbie-fp', 'herbie'],
    ['midwayjs', 'midway'],
    ['angular', 'angular'],
    ['expo', 'expo'],
    ['denoland', 'deno'],
    ['vitejs', 'vite'],
    ['discordjs', 'discord.js'],
    ['fabricjs', 'fabric.js'],
    ['vuejs', 'core'],
    ['markedjs', 'marked'],
    ['apify', 'crawlee'],
    ['electron', 'electron'],
    ['11ty', 'eleventy'],
    ['dagger', 'dagger'],
    ['cheeriojs', 'cheerio'],
    ['lint-staged', 'lint-staged'],
    ['metarhia', 'impress'],
    ['oven-sh', 'bun'],
    ['napi-rs', 'napi-rs'],
    ['pingcap', 'ossinsight'],
    ['sindresorhus', 'got'],
    ['twbs', 'bootstrap'],
    ['sveltejs', 'svelte'],
    ['prettier', 'prettier'],
    ['vercel', 'next.js'],
    ['biomejs', 'biome'],
    ['microsoft', 'vscode'],
    ['ether', 'etherpad-lite'],
    ['directus', 'directus'],
    ['spacedriveapp', 'spacedrive'],
    ['backstage', 'backstage'],
    ['sitespeedio', 'sitespeed.io'],
    ['webdriverio', 'webdriverio'],
    ['streetwriters', 'notesnook'],
    ['nomic-ai', 'gpt4all'],
    ['salesforce', 'lwc'],
    ['botpress', 'botpress']
  ]) {
    console.debug('repo: %s/%s...', repo.at(0), repo.at(1))
    const repository = await client.rest.repos.get({
      owner: repo.at(0),
      repo: repo.at(1)
    })
    repos.push(repository.data)
  }

  return repos
}

const githubAccessTok = process.env.GITHUB_ACCESS_TOKEN_CONTRIBS
if (!githubAccessTok) {
  throw new Error('missing Github access token')
}
const githubClient = new Octokit({ auth: githubAccessTok })

const repos = await getRepos(githubClient)
console.debug('repos: %d', repos.length)

let contribsn = 0
for (const repo of repos) {
  const repoOwner = repo.owner.login
  const repoName = repo.name
  console.debug('repo: %s/%s', repoOwner, repoName)

  const TMP_DIR = path.join('/tmp', `${repoOwner}_${repoName}`)

  console.debug('cloning: %s', repo.clone_url)
  execSync(`git clone -q --depth 1 --no-tags --filter=blob:limit=100k ${repo.clone_url} ${TMP_DIR}`)

  // Remove extraneous
  rmSync(`${TMP_DIR}/.git`, { force: true, maxRetries: 1, recursive: true })
  execSync(`find ${TMP_DIR}/ -name 'node_modules' -type d -prune -exec rm -rf '{}' +`)

  const contribs = []
  let locusn = 0
  let filesn = 0
  for (let file of await findNodeJsFiles(TMP_DIR, [])) {
    console.debug('file: %s', file)
    filesn += 1

    const buffer = readFileSync(file)
    const code = buffer.toString()
    const locus = extract(code)

    if (locus.length === 0) {
      continue
    }
    locusn += locus.length
    console.debug('locus: %d', locus.length)

    file = file.split(TMP_DIR).pop()
    const filepath = path.dirname(file)
    const filename = path.basename(file)

    contribs.push({
      locus,
      code,
      filepath,
      filename,
      repo_name: repoName,
      repo_owner: repoOwner
    })

    contribsn += 1
  }

  console.debug('contribs: %d', contribs.length)
  console.debug('locus: %d', locusn)
  console.debug('files: %d', filesn)

  rmSync(path.join('/tmp', `${repoOwner}_${repoName}`), { force: true, recursive: true })

  if (contribs.length === 0) {
    continue
  }

  // Delete existing and insert new contributions
  await mongoCollection.deleteMany({
    repo_owner: repoOwner,
    repo_name: repoName
  })
  await insertContribs(contribs)
}

console.debug('contribs: %d', contribsn)

const licencesSaved = await insertLicenses()
console.debug('licenses: %s', licencesSaved)

await insertCatalogue(contribsn, repos.length)

process.exit(0)
