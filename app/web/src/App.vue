<script setup>
import { provide } from 'vue'
import { useRoute } from 'vue-router'
import { useTheme } from 'vuetify/lib/framework.mjs'

const theme = useTheme()
const route = useRoute()

const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}

provide('setDocumentTitle', title => {
  document.title = `stdlibs.com - ${title}`
})
provide('setCanonicalUrl', url => {
  const canonical = document.querySelector('link[rel="canonical"]')
  canonical.href = `https://www.stdlibs.com${url}`
})
</script>

<template>
  <v-app id="inspire">
    <v-app-bar
      :color="route.path === '/' ? 'transparent' : ''"
      absolute
      flat>
      <v-container
        fluid
        class="mx-auto d-flex align-center justify-center">
        <v-btn
          v-if="route.path !== '/'"
          size="small"
          variant="text"
          icon="mdi-home"
          to="/" />

        <v-spacer />

        <v-btn
          size="small"
          variant="text"
          icon="mdi-theme-light-dark"
          @click="toggleTheme()" />
      </v-container>
    </v-app-bar>

    <Suspense>
      <router-view />
    </Suspense>

    <v-footer
      :color="route.path === '/' ? 'transparent' : ''"
      app
      absolute>
      <v-container
        id="app"
        fluid
        class="mx-auto d-flex align-center justify-center">
        <v-btn
          to="/impressum"
          size="small"
          variant="text">
          Impressum
        </v-btn>
        <v-btn
          to="/privacy"
          size="small"
          variant="text">
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
