import {
  defineConfig
} from 'vite'
import {
  fileURLToPath,
  URL
} from 'url'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  base: '/',
  build: {
    outDir: '../website'
  },
  plugins: [vue()],
  publicPath: './',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
