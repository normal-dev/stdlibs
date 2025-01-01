<script setup>
import {
  useTheme
} from 'vuetify/lib/framework.mjs'
import {
  getCatalogue,
  getLicenses,
  getRandomContributions
} from '../api'
import {
  defineAsyncComponent,
  provide,
  computed,
  ref,
  onMounted
} from 'vue'
import {
  isArray,
  mergeWith,
  reduce,
  sortBy,
  nth,
  eq
} from 'lodash'

const theme = useTheme()

document.title = 'stdlibs.com'

// Contributions
const CodeViewerCard = defineAsyncComponent(() => import('../components/CodeViewerCard.vue'))
// Fetch random contribution
const contributions = await getRandomContributions()
const contribution = nth(contributions, 0)
const contributionLocus = computed(() => {
  return sortBy(contribution.locus, 'line')
})
const cursor = ref(0)
const codeViewerlanguage = computed(() => {
  const fileExtension = contribution.filename.split('.').pop()
  switch (fileExtension) {
    case 'mjs':
    case 'cjs':
    case 'js':
      return 'javascript'
    case 'ts':
      return 'typescript'
    case 'go':
      return 'go'
    case 'py':
      return 'python'
  }

  return ''
})
const codeViewerLines = computed(() => {
  return reduce(sortBy(contribution.locus, 'line'), (locus, { line }) => {
    locus.push(line)
    return locus
  }, [])
})

// Licenses
const goLicenses = await getLicenses('go')
const nodeLicenses = await getLicenses('node')
const pythonLicenses = await getLicenses('python')
// Merge and pass licenses to child components
provide('licenses', mergeWith(
  {},
  goLicenses,
  nodeLicenses,
  pythonLicenses,
  (source, target) => {
    if (isArray(source)) {
      return source.concat(target)
    }

    return source
  }))

// Catalogue
const goCatalogue = ref({})
const nodeCatalogue = ref({})
const pythonCatalogue = ref({})
onMounted(async () => {
  Promise
    .all([
      getCatalogue('go'),
      getCatalogue('node'),
      getCatalogue('python')
    ])
    .then(([go, node, python]) => {
      goCatalogue.value = go
      nodeCatalogue.value = node
      pythonCatalogue.value = python
    })
})
</script>

