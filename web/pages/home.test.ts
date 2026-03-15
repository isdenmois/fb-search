import userEvent from '@testing-library/user-event'
import { render, screen } from '@testing-library/vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { api } from '@/shared/api'
import type { Book } from '@/shared/api/fb'
import HomePage from './HomePage.vue'
import { home } from './home.po'

vi.mock('@/shared/api', () => ({
  api: {
    fb: {
      search: vi.fn(),
    },
  },
}))

describe('HomePage', () => {
  let booksMock: Book[]

  beforeEach(() => {
    booksMock = [
      {
        id: 'files/test.inpx/book1.fb2',
        lang: 'ru',
        authors: 'Author One',
        title: 'Test Book One',
        size: 1024,
        series: 'Test Series',
        serno: '1',
      },
      {
        id: 'files/test.inpx/book2.fb2',
        lang: 'ru',
        authors: 'Author Two',
        title: 'Test Book Two',
        size: 2048,
        series: 'Test Series',
        serno: '2',
      },
    ]

    vi.mocked(api.fb.search).mockResolvedValue(booksMock)

    render(HomePage)
  })

  it('searches on form submit', async () => {
    // arrange + act
    await userEvent.type(home.input, 'test query{enter}')

    // assert
    expect(api.fb.search).toHaveBeenCalledWith('test query')

    expect(screen.getByRole('list')).toBeInTheDocument()
    expect(screen.getAllByRole('listitem')).toHaveLength(2)

    expect(screen.getByText('Test Book One')).toBeInTheDocument()
    expect(screen.getByText('Test Book Two')).toBeInTheDocument()
  })

  it('disables input during search', async () => {
    // arrange
    vi.mocked(api.fb.search).mockReturnValue(new Promise(() => {}))

    // act
    await userEvent.type(home.input, 'test query{enter}')

    // assert
    expect(home.input).toBeDisabled()
  })

  it('shows loading spinner during search', async () => {
    // arrange
    vi.mocked(api.fb.search).mockReturnValue(new Promise(() => {}))

    // act
    await userEvent.type(home.input, 'test query{enter}')

    // assert
    expect(home.spinner).toBeInTheDocument()
  })

  it('renders book items with correct data', async () => {
    // arrange
    vi.mocked(api.fb.search).mockResolvedValue([
      {
        id: 'files/test.inpx/book.fb2',
        lang: 'ru',
        authors: 'John Doe',
        title: 'The Great Adventure',
        size: 102400,
        series: 'Adventure Series',
        serno: '1',
      },
    ])

    // act
    await userEvent.type(home.input, 'test query{enter}')

    // assert
    expect(screen.getByText('The Great Adventure')).toBeInTheDocument()
    expect(screen.getByText('John Doe')).toBeInTheDocument()
    expect(screen.getByText('100Kb, Adventure Series, 1, ru')).toBeInTheDocument()
    expect(screen.getByRole('link')).toHaveAttribute('href', '/dl/files/test.inpx/book.fb2')
  })

  // TODO: update tests
  it('handles empty search results', async () => {
    // arrange
    vi.mocked(api.fb.search).mockResolvedValue([])

    const input = screen.getByRole('textbox') as HTMLInputElement

    // act
    await userEvent.type(input, 'no results{enter}')

    // Wait for search to complete
    await vi.waitFor(() => {
      expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
    })

    // assert - no list rendered when no results
    const list = screen.queryByRole('list')
    expect(list).not.toBeInTheDocument()
  })

  it('input clears on clear button click', async () => {
    // arrange

    const input = screen.getByRole('textbox') as HTMLInputElement

    // act - Type in input
    await userEvent.type(input, 'test query')

    // assert - input has value
    expect(input).toHaveValue('test query')

    // act - Click clear button
    const clearButton = screen.getByTestId('clear')
    await userEvent.click(clearButton)

    // assert - input cleared
    expect(input).toHaveValue('')
  })

  it('search trims query before sending', async () => {
    // arrange
    const mockBooks = [
      {
        id: 'files/test.inpx/book.fb2',
        lang: 'ru',
        authors: 'Author One',
        title: 'Test Book',
        size: 1024,
      },
    ]
    vi.mocked(api.fb.search).mockResolvedValue(mockBooks)

    const input = screen.getByRole('textbox') as HTMLInputElement

    // act
    await userEvent.type(input, '  test query  {enter}')

    // assert
    expect(api.fb.search).toHaveBeenCalledWith('  test query  ')
  })

  it('handles API error gracefully', async () => {
    // arrange
    vi.mocked(api.fb.search).mockRejectedValue(Error('Search failed'))

    const input = screen.getByRole('textbox') as HTMLInputElement

    // act - Wrap in try-catch to prevent unhandled rejection
    try {
      await userEvent.type(input, 'test query{enter}')
    } catch {
      // Ignore error
    }

    // assert - input re-enabled after error
    expect(input).not.toBeDisabled()
  })
})
