import { beforeEach, describe, expect, it, vi } from 'vitest'
import { http } from './client'
import { getProgress, rebuild } from './parse'

vi.mock('./client', () => ({
  http: {
    get: vi.fn().mockReturnThis(),
    url: vi.fn().mockReturnThis(),
    headers: vi.fn().mockReturnThis(),
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
    const arrangeRebuild = (mockJson: ReturnType<typeof vi.fn>) => {
      const mockPost = vi.fn().mockReturnValue({ json: mockJson })
      const mockHeaders = vi.fn().mockReturnValue({ post: mockPost })
      ;(http.url as ReturnType<typeof vi.fn>).mockReturnValue({ headers: mockHeaders })

      return { mockHeaders, mockPost }
    }

    it('should trigger rebuild successfully with the api key header', async () => {
      // Arrange
      const mockProgress = {
        files: 0,
        books: 0,
        time: '0s',
      }

      const { mockHeaders, mockPost } = arrangeRebuild(vi.fn().mockResolvedValue(mockProgress))

      // Act
      const result = await rebuild('test-key')

      // Assert
      expect(result).toEqual(mockProgress)
      expect(http.url).toHaveBeenCalledWith('/parse/rebuild')
      expect(mockHeaders).toHaveBeenCalledWith({ 'X-API-Key': 'test-key' })
      expect(mockPost).toHaveBeenCalledWith({})
    })

    it('should handle rebuild API error response', async () => {
      // Arrange
      arrangeRebuild(vi.fn().mockRejectedValue(new Error('HTTP error! status: 500')))

      // Act & Assert
      await expect(rebuild('test-key')).rejects.toThrow('HTTP error! status: 500')
    })

    it('should handle rebuild network error', async () => {
      // Arrange
      arrangeRebuild(vi.fn().mockRejectedValue(new Error('Network error')))

      // Act & Assert
      await expect(rebuild('test-key')).rejects.toThrow('Network error')
    })
  })
})
