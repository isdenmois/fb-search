import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/vue'
import BookItem from './BookItem.vue'

function createBook(overrides?: Partial<Book>): Book {
  return {
    id: '1',
    title: 'Test Book',
    lang: 'en',
    authors: 'Test Author',
    size: 102400,
    series: 'Test Series',
    serno: '123',
    ...overrides,
  }
}

describe('BookItem', () => {
  it('renders book title and authors', () => {
    // arrange
    const book = createBook()
    render(BookItem, { props: { book } })

    // act
    const title = screen.getByText(book.title)
    const authors = screen.getByText(book.authors!)

    // assert
    expect(title).toBeInTheDocument()
    expect(authors).toBeInTheDocument()
  })

  it('formats subtitle with size, series, serno, and language', () => {
    // arrange
    const book = createBook({ size: 102400 }) // 100Kb
    render(BookItem, { props: { book } })

    // act
    const subtitle = screen.getByText(/100Kb.*Test Series.*123.*en/)

    // assert
    expect(subtitle).toBeTruthy()
  })

  it('handles missing optional fields', () => {
    // arrange
    const book = createBook({
      size: undefined,
      series: undefined,
      serno: undefined,
      lang: 'ru',
    })
    render(BookItem, { props: { book } })

    // act
    const subtitle = screen.getByText('0Kb, ru')

    // assert
    expect(subtitle).toBeTruthy()
  })

  it('renders FB2 icon', () => {
    // arrange
    const book = createBook()
    render(BookItem, { props: { book } })

    // act
    const svg = screen.getByRole('img')

    // assert
    expect(svg).toBeTruthy()
    expect(svg.textContent).toMatch(/FB2/i)
  })

  it('displays authors in suptitle position', () => {
    // arrange
    const book = createBook({ authors: 'Multiple Authors' })
    render(BookItem, { props: { book } })

    // act
    const suptitle = screen.getByText('Multiple Authors')
    const title = screen.getByText(book.title)

    // assert
    expect(suptitle).toBeTruthy()
    // Verify it appears before title (suptitle position)
    expect(suptitle.compareDocumentPosition(title) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
  })
})
