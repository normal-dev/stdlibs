<script setup>
import {
  inject
} from 'vue'

const setDocumentTitle = inject('setDocumentTitle')
setDocumentTitle('404')
</script>

<template>
  <v-main>
    <v-container>
      <v-card
        flat
        variant="text">
        <v-card-title>
          404
        </v-card-title>
        <v-card-text>
          Sorry, we can't find this page.
        </v-card-text>
      </v-card>
    </v-container>
  </v-main>
</template>
