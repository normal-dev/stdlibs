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
  return sortBy(contribution.locus, 'line')
})
const codeViewerLines = computed(() => {
  return reduce(sortBy(contribution.locus, 'line'), (locus, { line }) => {
    locus.push(line)
    return locus
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
      <v-row>
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
