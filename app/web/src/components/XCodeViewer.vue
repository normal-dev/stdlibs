<!-- eslint-disable vue/no-mutating-props -->
<script setup>
import * as monaco from 'monaco-editor'
import { watch, ref, onMounted, inject } from 'vue'
import { useTheme } from 'vuetify/lib/framework.mjs'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import { find } from 'lodash'

const theme = useTheme()

const emit = defineEmits(['nextSlide'])
const props = defineProps({
  contribution: {
    required: true,
    type: Object
  },
  // go, javascript
  language: String,
  // [3, 6, 99]
  lines: {
    required: true,
    type: Array
  },
  slideshow: Boolean,
  noNavigation: Boolean,
  // Card variant
  variant: String
})

const licenses = inject('licenses')

const cursor = ref(0)
const codeViewerHtmlElement = ref(null)

const license = find(licenses.repos, ({ repo: repository }) => {
  const repositoryOwner = repository.at(0)
  const repositoryName = repository.at(1)
  return (repositoryOwner === props.contribution.repo_owner && repositoryName === props.contribution.repo_name)
})
let next = () => undefined
let previous = () => undefined

self.MonacoEnvironment = {
  getWorker (_, label) {
    switch (label) {
      case 'typescript':
      case 'javascript':
        return new tsWorker()

      default:
        return new editorWorker()
    }
  }
}

onMounted(async () => {
  const editor = monaco.editor.create(codeViewerHtmlElement.value, {
    automaticLayout: true,
    contextmenu: false,
    domReadOnly: true,
    fixedOverflowWidgets: true,
    language: props.language,
    minimap: { enabled: false },
    ovnamespaceerviewRulerLanes: 0,
    readOnly: true,
    scrollBeyondLastLine: false,
    theme: theme.global.current.value.dark ? 'vs-dark' : 'vs',
    value: props.contribution.code,
    wordWrap: 'on',
    wrappingStrategy: 'advanced'
  })

  let decorationsCollection = editor.createDecorationsCollection([
    {
      range: new monaco.Range(props.lines.at(0), 1, props.lines.at(0), 1),
      options: {
        isWholeLine: true,
        linesDecorationsClassName: 'code-viewer--line-decoration'
      }
    }
  ])
  // Show first line as default
  editor.revealLine(props.lines.at(cursor.value))

  next = () => {
    if (cursor.value > props.lines.length - 1) {
      return
    }

    cursor.value++

    editor.revealLine(props.lines.at(cursor.value))
    decorationsCollection.clear()
    decorationsCollection = editor.createDecorationsCollection([
      {
        range: new monaco.Range(props.lines.at(cursor.value), 1, props.lines.at(cursor.value), 1),
        options: {
          isWholeLine: true,
          linesDecorationsClassName: 'code-viewer--line-decoration'
        }
      }
    ])
    emit('nextSlide', cursor.value)
  }
  previous = () => {
    if (cursor.value === 0) {
      return
    }

    cursor.value--

    editor.revealLine(props.lines.at(cursor.value))
    decorationsCollection.clear()
    decorationsCollection = editor.createDecorationsCollection([
      {
        range: new monaco.Range(props.lines.at(cursor.value), 1, props.lines.at(cursor.value), 1),
        options: {
          isWholeLine: true,
          linesDecorationsClassName: 'code-viewer--line-decoration'
        }
      }
    ])
  }

  watch(() => theme.global.current.value.dark, () => {
    monaco.editor.setTheme(theme.global.current.value.dark ? 'vs-dark' : 'vs')
  })

  if (props.slideshow) {
    setInterval(() => {
      if (cursor.value >= props.lines.length - 1) {
        cursor.value = -1
      }

      next()
    }, 2875)
  }
})
</script>

<template>
  <v-card
    class="pa-1"
    color="transparent"
    :variant="props.variant ? props.variant : 'flat'">
    <v-card-title class="pl-4 pr-4 pt-4">
      {{ contribution.filepath }}{{ contribution.filepath !== '/' ? '/' : '' }}{{ contribution.filename }}
    </v-card-title>
    <v-card-subtitle class="pl-4 pr-4">
      <a
        class="text-medium-emphasis"
        target="_blank"
        :href="`https://www.github.com/${contribution.repo_owner}/${contribution.repo_name}`">
        {{ contribution.repo_owner }}/{{ contribution.repo_name }}
      </a> <v-icon
        size="x-small"
        icon="mdi-link" />
    </v-card-subtitle>
    <v-card-text
      class="pl-4 pr-4">
      <div
        ref="codeViewerHtmlElement"
        class="code-viewer mb-2" />

      <small>&copy; {{ license.author }}, {{ license.type }}</small>
    </v-card-text>

    <v-card-actions v-if="!props.noNavigation">
      <v-btn
        :disabled="cursor === 0"
        rounded="sm"
        icon="mdi-chevron-left"
        size="large"
        class="flex-grow-1"
        @click="previous()" />
      <v-btn
        :disabled="cursor === lines.length-1"
        rounded="sm"
        icon="mdi-chevron-right"
        size="large"
        class="flex-grow-1"
        @click="next()" />
    </v-card-actions>
  </v-card>
</template>

<style>
.code-viewer {
  height: 525px;
}

.code-viewer--line-decoration {
  background: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
  width: 5px !important;
  margin-left: 4px;
}
</style>
