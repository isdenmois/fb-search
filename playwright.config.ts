import { defineConfig } from '@playwright/test'

export default defineConfig({
  // Directory containing your test files
  testDir: './e2e/tests',

  // Maximum time one test can run for
  timeout: 30 * 1000,

  // Number of retries on failures
  retries: 2,

  // Use common settings for all tests
  use: {
    baseURL: 'http://localhost:5173',
    headless: true,
    viewport: { width: 1280, height: 720 },
    trace: 'on-first-retry',
    colorScheme: 'dark',
  },

  fullyParallel: true,

  // Start the Vite dev server before tests run
  // and also start the Go backend
  webServer: [
    {
      command: 'bun run web',
      port: 5173,
      timeout: 120 * 1000,
      reuseExistingServer: !process.env.CI,
    },
  ],
})
