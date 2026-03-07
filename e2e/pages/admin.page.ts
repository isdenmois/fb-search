import type { Page, Locator } from '@playwright/test'

export interface ParseProgress {
  files: number
  books: number
  time: string
}

export class AdminPage {
  readonly rebuildButton: Locator
  readonly progressSection: Locator
  readonly filesInfo: Locator
  readonly booksInfo: Locator
  readonly timeInfo: Locator

  constructor(private readonly page: Page) {
    this.rebuildButton = page.getByText('Rebuild Database')
    this.progressSection = page.locator('div:has-text("Files:")')
    this.filesInfo = page.locator('p:has-text("Files:")')
    this.booksInfo = page.locator('p:has-text("Books:")')
    this.timeInfo = page.locator('p:has-text("Time:")')
  }

  async goto(): Promise<void> {
    await this.page.goto('/admin')
  }

  async clickRebuild(): Promise<void> {
    await this.rebuildButton.click()
  }

  async getProgressInfo(): Promise<{ files: string; books: string; time: string } | null> {
    const isVisible = await this.progressSection.isVisible()
    if (!isVisible) {
      return null
    }

    const filesText = await this.filesInfo.textContent()
    const booksText = await this.booksInfo.textContent()
    const timeText = await this.timeInfo.textContent()

    return {
      files: filesText || '',
      books: booksText || '',
      time: timeText || '',
    }
  }

  async isRebuildButtonDisabled(): Promise<boolean> {
    return await this.rebuildButton.isDisabled()
  }

  async getRebuildButtonText(): Promise<string> {
    return (await this.rebuildButton.textContent()) || ''
  }
}
