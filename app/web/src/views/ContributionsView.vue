<script setup>
import CodeViewer from '../components/CodeViewerCard.vue'
import {
  getCatalogue,
  getApis,
  getContributions,
  getLicenses
} from '../api.js'
import {
  nextTick,
  onMounted,
  ref,
  watch,
  computed,
  inject,
  provide
} from 'vue'
import {
  useRoute,
  useRouter
} from 'vue-router'
import {
  get,
  find,
  reduce,
  eq,
  sortBy,
  isEmpty,
  filter,
  startsWith,
  some,
  size,
  replace,
  head
} from 'lodash'
import TechnologyCard from '../components/TechnologyCard.vue'

const or = (...predicates) => some(predicates, Boolean)

const route = useRoute()
const router = useRouter()

const setDocumentTitle = inject('setDocumentTitle')

// E. g.: "node", "python"
const technology = route.meta.technology

// Set initial document title
setDocumentTitle(technology)

// Maps technology to language
const techMapper = new Map()
techMapper.set('go', 'go')
techMapper.set('node', 'javascript')
techMapper.set('python', 'python')

// Namespaces
const namespaces = ref([])
const namespacesHtmlElement = ref(null)
const namespaceQuery = ref('')
const filteredNamespaces = ref([]) // TODO: Use "computed"
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

  const {
    per_page,
    total,
    contribs
  } = await getContributions(
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
  if (eq(total, 0)) {
    // No contributions found, cancel
    return
  }

  // Scroll to contributions
  document.getElementById('contributions').scrollIntoView()
}

