import { Page } from '@playwright/test'

const mockProgress = {
  files: 100,
  books: 50,
  time: '1234',
}

export const parserFixture = {
  async mockParse(page: Page, json: object = mockProgress, status = 200) {
    page.route('**/api/parse**', (route) => {
      const url = new URL(route.request().url())
      if (url.pathname === '/api/parse' || url.pathname === '/api/parse/rebuild') {
        route.fulfill({ json, status })
      } else {
        route.continue()
      }
    })
  },
}
