<script setup>
import {
  ref,
  provide,
  watch
} from 'vue'
import {
  useTheme
} from 'vuetify/lib/framework.mjs'

const theme = useTheme()

// Toggle dark/white theme
const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}
// Change theme icon depending on theme
const themeToggleIcon = ref('mdi-invert-colors')
watch(() => theme.global.current.value.dark, dark => {
  if (dark) {
    themeToggleIcon.value = 'mdi-invert-colors'
    return
  }

  themeToggleIcon.value = 'mdi-invert-colors-off'
})

// Provide document title setter
provide('setDocumentTitle', title => {
  document.title = `stdlibs.com - ${title}`
})
</script>

<template>
  <v-app id="inspire">
    <v-app-bar
      absolute
      color="transparent"
      flat>
      <v-container
        class="mx-auto d-flex align-center justify-center"
        fluid>
        <v-btn
          icon="mdi-home-outline"
          size="small"
          to="/"
          variant="plain" />
        <v-btn
          icon="mdi-source-fork"
          size="small"
          to="/repositories"
          variant="plain" />

        <v-spacer />

        <v-btn
          :icon="themeToggleIcon"
          size="small"
          variant="plain"
          @click="toggleTheme()" />
      </v-container>
    </v-app-bar>

    <Suspense>
      <router-view />
    </Suspense>

    <v-footer
      absolute
      app
      color="transparent">
      <v-container
        id="app"
        fluid>
        <v-btn
          color="dark"
          href="https://github.com/normal-dev/stdlibs"
          size="small"
          target="_blank"
          variant="plain">
          Source code
        </v-btn>
        <v-btn
          color="dark"
          size="small"
          to="/impressum"
          variant="plain">
          Impressum
        </v-btn>
        <v-spacer />
      </v-container>
    </v-footer>
  </v-app>
</template>

<style>
@media only screen and (min-width: 960px) {
  .v-main .v-container {
    max-width: 1920px !important;
  }
}
</style>