// Returns the lines of code without identifier
const findLines = locus => {
  return sortBy(reduce(locus, (lines, { ident, line }) => {
    if (eq(
      ident,
      `${selectedNamespace.value}.${head(selectedApi.value)}`
    )) {
      lines.push(line)
    }

    return lines
  }, []))
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

// Fetch and provide licenses for code viewer
const licenses = await getLicenses(technology)
provide('licenses', licenses)

onMounted(async () => {
  toggleIsLoadingNamespaces()
  toggleIsLoadingApis()

  // Fetch catalogue with included namespaces
  catalogue.value = await getCatalogue(technology)
  namespaces.value = sortBy(catalogue.value.ns)
  filteredNamespaces.value = namespaces.value

  // Select potential namespace
  if (route.query.ns) {
    // Take namespace from parameters
    selectedNamespace.value = [decodeURIComponent(route.query.ns)]
  } else {
    // No namespace given, preselect first one
    namespacesHtmlElement.value.select(
      head(namespaces.value),
      head(namespaces.value)
    )
    // Update Url
    router.replace({
      query: {
        ...route.query,
        ns: encodeURIComponent(head(namespaces.value))
      }
    })

    setDocumentTitle(`${technology}/${head(namespaces.value)}`)
  }

  // Resolving namespaces done
  toggleIsLoadingNamespaces()

  // Fetch Apis based on namespace
  apis.value = await getApis(technology, selectedNamespace.value)
  filteredApis.value = apis.value

  // Start namespace search watcher
  watch(selectedNamespace, async () => {
    toggleIsLoadingApis()

    selectedApi.value = []
    contributions.value = []
    if (isEmpty(selectedNamespace.value)) {
      // Cancel
      toggleIsLoadingApis()
      return
    }

    // Namespace selected, reset Api and page
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

    // Fetch Apis with given namespace
    apis.value = await getApis(technology, selectedNamespace.value)
    filteredApis.value = apis.value
    if (isEmpty(apis.value)) {
      // No apis found, cancel
      toggleIsLoadingApis()
      return
    }

    // Resolving Apis done
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

  // Select Api based on Url parameter
  if (route.query.api) {
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

  // Resolving Apis done
  toggleIsLoadingApis()

  watch(selectedApi, async () => {
    contributions.value = []
    if (isEmpty(selectedApi.value)) {
      return
    }

    selectedApiDocumentation.value = get(
      find(apis.value, ['name', head(selectedApi.value)]),
      'doc'
    )

    resetPagination()
    await getContributionsHandler()

    router.replace({
      query: {
        ...route.query,
        api: head(selectedApi.value),
        page: undefined
      }
    })
    setDocumentTitle(`${technology}/${selectedNamespace.value}/${selectedApi.value}`)
  }, { deep: true })

  // Query for namespace
  watch(namespaceQuery, () => {
    if (isEmpty(namespaceQuery.value)) {
      filteredNamespaces.value = namespaces.value
      return
    }

    filteredNamespaces.value = filter(namespaces.value, namespace => {
      return or(
        startsWith(namespace, namespaceQuery.value),
        eq(namespace, head(selectedNamespace.value))
      )
    })
  })

  // Query for Api
  watch(apisQuery, () => {
    if (isEmpty(apisQuery.value)) {
      filteredApis.value = apis.value
      return
    }

    filteredApis.value = filter(apis.value, api => {
      return or(
        startsWith(api.name, apisQuery.value),
        eq(api.name, head(selectedApi.value))
      )
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
          <TechnologyCard
            v-if="!isEmpty(catalogue.version)"
            :technology="technology"
            :version="catalogue.version" />

          <!-- Namespaces -->
          <v-card
            class="mt-4"
            flat
            :loading="isLoadingNamespaces">
            <v-card-title class="text-caption">
              Namespaces ({{ size(filteredNamespaces) }})
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
                :disabled="or(
                  isLoadingApis,
                  isLoadingNamespaces,
                  isLoadingContributions
                )"
                max-height="300"
                nav
                return-object>
                <v-list-item
                  v-for="namespace in filteredNamespaces"
                  :key="namespace"
                  :disabled="eq(namespace, head(selectedNamespace))"
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
              APIs ({{ size(filteredApis) }})
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
                :disabled="or(
                  isLoadingApis,
                  isLoadingNamespaces,
                  isLoadingContributions
                )"
                max-height="300"
                nav
                return-object>
                <v-list-item
                  v-for="api in filteredApis"
                  :key="api._id"
                  :disabled="eq(api.name, head(selectedApi))"
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
          <!-- Documentation -->
          <v-card
            v-if="!isEmpty(selectedApi)"
            flat>
            <v-card-title>
              {{ head(selectedNamespace) }}
            </v-card-title>
            <v-card-subtitle>
              {{ head(selectedApi) }} ({{ pagination.total }})
            </v-card-subtitle>
            <v-card-text>
              <div v-show="eq(technology, 'go')">
                <v-icon
                  icon="mdi-api"
                  size="x-small" /> <a
                    class="text-medium-emphasis"
                    :href="`https://pkg.go.dev//${head(selectedNamespace)}#${head(selectedApi)}`"
                    target="_blank">
                    Go doc</a>
              </div>

              <div v-show="eq(technology, 'node')">
                <v-icon
                  icon="mdi-api"
                  size="x-small" /> <a
                    class="text-medium-emphasis"
                    :href="`https://nodejs.org/api/${replace(head(selectedNamespace), 'node:', '')}.html`"
                    target="_blank">
                    Node.js documentation</a>
              </div>

              <div v-show="eq(technology, 'python')">
                <v-icon
                  icon="mdi-api"
                  size="x-small" /> <a
                    class="text-medium-emphasis"
                    :href="`https://docs.python.org/3/library/${head(selectedNamespace)}.html`"
                    target="_blank">
                    Python documentation</a>
              </div>
            </v-card-text>
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
            <CodeViewer
              :contribution="contribution"
              :language="techMapper.get(technology)"
              :lines="findLines(contribution.locus)" />
          </div>

          <!-- Pagination -->
          <v-pagination
            v-if="!isEmpty(contributions)"
            v-model="pagination.page"
            class="mt-4"
            density="comfortable"
            :disabled="or(
              isLoadingContributions,
              isLoadingApis,
              isLoadingNamespaces
            )"
            :length="Math.ceil(pagination.total / pagination.perPage) < 0 ? 1 : Math.ceil(pagination.total / pagination.perPage)"
            rounded="circle"
            size="small"
            :total-visible="3" />
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
