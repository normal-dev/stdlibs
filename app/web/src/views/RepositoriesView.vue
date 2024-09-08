<script setup>
import {
  inject
} from 'vue'
import {
  getRepositories
} from '../api'
import {
  format
} from 'date-fns'

const setDocumentTitle = inject('setDocumentTitle')

setDocumentTitle('Repositories')

const repositories = await getRepositories()
</script>

<template>
  <v-main>
    <v-container>
      <v-row>
        <v-col
          v-for="(repository, index) in repositories"
          :key="index"
          cols="12"
          lg="3"
          md="4"
          sm="12">
          <v-card flat>
            <v-card-item>
              <v-card-title>{{ repository.repo_owner }}/{{ repository.repo_name }}</v-card-title>
              <v-card-subtitle>Updated at {{ format(repository.updated, 'dd/MM/yyyy') }}</v-card-subtitle>
            </v-card-item>

            <v-card-text>
              {{ repository.locusn }} contributions
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
