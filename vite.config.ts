import { fileURLToPath } from 'node:url'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import { defineConfig } from 'vite'
import { configDefaults } from 'vitest/config'

// https://vitejs.dev/config/
export default defineConfig(() => {
  return {
    plugins: [vue(), UnoCSS()],
    server: {
      proxy: {
        '^/(api|dl)': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          secure: false,
        },
      },
    },
    root: 'web',
    build: {
      assetsInlineLimit: 64,
      outDir: '../public',
      emptyOutDir: true,
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./web', import.meta.url)),
      },
    },
    test: {
      environment: 'happy-dom',
      exclude: [...configDefaults.exclude, 'e2e/**'],
      root: fileURLToPath(new URL('./', import.meta.url)),
      setupFiles: ['./web/vitest.setup.ts'],
    },
  }
})
