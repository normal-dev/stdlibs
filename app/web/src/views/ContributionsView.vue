<script setup>
import XCodeViewer from '../components/XCodeViewer.vue'
import { getCatalogue, getApis, getContributions, getLicenses } from '../api.js'
import { nextTick, onMounted, ref, watch, computed, inject, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, find, reduce } from 'lodash'

const route = useRoute()
const router = useRouter()

const setDocumentTitle = inject('setDocumentTitle')

// Node.js, Go, etc.
const technology = route.meta.technology

// Set initial document title
setDocumentTitle(technology)

// Maps technology to language
const techMapper = new Map()
techMapper.set('go', 'go')
techMapper.set('node', 'javascript')

// Namespaces
const namespaces = ref([])
const namespacesHtmlElement = ref(null)
const namespaceQuery = ref('')
const filteredNamespaces = ref([])
const isLoadingNamespaces = ref(false)
const selectedNamespace = ref([])
const toggleIsLoadingNamespaces = () => {
  isLoadingNamespaces.value = !isLoadingNamespaces.value
}

// Apis
const apis = ref([])
const apisQuery = ref('')
const apisHtmlElement = ref(null)
const filteredApis = ref([])
const isLoadingApis = ref(false)
const selectedApi = ref([])
const selectedApiDocumentation = ref(null)
const toggleIsLoadingApis = () => {
  isLoadingApis.value = !isLoadingApis.value
}

// Contributions
const catalogue = ref({})
const contributions = ref([])
const isLoadingContributions = ref(false)
const toggleIsLoadingContributions = () => {
  isLoadingContributions.value = !isLoadingContributions.value
}
const getContributionsHandler = async () => {
  toggleIsLoadingContributions()

  const { per_page, total, contribs } = await getContributions(
    technology,
    selectedNamespace.value,
    selectedApi.value,
    pagination.value.page
  )
  pagination.value.perPage = per_page
  pagination.value.total = total
  contributions.value = contribs

  toggleIsLoadingContributions()

  await nextTick()
  if (total === 0) {
    return
  }
  // Scroll to contributions
  document.getElementById('contributions').scrollIntoView()
}
const findLines = locus => {
  return reduce(locus, (locus, { ident, line }) => {
    if (ident === `${selectedNamespace.value}.${selectedApi.value.at(0)}`) {
      locus.push(line)
    }
    return locus
  }, []).sort()
}

// Pagination
const pagination = ref({
  page: 1,
  perPage: 6,
  total: 0
})
const resetPagination = () => {
  pagination.value.page = 1
  pagination.value.perPage = 6
  pagination.value.total = 0
}

// Get and provide licenses for code viewer
const licenses = await getLicenses(technology)
provide('licenses', licenses)

