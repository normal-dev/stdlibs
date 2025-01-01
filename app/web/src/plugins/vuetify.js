import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import {
  createVuetify
} from 'vuetify'
import {
  md3
} from 'vuetify/blueprints'

export default createVuetify({
  blueprint: md3,
  components: {
    ...components
  },
  directives,
  theme: {
    defaultTheme: 'dark'
  }
})
