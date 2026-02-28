import { describe, it, expect, vi, beforeEach } from 'vitest'
import { search } from './fb'
import { http } from './client'

vi.mock('./client', () => ({
  http: {
    url: vi.fn(() => ({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      post: vi.fn().mockReturnThis(),
      json: vi.fn(),
    })),
    get: vi.fn().mockReturnThis(),
  },
}))

describe('search API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should search books with valid query', async () => {
    // Arrange
    const mockBooks = [
      {
        id: 1,
        lang: 'ru',
        authors: 'Author One',
        title: 'Test Book One',
        size: 1024,
        series: 'Test Series',
        serno: '1',
      },
      {
        id: 2,
        lang: 'en',
        authors: 'Author Two',
        title: 'Test Book Two',
        size: 2048,
        series: 'Test Series',
        serno: '2',
      },
    ]

    const mockJson = vi.fn().mockResolvedValue(mockBooks)
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act
    const result = await search('test query')

    // Assert
    expect(result).toEqual(mockBooks)
    expect(http.url).toHaveBeenCalledWith('/search')
  })

  it('should handle search with empty query', async () => {
    // Arrange
    const mockJson = vi.fn().mockResolvedValue([])
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act
    const result = await search('')

    // Assert
    expect(result).toEqual([])
  })

  it('should handle search with whitespace-only query', async () => {
    // Arrange
    const mockJson = vi.fn().mockResolvedValue([])
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act
    const result = await search('   ')

    // Assert
    expect(result).toEqual([])
  })

  it('should trim and format query properly', async () => {
    // Arrange
    const mockBooks = [{ id: 1, lang: 'ru', title: 'Test' }]

    const mockJson = vi.fn().mockResolvedValue(mockBooks)
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act
    const query = '  test query  '
    await search(query)

    // Assert
    expect(http.url).toHaveBeenCalledWith('/search')
  })

  it('should handle API error response', async () => {
    // Arrange
    const mockJson = vi.fn().mockRejectedValue(new Error('HTTP error! status: 500'))
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act & Assert
    await expect(search('test')).rejects.toThrow('HTTP error! status: 500')
  })

  it('should handle network error', async () => {
    // Arrange
    const mockJson = vi.fn().mockRejectedValue(new Error('Network error'))
    ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
      query: vi.fn().mockReturnThis(),
      get: vi.fn().mockReturnThis(),
      json: mockJson,
    })

    // Act & Assert
    await expect(search('test')).rejects.toThrow('Network error')
  })
})
