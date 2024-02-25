import { execSync } from 'node:child_process'
import { readFileSync, rmSync, readdirSync, statSync } from 'node:fs'
import path from 'node:path'
import { Octokit } from 'octokit'
import extract, { publicBuiltinModules } from './extract.mjs'
import mongoClient from './db.mjs'

// MongoDB Ids
const CAT_ID = '_cat'
const LICENSES_ID = "_licenses"
const TEMPORARY_DIRECTORY = path.join('/tmp', `contribs-node${Date.now()}`)

const mongoCollection = mongoClient.db('contribs').collection('node')

const findNodeJsFiles = async (directory, files) => {
  for (const filename of readdirSync(directory)) {
    try {
      const file = path.join(directory, filename)

      const statistics = statSync(file)
      if (statistics.isDirectory()) {
        files = await findNodeJsFiles(file, files)
        continue
      }

      switch (path.extname(filename)) {
        case '.js':
        case '.mjs':
        case '.cjs':
        case '.ts':
          files.push(file)
      }
    } catch (error) {
      // Symbolic links result into an error and can be ignored
      console.error(error)
    }
  }

  return files
}

const deleteTemporaryDirectory = () => {
  console.debug('deleting temp dir %s...', TEMPORARY_DIRECTORY)
  rmSync(TEMPORARY_DIRECTORY, { force: true, recursive: true })
}

const saveLicenses = async () => {
  await mongoCollection.deleteOne({ _id: LICENSES_ID })
  const insertOneResult = await mongoCollection.insertOne({
    _id: LICENSES_ID,
    repos: [
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
      }
    ]
  })
  return insertOneResult.acknowledged
}

const saveContributions = async contributions => {
  const insertManyResult = await mongoCollection.insertMany(contributions)
  return insertManyResult.insertedCount
}

const saveCatalogue = async (amountContributions, amountRepositories) => {
  await mongoCollection.deleteOne({ _id: CAT_ID })
  const insertOneResult = await mongoCollection.insertOne({
    _id: CAT_ID,
    n_contribs: amountContributions,
    n_repos: amountRepositories
  })
  return insertOneResult.acknowledged
}

const getHandpickedRepositories = async githubClient => {
  const repositories = []
  for (const handpickedRepository of [
    ['socketio', 'socket.io'],
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
    ['axios', 'axios'],
    ['prettier', 'prettier'],
    ['vercel', 'next.js'],
    ['biomejs', 'biome']
  ]) {
    console.debug('fetching repo %s/%s...', handpickedRepository.at(0), handpickedRepository.at(1))
    const repository = await githubClient.rest.repos.get({
      owner: handpickedRepository.at(0),
      repo: handpickedRepository.at(1)
    })
    repositories.push(repository.data)
  }

  return repositories
}

const cleanRepository = async (repositoryOwner, repositoryName) => {
  try {
    deleteTemporaryDirectory()

    await mongoCollection.deleteMany({
      repo_owner: repositoryOwner,
      repo_name: repositoryName,
    })
  } catch (error) {
    console.error(error)
    process.exit(1)
  }
}

try {
  const githubAccessToken = process.env.GITHUB_ACCESS_TOKEN_CONTRIBS
  if (!githubAccessToken) {
    throw new Error('missing Github access token')
  }
  const githubClient = new Octokit({ auth: githubAccessToken })

  const handpickedRepositories = await getHandpickedRepositories(githubClient)
  console.debug('found %d repos', handpickedRepositories.length)

  let amountContributions = 0
  for (const repository of handpickedRepositories) {
    const repositoryOwner = repository.owner.login
    const repositoryName = repository.name
    console.debug('cleaning...')
    await cleanRepository(repositoryOwner, repositoryName)

    console.debug('cloning repo %s to %s...', repository.clone_url, TEMPORARY_DIRECTORY)
    execSync(`git clone -q --depth 1 --no-tags --filter=blob:limit=100k ${repository.clone_url} ${TEMPORARY_DIRECTORY}`)

    console.debug('cleaning repo files...')
    rmSync(`${TEMPORARY_DIRECTORY}/.git`, { force: true, maxRetries: 1, recursive: true })
    execSync(`find ${TEMPORARY_DIRECTORY}/ -name 'node_modules' -type d -prune -exec rm -rf '{}' +`)

    console.debug('searching Node.js files...')
    const contributions = []
    for (let file of await findNodeJsFiles(TEMPORARY_DIRECTORY, [])) {
      const buffer = readFileSync(file)
      const code = buffer.toString()
      const apis = extract(code)

      if (apis.length === 0) {
        continue
      }

      file = file.split(TEMPORARY_DIRECTORY).pop()
      const filepath = path.dirname(file)
      const filename = path.basename(file)

      contributions.push({
        apis,
        code,
        filepath,
        filename,
        repo_name: repositoryName,
        repo_owner: repositoryOwner
      })

      amountContributions += 1
    }

    console.debug('found %d contributions', contributions.length)
    const contributionsSaved = await saveContributions(contributions)
    console.debug('%d contributions saved', contributionsSaved)
  }

  const licencesSaved = await saveLicenses()
  console.debug('licenses saved: %s', licencesSaved)

  const catalogueSaved = await saveCatalogue(amountContributions, handpickedRepositories.length)
  console.debug('catalogue saved: %s', catalogueSaved)

  process.exit(0)
} catch (error) {
  console.error(error)
  process.exit(1)
}
