<!-- eslint-disable vue/no-v-html -->
<template>
  <v-main>
    <v-container>
      <v-row>
        <v-col
          cols="12"
          xs="12"
          sm="12"
          md="4"
          lg="2"
          xl="2">
          <!-- Technology -->
          <v-card
            flat>
            <v-card-item>
              <template #prepend>
                <v-icon
                  v-if="technology === 'go'"
                  size="x-large"
                  icon="mdi-language-go"
                  color="dark" />
                <v-icon
                  v-if="technology === 'node'"
                  size="x-large"
                  icon="mdi-nodejs"
                  color="dark" />
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
            flat
            :loading="isLoadingNamespaces"
            class="mt-4">
            <v-card-title class="text-caption">
              Namespaces ({{ filteredNamespaces.length }})
            </v-card-title>
            <v-card-text>
              <v-text-field
                v-model="namespaceQuery"
                prepend-inner-icon="mdi-magnify"
                density="compact"
                bg-color="transparent"
                variant="plain"
                label="Search namespaces" />
              <v-list
                ref="namespacesHtmlElement"
                v-model:selected="selectedNamespace"
                :disabled="isLoadingApis || isLoadingNamespaces || isLoadingContributions"
                density="compact"
                return-object
                nav
                max-height="300">
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
            flat
            :loading="isLoadingApis"
            class="mt-4">
            <v-card-title class="text-caption">
              APIs ({{ filteredApis.length }})
            </v-card-title>
            <v-card-text>
              <v-text-field
                v-model="apisQuery"
                prepend-inner-icon="mdi-magnify"
                density="compact"
                label="Search APIs"
                variant="plain" />
              <v-list
                ref="apisHtmlElement"
                v-model:selected="selectedApi"
                :disabled="isLoadingApis || isLoadingNamespaces || isLoadingContributions"
                bg-color="transparent"
                density="compact"
                return-object
                nav
                max-height="300">
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
          xs="12"
          sm="12"
          md="8"
          lg="6"
          xl="6">
          <!-- Results information -->
          <v-card
            v-if="selectedNamespace.length > 0 && catalogue.vids[selectedNamespace.at(0)] != undefined"
            flat>
            <v-card-text>
              <iframe
                width="100%"
                height="400"
                :src="`https://www.youtube.com/embed/${catalogue.vids[selectedNamespace.at(0)]}?amp;controls=0`"
                frameborder="0"
                allow="autoplay; encrypted-media; gyroscope; picture-in-picture; web-share"
                allowfullscreen />
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
                size="x-small"
                icon="mdi-api" /> <a
                  class="text-medium-emphasis"
                  target="_blank"
                  :href="`https://pkg.go.dev//${selectedNamespace.at(0)}#${selectedApi.at(0)}`">
                  Go doc</a>
            </v-card-subtitle>
            <v-card-subtitle
              v-show="technology === 'node' && technology !== 'node'">
              <v-icon
                size="x-small"
                icon="mdi-api" /> <a
                  class="text-medium-emphasis"
                  target="_blank"
                  :href="`https://nodejs.org/api/${selectedNamespace.at(0).replace('node:', '')}.html`">
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

          <!-- Results -->
          <div
            v-for="(contribution, index) in contributions"
            :key="contribution._id">
            <div class="text-right text-caption mb-2 mt-4">
              #{{ (index+1)+pagination.perPage*(pagination.page-1) }}
            </div>
            <XCodeViewer
              :lines="findLines(contribution.apis)"
              :language="technologyToLanguageMapper.get(technology)"
              :contribution="contribution" />
          </div>

          <v-pagination
            v-if="contributions.length > 0"
            v-model="pagination.page"
            class="mt-4"
            density="comfortable"
            size="small"
            rounded="circle"
            :disabled="isLoadingContributions || isLoadingApis || isLoadingNamespaces"
            :total-visible="3"
            :length="Math.ceil(pagination.total / pagination.perPage) < 0 ? 1 : Math.ceil(pagination.total / pagination.perPage)" />
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script setup>
import XCodeViewer from '../components/XCodeViewer.vue'
import { getCatalogue, getApis, getContributions, getLicenses } from '../api.js'
import { nextTick, onMounted, ref, watch, computed, onBeforeMount, inject, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, find, reduce } from 'lodash'

const route = useRoute()
const router = useRouter()

const technology = route.meta.technology

