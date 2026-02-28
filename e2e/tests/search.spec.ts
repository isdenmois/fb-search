import { test, expect } from '@playwright/test'
import { HomePage } from '../pages/home.page'

const mockBooks = [
  {
    id: 1,
    lang: 'en',
    authors: 'John Doe',
    title: 'Test Book One',
    size: 1024,
    series: 'Test Series',
    serno: '1',
  },
  {
    id: 2,
    lang: 'en',
    authors: 'Jane Smith',
    title: 'Test Book Two',
    size: 2048,
    series: 'Test Series',
    serno: '2',
  },
]

test.describe('Book Search', () => {
  test('should search books and display results', async ({ page }) => {
    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ json: mockBooks })
    })

    const homePage = new HomePage(page)
    await homePage.goto()

    await expect(homePage.searchInput).toBeVisible()

    await homePage.search('anything')

    await expect(homePage.results.getByRole('link')).toHaveCount(2)
  })

  test('should display loading state while searching', async ({ page }) => {
    await page.route('**/api/search*', async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 100))
      await route.fulfill({ json: mockBooks })
    })

    const homePage = new HomePage(page)
    await homePage.goto()

    await homePage.search('test')

    await expect(homePage.results.getByRole('link')).toHaveCount(2)
  })

  test('should show empty results for no matches', async ({ page }) => {
    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ json: [] })
    })

    const homePage = new HomePage(page)
    await homePage.goto()

    await homePage.search('nonexistentbook12345')

    const bookItems = await homePage.getBookItems()
    expect(bookItems).toHaveLength(0)
  })
})
