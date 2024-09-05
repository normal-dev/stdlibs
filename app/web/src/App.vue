<script setup>
import { ref, provide, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useTheme } from 'vuetify/lib/framework.mjs'

const theme = useTheme()

const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

const themeToggleIcon = ref('mdi-invert-colors')

provide('setDocumentTitle', title => {
  document.title = `stdlibs.com - ${title}`
})

watch(() => theme.global.current.value.dark, (dark) => {
  if (dark) {
    themeToggleIcon.value = 'mdi-invert-colors'
  } else {
    themeToggleIcon.value = 'mdi-invert-colors-off'
  }
})
</script>

<template>
  <v-app id="inspire">
    <v-app-bar
      color="transparent"
      absolute
      flat>
      <v-container
        fluid
        class="mx-auto d-flex align-center justify-center">
        <v-btn
          size="small"
          variant="plain"
          icon="mdi-home-outline"
          to="/" />
        <v-btn
          size="small"
          variant="plain"
          icon="mdi-newspaper-variant-outline"
          to="/news" />

        <v-spacer />

        <v-btn
          size="small"
          variant="plain"
          :icon="themeToggleIcon"
          @click="toggleTheme()" />
      </v-container>
    </v-app-bar>

    <Suspense>
      <router-view />
    </Suspense>

    <v-footer
      color="transparent"
      app
      absolute>
      <v-container
        id="app"
        fluid>
        <v-btn
          color="dark"
          to="/impressum"
          size="small"
          variant="plain">
          Impressum
        </v-btn>
        <v-btn
          color="dark"
          to="/privacy"
          size="small"
          variant="plain">
          Privacy
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
