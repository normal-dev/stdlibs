<script setup>
const {
  technology,
  version
} = defineProps({
  technology: {
    type: String,
    required: true
  },

  version: {
    type: String,
    required: true
  }
})

let icon = null
let title = null
let subtitle = null
let text = null
switch (technology) {
  case 'go':
    icon = 'mdi-language-go'
    title = 'Go'
    subtitle = 'The Go Authors'
    text = 'Go is a statically typed, compiled high-level programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson.'
    break

  case 'node':
    icon = 'mdi-nodejs'
    title = 'Node.js'
    subtitle = 'OpenJS Foundation'
    text = 'Node.js is a cross-platform, open-source JavaScript runtime environment that can run on Windows, Linux, Unix, macOS, and more.'
    break

  case 'python':
    icon = 'mdi-language-python'
    title = 'Python'
    subtitle = 'Python Software Foundation'
    text = 'Python is a programming language that lets you work more quickly and integrate your systems more effectively.'
    break
}
</script>

<template>
  <v-card
    flat>
    <v-card-item>
      <template #prepend>
        <v-icon
          color="dark"
          :icon="icon"
          size="x-large" />
      </template>
      <v-card-title>
        {{ title }} ({{ version }})
      </v-card-title>
      <v-card-subtitle>
        by {{ subtitle }}
      </v-card-subtitle>
    </v-card-item>

    <v-card-text>
      {{ text }}
    </v-card-text>
  </v-card>
</template>