onMounted(async () => {
  toggleIsLoadingNamespaces()
  toggleIsLoadingApis()

  // Fetch catalogue with included namespaces
  catalogue.value = await getCatalogue(technology)
  namespaces.value = catalogue.value.ns.sort()
  filteredNamespaces.value = namespaces.value

  // Select potential namespace
  if (route.query.ns) {
    // Take namespace from parameters
    selectedNamespace.value = [decodeURIComponent(route.query.ns)]
  } else {
    // No namespace given, preselect first one
    namespacesHtmlElement.value.select(namespaces.value.at(0), namespaces.value.at(0))
    router.replace({
      query: {
        ...route.query,
        ns: encodeURIComponent(namespaces.value.at(0))
      }
    })

    setDocumentTitle(`${technology}/${namespaces.value.at(0)}`)
  }

  // Resolving namespaces done
  toggleIsLoadingNamespaces()

  // Fetch apis based on namespace
  apis.value = await getApis(technology, selectedNamespace.value)
  filteredApis.value = apis.value

  // Start namespace search watcher
  watch(selectedNamespace, async () => {
    toggleIsLoadingApis()

    selectedApi.value = []
    contributions.value = []
    if (selectedNamespace.value.length === 0) {
      // Cancel
      toggleIsLoadingApis()
      return
    }

    // Namespace selected, reset api and page
    router.replace({
      query: {
        ...route.query,
        api: undefined,
        page: undefined,
        ns: encodeURIComponent(selectedNamespace.value)
      }
    })
    resetPagination()
    setDocumentTitle(`${technology}/${selectedNamespace.value}`)

    // Fetch apis with given namespace
    apis.value = await getApis(technology, selectedNamespace.value)
    filteredApis.value = apis.value
    if (apis.value.length === 0) {
      throw new Error('can\'t find apis')
    }

    // Resolving apis done
    toggleIsLoadingApis()

    document.getElementById('apis').scrollIntoView()
  }, { deep: true })

  // Load contributions based on pagination
  watch(computed(() => pagination.value.page), async () => {
    contributions.value = []

    router.replace({
      query: {
        ...route.query,
        page: pagination.value.page
      }
    })

    await getContributionsHandler()
    setDocumentTitle(`${technology}/${selectedNamespace.value}/${selectedApi.value} #${pagination.value.page}`)
  })

  // Select potential api
  if (route.query.api) {
    // Take namespace from Url parameters
    const api = route.query.api
    selectedApi.value = [decodeURIComponent(api)]
    selectedApiDocumentation.value = get(find(apis.value, ['name', api]), 'doc')

    // Jump to potential page
    if (route.query.page) {
      pagination.value.page = Number(route.query.page)
      await nextTick()
    }

    await getContributionsHandler()

    router.replace({
      query: {
        ...route.query,
        api: selectedApi.value
      }
    })
    setDocumentTitle(`${technology}/${selectedNamespace.value}/${selectedApi.value}`)
  }

  // Resolving apis done
  toggleIsLoadingApis()

  // Start selected api watcher
  watch(selectedApi, async () => {
    contributions.value = []
    if (selectedApi.value.length === 0) {
      return
    }

    selectedApiDocumentation.value = get(find(apis.value, ['name', selectedApi.value.at(0)]), 'doc')

    resetPagination()
    await getContributionsHandler()

    router.replace({
      query: {
        ...route.query,
        api: selectedApi.value.at(0),
        page: undefined
      }
    })
    setDocumentTitle(`${technology}/${selectedNamespace.value}/${selectedApi.value}`)
  }, { deep: true })

  // Query for namespace
  watch(namespaceQuery, () => {
    if (namespaceQuery.value === '') {
      filteredNamespaces.value = namespaces.value
      return
    }

    filteredNamespaces.value = namespaces.value.filter(namespace => {
      return namespace.startsWith(namespaceQuery.value) || namespace === selectedNamespace.value.at(0)
    })
  })
  // Query for api
  watch(apisQuery, () => {
    if (apisQuery.value === '') {
      filteredApis.value = apis.value
      return
    }

    filteredApis.value = apis.value.filter(api => {
      return api.name.startsWith(apisQuery.value) || api.name === selectedApi.value.at(0)
    })
  })
})
</script>

