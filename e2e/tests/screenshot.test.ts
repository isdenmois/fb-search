import { test, expect } from '@playwright/test'
import { HomePage } from '../pages/home.page'

const mockBooks = [
  {
    id: '1',
    lang: 'en',
    authors: 'John Doe',
    title: 'Test Book One',
    size: 1024,
    series: 'Test Series',
    serno: '1',
  },
  {
    id: '2',
    lang: 'en',
    authors: 'Jane Smith',
    title: 'Test Book Two',
    size: 2048,
    series: 'Test Series',
    serno: '2',
  },
  {
    id: '3',
    lang: 'ru',
    authors: 'Иван Петров',
    title: 'Война и мир',
    size: 5120,
    series: 'Классика',
    serno: '1',
  },
]

test.describe('Visual Regression Tests', () => {
  test('should match homepage snapshot', async ({ page }) => {
    // arrange
    const homePage = new HomePage(page)
    await homePage.goto()

    // assert
    await expect(page).toHaveScreenshot('homepage-initial.png', {
      fullPage: false,
      mask: [page.locator('footer')], // Optional: mask dynamic regions
    })
  })

  test('should match search results snapshot', async ({ page }) => {
    // arrange
    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ json: mockBooks })
    })

    const homePage = new HomePage(page)
    await homePage.goto()
    await homePage.search('test')

    // Wait for results to render
    await expect(homePage.results.getByRole('link')).toHaveCount(3)

    // assert
    await expect(page).toHaveScreenshot('search-results.png', {
      fullPage: false,
    })
  })

  test('should match empty results snapshot', async ({ page }) => {
    // arrange
    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ json: [] })
    })

    const homePage = new HomePage(page)
    await homePage.goto()
    await homePage.search('nonexistentbook12345')

    // assert
    await expect(page).toHaveScreenshot('empty-results.png', {
      fullPage: false,
    })
  })

  test('should match error state snapshot', async ({ page }) => {
    // arrange
    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ status: 500, json: { error: 'Internal server error' } })
    })

    const homePage = new HomePage(page)
    await homePage.goto()
    await homePage.search('error test')

    // Wait for error to appear
    await expect(page.locator('.text-red')).toBeVisible()

    // assert
    await expect(page).toHaveScreenshot('error-state.png', {
      fullPage: false,
    })
  })

  test('should match search results with Russian text', async ({ page }) => {
    // arrange
    const russianBooks = [
      {
        id: '1',
        lang: 'ru',
        authors: 'Лев Толстой',
        title: 'Анна Каренина',
        size: 3072,
        series: 'Романы',
        serno: '1',
      },
      {
        id: '2',
        lang: 'ru',
        authors: 'Фёдор Достоевский',
        title: 'Преступление и наказание',
        size: 4096,
        series: 'Романы',
        serno: '1',
      },
    ]

    await page.route('**/api/search*', async (route) => {
      await route.fulfill({ json: russianBooks })
    })

    const homePage = new HomePage(page)
    await homePage.goto()
    await homePage.search('толстой')

    // Wait for results to render
    await expect(homePage.results.getByRole('link')).toHaveCount(2)

    // assert
    await expect(page).toHaveScreenshot('russian-text-results.png', {
      fullPage: false,
    })
  })
})
