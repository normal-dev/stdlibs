let baseUrl = window.location.origin // https://www.stdlibs.com
if (process.env.NODE_ENV === 'development') {
  baseUrl = 'http://localhost:8080'
}

export const getRepositories = async () => {
  const url = `${baseUrl}/api/seo/repositories`
  const response = await fetch(url)
  return await response.json()
}

export const getRandomContributions = async () => {
  const url = `${baseUrl}/api/gen`
  const response = await fetch(url)
  return await response.json()
}

export const getCatalogue = async technology => {
  const url = `${baseUrl}/api/${technology}`
  const response = await fetch(url)
  return await response.json()
}

export const getLicenses = async technology => {
  const url = `${baseUrl}/api/${technology}/licenses`
  const response = await fetch(url)
  return await response.json()
}

export const getApis = async (technology, namespace) => {
  const url = `${baseUrl}/api/${technology}/${encodeURIComponent(namespace)}`
  const response = await fetch(url)
  return await response.json()
}

export const getContributions = async (technology, namespace, api, page) => {
  const url = `${baseUrl}/api/${technology}/${encodeURIComponent(namespace)}/${encodeURIComponent(api)}?page=${page}`
  const response = await fetch(url)
  return await response.json()
}