<template>
  <div
    id="stars-sm"
    :class="{
      invert: !theme.global.current.value.dark
    }" />
  <div
    id="stars-m"
    :class="{
      invert: !theme.global.current.value.dark
    }" />
  <div
    id="stars-lg"
    :class="{
      invert: !theme.global.current.value.dark
    }" />
  <v-main>
    <v-container>
      <v-row>
        <v-col align="center">
          <v-card
            style="max-width: 960px !important;"
            variant="text">
            <v-card-item>
              <v-card-title class="text-h4">
                stdlibs.com
              </v-card-title>
              <v-card-subtitle>
                Hand-picked API examples for your favorite technology including
                Go, Node.js, Python and more.
              </v-card-subtitle>
            </v-card-item>

            <v-card-text class="text-justify">
              stdlibs.com employs static code analysis to identify and extract
              real-world examples of standard library usage from curated
              open-source repositories. This resource aids developers in
              understanding how to effectively utilize specific functions within
              their chosen programming language or technology. By selecting a
              namespace and API, users can access relevant and prominent usage
              examples, supporting their development process.
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <!-- Go -->
        <v-col
          cols="12"
          lg="3"
          md="4">
          <v-card
            to="/go"
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  color="dark"
                  icon="mdi-language-go"
                  size="x-large" />
              </template>
              <v-card-title>
                Go
              </v-card-title>
              <v-card-subtitle>
                by The Go Authors
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              Go is a statically typed, compiled high-level programming
              language designed at Google by Robert Griesemer, Rob Pike, and Ken
              Thompson.

              <br>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-package-variant-closed"
                  size="small"
                  start />{{ goCatalogue.n_ns }} packages
              </v-chip>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-code-braces"
                  size="small"
                  start />{{ goCatalogue.n_apis }} APIs
              </v-chip>
              <v-chip
                class="mt-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-source-fork"
                  size="small"
                  start />{{ goCatalogue.n_repos }} repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Node.js -->
        <v-col
          cols="12"
          lg="3"
          md="4">
          <v-card
            to="/node"
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  color="dark"
                  icon="mdi-nodejs"
                  size="x-large" />
              </template>
              <v-card-title>
                Node.js
              </v-card-title>
              <v-card-subtitle>
                by OpenJS Foundation
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              Node.js is a cross-platform, open-source JavaScript runtime
              environment that can run on Windows, Linux, Unix, macOS, and more.

              <br>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-package-variant-closed"
                  size="small"
                  start />{{ nodeCatalogue.n_ns }} modules
              </v-chip>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-code-braces"
                  size="small"
                  start />{{ nodeCatalogue.n_apis }} APIs
              </v-chip>
              <v-chip
                class="mt-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-source-fork"
                  size="small"
                  start />{{ nodeCatalogue.n_repos }} repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Python -->
        <v-col
          cols="12"
          lg="3"
          md="4">
          <v-card
            to="/python"
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  color="dark"
                  icon="mdi-language-python"
                  size="x-large" />
              </template>
              <v-card-title>
                Python
              </v-card-title>
              <v-card-subtitle>
                by Python Software Foundation
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              Python is a high-level, general-purpose programming language. Its
              design philosophy emphasizes code readability with the use of
              significant indentation.

              <br>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-package-variant-closed"
                  size="small"
                  start />{{ pythonCatalogue.n_ns }} modules
              </v-chip>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-code-braces"
                  size="small"
                  start />{{ pythonCatalogue.n_apis }} APIs
              </v-chip>
              <v-chip
                class="mt-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-source-fork"
                  size="small"
                  start />{{ pythonCatalogue.n_repos }} repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Rust -->
        <v-col
          cols="12"
          lg="3"
          md="4">
          <v-card
            disabled
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  color="dark"
                  icon="mdi-language-rust"
                  size="x-large" />
              </template>
              <v-card-title>
                Rust <v-chip size="x-small">
                  Coming soon
                </v-chip>
              </v-card-title>
              <v-card-subtitle>
                by The Rust Team
              </v-card-subtitle>
            </v-card-item>

            <v-card-text>
              Hundreds of companies around the world are using Rust in
              production today for fast, low-resource, cross-platform solutions.
              Software you know and love, like Firefox, Dropbox, and Cloudflare,
              uses Rust.

              <br>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-package-variant-closed"
                  size="small"
                  start />0 packages
              </v-chip>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-code-braces"
                  size="small"
                  start />0 APIs
              </v-chip>
              <v-chip
                class="mt-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-source-fork"
                  size="small"
                  start />0 repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Contributions -->
      <v-lazy
        :min-height="200"
        :options="{ threshold: 0.25}"
        transition="fade-transition">
        <v-row class="mt-4">
          <v-col class="d-flex justify-center align-center">
            <v-list
              bg-color="transparent"
              lines="two">
              <v-list-subheader>APIs</v-list-subheader>
              <v-list-item
                v-for="(api, index) in contributionLocus"
                :key="index"
                :active="eq(index, cursor)"
                rounded
                :subtitle="api.line"
                :title="api.ident" />
            </v-list>
          </v-col>
          <v-col class="d-flex justify-center align-center">
            <CodeViewerCard
              :contribution="contribution"
              :language="codeViewerlanguage"
              :lines="codeViewerLines"
              no-navigation
              slideshow
              variant="text"
              width="800px"
              @next-slide="$event => cursor = $event" />
          </v-col>
        </v-row>
      </v-lazy>
    </v-container>
  </v-main>
</template>

<style scoped>
.invert {
  filter: invert(1);
}
</style>

<style lang="sass">
@function multiple-box-shadow ($n)
  $value: '#{random(3840)}px #{random(3840)}px #FFF'
  @for $i from 2 through $n
    $value: '#{$value}, #{random(3840)}px #{random(3840)}px #FFF'

  @return unquote($value)

$shadows-small:  multiple-box-shadow(200)
$shadows-medium: multiple-box-shadow(100)
$shadows-big:    multiple-box-shadow(25)

#stars-sm
  animation: translate-y-2048 50s linear infinite
  background: transparent
  box-shadow: $shadows-small
  height: 1px
  width: 1px

  &:after
    width: 1px
    height: 1px
      content: " "
      position: absolute
      background: transparent
      box-shadow: $shadows-small

#stars-m
  animation: translate-y-2048 100s linear infinite
  background: transparent
  box-shadow: $shadows-medium
  height: 2px
  width: 2px

  &:after
    background: transparent
    box-shadow: $shadows-medium
    content: " "
    height: 2px
    position: absolute
    width: 2px

#stars-lg
  animation: translate-y-2048 150s linear infinite
  background: transparent
  box-shadow: $shadows-big
  height: 3px
  width: 3px

  &:after
    background: transparent
    box-shadow: $shadows-big
    content: " "
    height: 3px
    position: absolute
    width: 3px

@keyframes translate-y-2048
  from
    transform: translateY(0px)
  to
    transform: translateY(-2048px)
</style>
