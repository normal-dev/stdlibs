<!-- eslint-disable vue/no-mutating-props -->
<script setup>
import * as monaco from 'monaco-editor'
import {
  watch,
  ref,
  onMounted,
  inject
} from 'vue'
import {
  useTheme
} from 'vuetify/lib/framework.mjs'
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import tsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker?worker'
import {
  eq,
  every,
  find,
  gt,
  gte,
  isUndefined,
  nth,
  size,
  subtract
} from 'lodash'

const and = (...predicates) => every(predicates, Boolean)

const theme = useTheme()

// Slideshow event
const emit = defineEmits(['nextSlide'])

const {
  contribution,
  language,
  lines,
  slideshow,
  noNavigation
} = defineProps({
  contribution: {
    required: true,
    type: Object,
    default: () => ({
      repo_owner: '',
      repo_name: '',
      filepath: '',
      filename: '',
      code: ''
    })
  },

  // E. g.: "go", "javascript"
  language: String,

  // E. g.: "[3, 6, 99]"
  lines: {
    required: true,
    type: Array
  },

  // Part of a slideshow
  slideshow: Boolean,

  // Hide navigation buttons
  noNavigation: Boolean,

  // Card variant
  variant: String
})

// Code viewer HTML element
const codeViewer = ref(null)

// Licenses
const licenses = inject('licenses')
// Find license in licenses
const license = find(licenses.repos, ({ repo: repository }) => {
  return (and(
    eq(nth(repository), contribution.repo_owner),
    eq(nth(repository, 1), contribution.repo_name)
  ))
})

// Lines of code navigation
const cursor = ref(0)
let next = () => undefined
let previous = () => undefined

// Monaco code editor worker
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
  // Create Monaco editor instance
  const editor = monaco.editor.create(codeViewer.value, {
    automaticLayout: true,
    contextmenu: false,
    domReadOnly: true,
    fixedOverflowWidgets: true,
    language,
    minimap: { enabled: false },
    ovnamespaceerviewRulerLanes: 0,
    readOnly: true,
    scrollBeyondLastLine: false,
    theme: theme.global.current.value.dark ? 'vs-dark' : 'vs',
    value: contribution.code,
    wordWrap: 'on',
    wrappingStrategy: 'advanced'
  })

  // Highlight lines of code
  let decorationsCollection = editor.createDecorationsCollection([
    {
      range: new monaco.Range(
        nth(lines),
        1,
        nth(lines),
        1
      ),
      options: {
        isWholeLine: true,
        linesDecorationsClassName: 'code-viewer--line-decoration'
      }
    }
  ])
  // Show first line as default
  editor.revealLine(nth(lines, cursor.value))

  // Advance lines of code navigation
  next = () => {
    if (gt(
      cursor.value,
      subtract(size(lines), 1))
    ) {
      // No more lines left, cancel
      return
    }

    // Advance cursor
    cursor.value++

    // Reveal next line
    editor.revealLine(nth(lines, cursor.value))
    decorationsCollection.clear()
    decorationsCollection = editor.createDecorationsCollection([
      {
        range: new monaco.Range(
          nth(lines, cursor.value),
          1,
          nth(lines, cursor.value),
          1
        ),
        options: {
          isWholeLine: true,
          linesDecorationsClassName: 'code-viewer--line-decoration'
        }
      }
    ])

    // Propagate new cursor value
    emit('nextSlide', cursor.value)
  }

  // Go back to previous line of code
  previous = () => {
    if (eq(cursor.value, 0)) {
      // No more lines left, cancel
      return
    }

    // Go back to previous line
    cursor.value--

    // Reveal previous line
    editor.revealLine(nth(lines, cursor.value))
    decorationsCollection.clear()
    decorationsCollection = editor.createDecorationsCollection([
      {
        range: new monaco.Range(
          nth(lines, cursor.value),
          1,
          nth(lines, cursor.value),
          1
        ),
        options: {
          isWholeLine: true,
          linesDecorationsClassName: 'code-viewer--line-decoration'
        }
      }
    ])
  }

  // Set Monaco theme according to global theme
  watch(() => theme.global.current.value.dark, () => {
    monaco.editor.setTheme(theme.global.current.value.dark ? 'vs-dark' : 'vs')
  })

  // Auto-advance if slideshow is enabled
  if (slideshow) {
    setInterval(() => {
      if (gte(
        cursor.value,
        subtract(size(lines), 1)
      )) {
        cursor.value -= 1
      }

      next()
    }, 2875)
  }
})
</script>

<template>
  <v-card
    flat>
    <v-card-title>
      {{ contribution.filepath }}{{ eq(contribution.filepath, '/') ? '/' : '' }}{{ contribution.filename }} ({{ size(lines) }})
    </v-card-title>
    <v-card-subtitle>
      <v-icon
        icon="mdi-web"
        size="x-small" /> <a
          class="text-medium-emphasis"
          :href="`https://www.github.com/${contribution.repo_owner}/${contribution.repo_name}`"
          target="_blank">
          {{ contribution.repo_owner }}/{{ contribution.repo_name }}</a>
    </v-card-subtitle>
    <v-card-text>
      <div
        ref="codeViewer"
        class="code-viewer" />
      <div
        v-if="!isUndefined(license)"
        class="mt-4">
        <small>&copy; {{ license.author }}, {{ license.type }}</small>
      </div>
    </v-card-text>

    <v-card-actions v-if="!noNavigation">
      <v-btn
        class="flex-grow-1"
        color="dark"
        :disabled="eq(cursor, 0)"
        icon="mdi-chevron-left"
        @click="previous()" />
      <v-btn
        class="flex-grow-1"
        color="dark"
        :disabled="eq(
          cursor,
          subtract(size(lines), 1)
        )"
        icon="mdi-chevron-right"
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
