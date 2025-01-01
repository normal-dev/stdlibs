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

// Map technology
const techMapper = new Map()
techMapper.set('go', 'Go')
techMapper.set('node', 'Node.js')
techMapper.set('python', 'Python')
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
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-file-code"
                  size="small"
                  start /> {{ techMapper.get(repository.tech) }}
              </v-chip>
              <v-chip
                class="mt-2 mr-2"
                label
                size="small"
                variant="text">
                <v-icon
                  icon="mdi-xml"
                  size="small"
                  start /> {{ repository.locusn.toLocaleString(
                    undefined,
                    { minimumFractionDigits: 0 }
                  ) }} contributions
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>
