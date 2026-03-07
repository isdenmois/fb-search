import type { Page, Locator } from '@playwright/test'

export interface Book {
  id: string
  lang: string
  authors?: string
  title: string
  size?: number
  series?: string
  serno?: string
}

export class HomePage {
  readonly searchInput: Locator
  readonly results: Locator
  readonly loadingSpinner: Locator

  constructor(private readonly page: Page) {
    this.searchInput = page.getByRole('textbox')
    this.results = page.getByRole('list')
    this.loadingSpinner = page.getByRole('progressbar')
  }

  async goto(): Promise<void> {
    await this.page.goto('/')
  }

  async search(query: string): Promise<void> {
    await this.searchInput.fill(query)
    await this.searchInput.press('Enter')
  }

  async getBookItems(): Promise<Locator[]> {
    return this.results.getByRole('link').all()
  }

  async getFirstBookTitle(): Promise<string> {
    const firstBook = this.results.getByRole('link').first()
    return (await firstBook.textContent()) || ''
  }
}
