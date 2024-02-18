<script setup>
import { watch, provide, ref } from 'vue'
import { useTheme } from 'vuetify/lib/framework.mjs'
import { throttle } from 'lodash'
import { search } from './api'
import { useRoute } from 'vue-router'

const theme = useTheme()
const route = useRoute()

// TEST START
const themeToggleIcon = ref('mdi-white-balance-sunny')
// TEST END
const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
}
// TEST START
watch(() => theme.global.current.value.dark, (dark) => {
  if (dark) {
    themeToggleIcon.value = 'mdi-white-balance-sunny'
  } else {
    themeToggleIcon.value = 'mdi-moon-waning-crescent'
  }
})
// TEST END

// TEST START
const query = ref('')
watch(query, throttle(async () => {
  const searchResults = await search(query)
  console.log(searchResults)
}, 1000, { leading: false }))
// TEST END

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
        <!-- <v-btn
          v-if="route.path !== '/'"
          size="small"
          variant="text"
          to="/about">
          About
        </v-btn> -->

        <v-spacer />

        <v-btn
          v-if="false"
          size="small"
          variant="text">
          Search
          <v-overlay
            activator="parent"
            location-strategy="connected"
            scroll-strategy="none">
            <v-card
              class="pa-2"
              min-width="200">
              <v-text-field
                v-model="query"
                label="Query" />
            </v-card>
          </v-overlay>
        </v-btn>
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
