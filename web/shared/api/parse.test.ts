import { beforeEach, describe, expect, it, vi } from 'vitest'
import { http } from './client'
import { getProgress, rebuild } from './parse'

vi.mock('./client', () => ({
  http: {
    get: vi.fn().mockReturnThis(),
    url: vi.fn().mockReturnThis(),
    post: vi.fn().mockReturnThis(),
    json: vi.fn(),
  },
}))

describe('parse API', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getProgress', () => {
    it('should get parse progress successfully', async () => {
      // Arrange
      const mockProgress = {
        files: 100,
        books: 50,
        time: '1.23s',
      }

      const mockJson = vi.fn().mockResolvedValue(mockProgress)
      ;(http.get as ReturnType<typeof vi.fn>).mockReturnValue({
        json: mockJson,
      })

      // Act
      const result = await getProgress()

      // Assert
      expect(result).toEqual(mockProgress)
      expect(http.get).toHaveBeenCalledWith('/parse')
    })

    it('should handle API error response', async () => {
      // Arrange
      const mockJson = vi.fn().mockRejectedValue(new Error('HTTP error! status: 404'))
      ;(http.get as ReturnType<typeof vi.fn>).mockReturnValue({
        json: mockJson,
      })

      // Act & Assert
      await expect(getProgress()).rejects.toThrow('HTTP error! status: 404')
    })

    it('should handle network error', async () => {
      // Arrange
      const mockJson = vi.fn().mockRejectedValue(new Error('Network error'))
      ;(http.get as ReturnType<typeof vi.fn>).mockReturnValue({
        json: mockJson,
      })

      // Act & Assert
      await expect(getProgress()).rejects.toThrow('Network error')
    })
  })

  describe('rebuild', () => {
    it('should trigger rebuild successfully', async () => {
      // Arrange
      const mockProgress = {
        files: 0,
        books: 0,
        time: '0s',
      }

      const mockUrl = vi.fn().mockReturnValue({
        post: vi.fn().mockReturnValue({
          json: vi.fn().mockResolvedValue(mockProgress),
        }),
      })

      const mockPost = vi.fn().mockReturnValue({
        json: vi.fn().mockResolvedValue(mockProgress),
      })

      http.url = mockUrl as typeof http.url
      ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
        post: mockPost,
        json: vi.fn().mockResolvedValue(mockProgress),
      })

      // Act
      const result = await rebuild()

      // Assert
      expect(result).toEqual(mockProgress)
      expect(http.url).toHaveBeenCalledWith('/parse/rebuild')
    })

    it('should handle rebuild API error response', async () => {
      // Arrange
      const mockPost = vi.fn().mockReturnValue({
        json: vi.fn().mockRejectedValue(new Error('HTTP error! status: 500')),
      })

      ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
        post: mockPost,
        json: vi.fn().mockResolvedValue({
          files: 0,
          books: 0,
          time: '0s',
        }),
      })

      // Act & Assert
      await expect(rebuild()).rejects.toThrow('HTTP error! status: 500')
    })

    it('should handle rebuild network error', async () => {
      // Arrange
      const mockPost = vi.fn().mockReturnValue({
        json: vi.fn().mockRejectedValue(new Error('Network error')),
      })

      ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({
        post: mockPost,
        json: vi.fn().mockResolvedValue({
          files: 0,
          books: 0,
          time: '0s',
        }),
      })

      // Act & Assert
      await expect(rebuild()).rejects.toThrow('Network error')
    })
  })
})
