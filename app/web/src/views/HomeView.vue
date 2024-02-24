<script setup>
import { useTheme } from 'vuetify/lib/framework.mjs'
import { getCatalogue, getLicenses, getRandomContributions } from '../api'
import { defineAsyncComponent, provide, computed, ref, onMounted } from 'vue'
import { isArray, mergeWith, reduce, sortBy } from 'lodash'

document.title = 'stdlibs.com'

const theme = useTheme()

const contributions = await getRandomContributions()
const contribution = contributions.at(0)
const cursor = ref(0)

const goLicenses = await getLicenses('go')
const nodeLicenses = await getLicenses('node')
provide('licenses', mergeWith(goLicenses, nodeLicenses, (source, target) => {
  if (isArray(source)) {
    return source.concat(target)
  }

  return source
}))

const goCatalogue = ref({})
const nodeCatalogue = ref({})

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
  }

  return ''
})
const contributionApis = computed(() => {
  return sortBy(contribution.apis, 'line')
})
const codeViewerLines = computed(() => {
  return reduce(sortBy(contribution.apis, 'line'), (apis, api) => {
    apis.push(api.line)
    return apis
  }, [])
})

const XCodeViewer = defineAsyncComponent(() => import('../components/XCodeViewer.vue'))

onMounted(async () => {
  goCatalogue.value = await getCatalogue('go')
  nodeCatalogue.value = await getCatalogue('node')
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
      <!-- Logo -->
      <v-row class="mx-auto text-center">
        <v-col>
          <v-card
            variant="text">
            <v-card-title class="text-h3">
              stdlibs.com
            </v-card-title>
            <v-card-subtitle class="font-weight-medium text-subtitle-1">
              Hand-picked API examples for your favorite technology including
              Go, Node.js, Python and more
            </v-card-subtitle>
          </v-card>
        </v-col>
      </v-row>

      <!-- Languages -->
      <v-row class="mt-2 pa-4">
        <!-- Go -->
        <v-col
          cols="12"
          lg="4">
          <v-card
            to="/go"
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  size="x-large"
                  icon="mdi-language-go"
                  color="dark" />
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
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-package-variant-closed" />{{ goCatalogue.n_ns }} packages
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-code-braces" />{{ goCatalogue.n_apis }} APIs
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-source-fork" />{{ goCatalogue.n_repos }} repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Node.js -->
        <v-col
          cols="12"
          lg="4">
          <v-card
            to="/node"
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  size="x-large"
                  icon="mdi-nodejs"
                  color="dark" />
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
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-package-variant-closed" />{{ nodeCatalogue.n_ns }} modules
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-code-braces" />{{ nodeCatalogue.n_apis }} APIs
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-source-fork" />{{ nodeCatalogue.n_repos }} repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Python -->
        <v-col
          cols="12"
          lg="4">
          <v-card
            disabled
            variant="text">
            <v-card-item>
              <template #prepend>
                <v-icon
                  size="x-large"
                  icon="mdi-language-python"
                  color="dark" />
              </template>
              <v-card-title>
                Python <v-chip size="x-small">
                  Coming soon
                </v-chip>
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
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-package-variant-closed" />0 modules
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2 mr-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-code-braces" />0 APIs
              </v-chip>
              <v-chip
                variant="text"
                label
                size="small"
                class="mt-2">
                <v-icon
                  size="small"
                  start
                  icon="mdi-source-fork" />0 repositories
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Texts -->
      <v-row class="mt-2 pa-4">
        <v-col
          xxl="6"
          xl="6"
          lg="6"
          md="12"
          sm="12">
          <p class="text-body-1">
            <span class="font-weight-medium">stdlibs.com</span> helps you to
            better understand the so-called standard library of your daily
            programming language or technology, or more precisely, its API.
            Standard libraries "<span class="font-italic">typically include
              definitions for commonly used algorithms, data structures, and
              mechanisms for input and output.</span>"
          </p>
        </v-col>
      </v-row>
      <v-row class="mt-2 pa-4">
        <v-col
          xxl="6"
          xl="6"
          lg="6"
          md="12"
          sm="12"
          offset-md="0"
          offset-sm="0"
          offset-lg="6">
          <p class="text-body-1 text-right">
            We conduct regularly static code analysis to compile fresh lists of API
            usages and their line numbers from industry-proven open-source
            repositories. You can navigate through these repository files with an
            specifically for this purpose configured code viewer, to examine how
            the API is consumed
          </p>
        </v-col>
      </v-row>
      <v-row class="mt-2 pa-4">
        <v-col
          xxl="6"
          xl="6"
          lg="6"
          md="12"
          sm="12">
          <p class="text-body-1">
            Using our services is more precise and predictable compared to other
            research methods: You spend less time on low-quality answers from
            the internet, finding an API is easy thanks to the automatically
            generated catalogues, we guarentee a high quality of open-source
            contributions
          </p>
        </v-col>
      </v-row>

      <!-- Contributions -->
      <v-row class="mt-2 pa-4">
        <v-col class="d-flex justify-center align-center">
          <v-list
            bg-color="transparent"
            lines="two">
            <v-list-subheader>APIs</v-list-subheader>
            <v-list-item
              v-for="(api, index) in contributionApis"
              :key="index"
              rounded
              :active="index === cursor"
              :title="api.ident"
              :subtitle="api.line" />
          </v-list>
        </v-col>
        <v-col class="d-flex justify-center align-center">
          <XCodeViewer
            width="800px"
            no-navigation
            slideshow
            variant="text"
            :language="codeViewerlanguage"
            :lines="codeViewerLines"
            :contribution="contribution"
            @next-slide="$event => cursor = $event" />
        </v-col>
      </v-row>
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