const technologyToLanguageMapper = new Map()
technologyToLanguageMapper.set('go', 'go')
technologyToLanguageMapper.set('node', 'javascript')

const apis = ref([])
const apisQuery = ref('')
const apisHtmlElement = ref(null)
const catalogue = ref({})
const contributions = ref([])
const filteredApis = ref([])
const filteredNamespaces = ref([])
const isLoadingContributions = ref(false)
const isLoadingApis = ref(false)
const isLoadingNamespaces = ref(false)
const namespaces = ref([])
const namespacesHtmlElement = ref(null)
const namespaceQuery = ref('')
const pagination = ref({
  page: 1,
  perPage: 6,
  total: 0
})
const selectedApi = ref([])
const selectedApiDocumentation = ref(null)
const selectedNamespace = ref([])

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
const findLines = apis => {
  return reduce(apis, (apis, api) => {
    if (api.ident === `${selectedNamespace.value}.${selectedApi.value.at(0)}`) {
      apis.push(api.line)
    }
    return apis
  }, []).sort()
}
const resetPagination = () => {
  pagination.value.page = 1
  pagination.value.perPage = 6
  pagination.value.total = 0
}
const toggleIsLoadingApis = () => {
  isLoadingApis.value = !isLoadingApis.value
}
const toggleIsLoadingContributions = () => {
  isLoadingContributions.value = !isLoadingContributions.value
}
const toggleIsLoadingNamespaces = () => {
  isLoadingNamespaces.value = !isLoadingNamespaces.value
}

const setDocumentTitle = inject('setDocumentTitle')

const licenses = await getLicenses(technology)
provide('licenses', licenses)

onBeforeMount(async () => {
  setDocumentTitle(technology)
})

onMounted(async () => {
  toggleIsLoadingNamespaces()
  toggleIsLoadingApis()

  catalogue.value = await getCatalogue(technology)
  namespaces.value = catalogue.value.ns.sort()
  filteredNamespaces.value = namespaces.value

  // Resolve existing namespace query request
  if (route.query.ns) {
    // Take namespace from parameters
    selectedNamespace.value = [decodeURIComponent(route.query.ns)]
  } else {
    // Preselect first namespace
    namespacesHtmlElement.value.select(namespaces.value.at(0), namespaces.value.at(0))
    router.replace({
      query: {
        ...route.query,
        ns: encodeURIComponent(namespaces.value.at(0))
      }
    })
    setDocumentTitle(`${technology}/${namespaces.value.at(0)}`)
  }

  toggleIsLoadingNamespaces()

  // Apis
  apis.value = await getApis(technology, selectedNamespace.value)
  filteredApis.value = apis.value

  watch(selectedNamespace, async () => {
    toggleIsLoadingApis()

    selectedApi.value = []
    contributions.value = []
    if (selectedNamespace.value.length === 0) {
      toggleIsLoadingApis()
      return
    }

    router.replace({
      query: {
        ...route.query,
        api: undefined,
        page: undefined,
        ns: encodeURIComponent(selectedNamespace.value)
      }
    })
    setDocumentTitle(`${technology}/${selectedNamespace.value}`)

    resetPagination()
    apis.value = await getApis(technology, selectedNamespace.value)
    filteredApis.value = apis.value
    if (apis.value.length === 0) {
      throw new Error('can\'t find apis')
    }

    toggleIsLoadingApis()

    document.getElementById('apis').scrollIntoView()
  }, { deep: true })

  watch(computed(() => pagination.value.page), async () => {
    contributions.value = []

    router.replace({
      query: {
        ...route.query,
        page: pagination.value.page
      }
    })

    setDocumentTitle(`${technology}/${selectedNamespace.value}/${selectedApi.value} #${pagination.value.page}`)

    await getContributionsHandler()
  })

  // Resolve existing query request
  if (route.query.api) {
    // Take namespace from Url parameters
    const api = route.query.api
    selectedApi.value = [decodeURIComponent(api)]
    selectedApiDocumentation.value = get(find(apis.value, ['name', api]), 'doc')

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

  toggleIsLoadingApis()

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

  watch(namespaceQuery, () => {
    if (namespaceQuery.value === '') {
      filteredNamespaces.value = namespaces.value
      return
    }

    filteredNamespaces.value = namespaces.value.filter(namespace => {
      return namespace.startsWith(namespaceQuery.value) || namespace === selectedNamespace.value.at(0)
    })
  })
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