<template>
  <v-main>
    <v-container>
      <v-row>
        <v-col
          cols="12"
          lg="2"
          md="4"
          sm="12"
          xl="2"
          xs="12">
          <!-- Technology -->
          <v-card
            flat>
            <v-card-item>
              <template #prepend>
                <v-icon
                  v-if="technology === 'go'"
                  color="dark"
                  icon="mdi-language-go"
                  size="x-large" />
                <v-icon
                  v-if="technology === 'node'"
                  color="dark"
                  icon="mdi-nodejs"
                  size="x-large" />
              </template>
              <v-card-title v-if="technology === 'go'">
                Go ({{ catalogue.version }})
              </v-card-title>
              <v-card-title v-if="technology === 'node'">
                Node.js ({{ catalogue.version }})
              </v-card-title>
              <v-card-subtitle v-if="technology === 'go'">
                by The Go Authors
              </v-card-subtitle>
              <v-card-subtitle v-if="technology === 'node'">
                by OpenJS Foundation
              </v-card-subtitle>
            </v-card-item>

            <v-card-text v-if="technology === 'go'">
              Go is a statically typed, compiled high-level programming
              language designed at Google by Robert Griesemer, Rob Pike, and Ken
              Thompson.
            </v-card-text>
            <v-card-text v-if="technology === 'node'">
              Node.js is a cross-platform, open-source JavaScript runtime
              environment that can run on Windows, Linux, Unix, macOS, and more.
            </v-card-text>
          </v-card>

          <!-- Namespaces -->
          <v-card
            class="mt-4"
            flat
            :loading="isLoadingNamespaces">
            <v-card-title class="text-caption">
              Namespaces ({{ filteredNamespaces.length }})
            </v-card-title>
            <v-card-text>
              <v-text-field
                v-model="namespaceQuery"
                bg-color="transparent"
                density="compact"
                label="Search namespaces"
                prepend-inner-icon="mdi-magnify"
                variant="plain" />
              <v-list
                ref="namespacesHtmlElement"
                v-model:selected="selectedNamespace"
                density="compact"
                :disabled="isLoadingApis || isLoadingNamespaces || isLoadingContributions"
                max-height="300"
                nav
                return-object>
                <v-list-item
                  v-for="namespace in filteredNamespaces"
                  :key="namespace"
                  :disabled="namespace === selectedNamespace.at(0)"
                  :value="namespace">
                  {{ namespace }}
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>

          <!-- APIs -->
          <v-card
            id="apis"
            class="mt-4"
            flat
            :loading="isLoadingApis">
            <v-card-title class="text-caption">
              APIs ({{ filteredApis.length }})
            </v-card-title>
            <v-card-text>
              <v-text-field
                v-model="apisQuery"
                density="compact"
                label="Search APIs"
                prepend-inner-icon="mdi-magnify"
                variant="plain" />
              <v-list
                ref="apisHtmlElement"
                v-model:selected="selectedApi"
                bg-color="transparent"
                density="compact"
                :disabled="isLoadingApis || isLoadingNamespaces || isLoadingContributions"
                max-height="300"
                nav
                return-object>
                <v-list-item
                  v-for="api in filteredApis"
                  :key="api._id"
                  :disabled="api.name === selectedApi.at(0)"
                  :value="api.name">
                  {{ api.name }} <v-chip
                    class="ma-2"
                    size="x-small">
                    {{ api.type }}
                  </v-chip>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Contributions -->
        <v-col
          id="contributions"
          lg="6"
          md="8"
          sm="12"
          xl="6"
          xs="12">
          <!-- Video -->
          <v-card
            v-if="selectedNamespace.length > 0 && catalogue.vids[selectedNamespace.at(0)] != undefined"
            flat>
            <v-card-text>
              <iframe
                allow="autoplay; encrypted-media; gyroscope; picture-in-picture; web-share"
                allowfullscreen
                frameborder="0"
                height="400"
                :src="`https://www.youtube.com/embed/${catalogue.vids[selectedNamespace.at(0)]}?amp;controls=0`"
                width="100%" />
            </v-card-text>
          </v-card>

          <v-card
            v-if="selectedApi.length > 0"
            :class="{ 'mt-4': catalogue.vids[selectedNamespace.at(0)] != undefined }"
            flat>
            <v-card-title>
              {{ `${selectedNamespace.at(0)}.${selectedApi.at(0)}` }} ({{ pagination.total }})
            </v-card-title>
            <v-card-subtitle
              v-show="technology === 'go'">
              <v-icon
                icon="mdi-api"
                size="x-small" /> <a
                  class="text-medium-emphasis"
                  :href="`https://pkg.go.dev//${selectedNamespace.at(0)}#${selectedApi.at(0)}`"
                  target="_blank">
                  Go doc</a>
            </v-card-subtitle>
            <v-card-subtitle
              v-show="technology === 'node' && technology !== 'node'">
              <v-icon
                icon="mdi-api"
                size="x-small" /> <a
                  class="text-medium-emphasis"
                  :href="`https://nodejs.org/api/${selectedNamespace.at(0).replace('node:', '')}.html`"
                  target="_blank">
                  Node.js documentation</a>
            </v-card-subtitle>
            <v-card-text
              v-if="technology !== 'go' && technology !== 'node'">
              <p v-html="selectedApiDocumentation" />
            </v-card-text>
            <br>
          </v-card>

          <v-skeleton-loader
            v-if="isLoadingContributions"
            class="mt-4"
            type="paragraph" />

          <!-- Contributions code -->
          <div
            v-for="(contribution, index) in contributions"
            :key="contribution._id">
            <div class="text-right text-caption mb-2 mt-4">
              #{{ (index+1)+pagination.perPage*(pagination.page-1) }}
            </div>
            <XCodeViewer
              :contribution="contribution"
              :language="techMapper.get(technology)"
              :lines="findLines(contribution.locus)" />
          </div>

          <v-pagination
            v-if="contributions.length > 0"
            v-model="pagination.page"
            class="mt-4"
            density="comfortable"
            :disabled="isLoadingContributions || isLoadingApis || isLoadingNamespaces"
            :length="Math.ceil(pagination.total / pagination.perPage) < 0 ? 1 : Math.ceil(pagination.total / pagination.perPage)"
            rounded="circle"
            size="small"
            :total-visible="3" />
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
