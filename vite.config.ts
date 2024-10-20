import { defineConfig, loadEnv } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
// import { analyzer } from 'vite-bundle-analyzer'
import UnoCSS from 'unocss/vite'
import { presetUno } from 'unocss'
import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  return {
    plugins: [
      svelte(),
      UnoCSS({
        presets: [presetUno({ preflight: false })],
      }),
      // analyzer(),
    ],
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:3000',
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
        pages: resolve(__dirname, 'web/pages'),
        features: resolve(__dirname, 'web/features'),
        entities: resolve(__dirname, 'web/entities'),
        shared: resolve(__dirname, 'web/shared'),
      },
    },
    test: {
      environment: 'happy-dom',
      isolate: false,
      fileParallelism: false,
      poolOptions: {
        forks: {
          isolate: false,
        },
      },
      setupFiles: ['@testing-library/svelte/vitest', 'vi-fetch/setup'],
      alias: {
        'svelte-routing': resolve('./src/shared/test/svelte-routing'),
      },
    },
  }
})
