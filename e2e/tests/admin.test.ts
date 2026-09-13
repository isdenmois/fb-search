import { test, expect } from '@playwright/test'
import { AdminPage } from '../pages/admin.page'
import { parserFixture } from '../fixtures'

test.describe('Admin Page', () => {
  test('should display admin page with rebuild button', async ({ page }) => {
    // arrange
    const adminPage = new AdminPage(page)

    await parserFixture.mockParse(page)

    // act
    await adminPage.goto()

    // assert
    await expect(page.locator('h1')).toHaveText('Admin')
    await expect(page.locator('button')).toBeVisible()
    await expect(page.locator('button')).toHaveText('Rebuild Database')
  })

  test('should display progress after page load', async ({ page }) => {
    // arrange
    const adminPage = new AdminPage(page)

    await parserFixture.mockParse(page)

    // act
    await adminPage.goto()

    // assert
    await expect(adminPage.filesInfo).toBeVisible()
    const filesText = await adminPage.filesInfo.textContent()
    expect(filesText).toContain('100')
  })

  test('should trigger database rebuild on button click', async ({ page }) => {
    // arrange
    const adminPage = new AdminPage(page)

    await parserFixture.mockParse(page)
    await adminPage.goto()

    // act
    await adminPage.clickRebuild()

    // assert - verify rebuild completed and progress is shown
    await expect(adminPage.filesInfo).toBeVisible()
    const filesText = await adminPage.filesInfo.textContent()
    expect(filesText).toContain('100') // Verify the mocked progress data appeared
  })

  test('should handle rebuild API error gracefully', async ({ page }) => {
    // arrange
    const adminPage = new AdminPage(page)

    await parserFixture.mockParse(page, { error: 'Internal Server Error' }, 500)

    await adminPage.goto()

    // act
    await adminPage.clickRebuild()

    // assert
    const button = page.locator('button')
    await expect(button).toBeEnabled()
  })

  test('should send api key header and stay enabled on unauthorized rebuild', async ({ page }) => {
    // arrange
    const adminPage = new AdminPage(page)

    await parserFixture.mockParse(page, { error: 'Unauthorized' }, 401)

    await adminPage.goto()
    await adminPage.fillApiKey('wrong-key')

    const rebuildRequest = page.waitForRequest(
      (req) => req.url().includes('/api/parse/rebuild') && req.method() === 'POST',
    )

    // act
    await adminPage.clickRebuild()
    const request = await rebuildRequest

    // assert
    expect(request.headers()['x-api-key']).toBe('wrong-key')

    const button = page.locator('button')
    await expect(button).toBeEnabled()
  })
})
