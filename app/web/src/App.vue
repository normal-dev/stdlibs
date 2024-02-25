<script setup>
import { provide } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

provide('setDocumentTitle', title => {
  document.title = `stdlibs.com - ${title}`
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
          v-if="route.path === '/impressum' || route.path === '/privacy'"
          size="small"
          variant="text"
          icon="mdi-arrow-left"
          to="/" />
        <v-btn
          v-else
          size="small"
          variant="text"
          icon="mdi-home"
          to="/" />

        <v-spacer />

        <v-btn
          size="small"
          variant="text"
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
          variant="text">
          Impressum
        </v-btn>
        <v-btn
          color="dark"
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
