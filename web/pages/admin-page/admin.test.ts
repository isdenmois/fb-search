import userEvent from '@testing-library/user-event'
import { render, screen } from '@testing-library/vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/shared/api', () => ({
  api: {
    parse: {
      getProgress: vi.fn().mockResolvedValue(null),
      rebuild: vi.fn().mockResolvedValue(null),
    },
  },
}))

// Fresh module graph per test so the api-key singleton re-reads sessionStorage.
const importAdmin = () => import('./AdminPage.vue').then((m) => m.default)

describe('Admin', () => {
  beforeEach(() => {
    vi.resetModules()
    sessionStorage.clear()
  })

  it('renders masked api key input and rebuild button', async () => {
    // arrange
    const Admin = await importAdmin()

    // act
    render(Admin)

    // assert
    const input = screen.getByPlaceholderText('API key')
    expect(input).toHaveAttribute('type', 'password')
    expect(screen.getByRole('button')).toHaveTextContent('Rebuild Database')
  })

  it('stores the entered key in sessionStorage', async () => {
    // arrange
    const Admin = await importAdmin()
    render(Admin)
    const input = screen.getByPlaceholderText('API key')

    // act
    await userEvent.type(input, 'secret-key')

    // assert
    expect(sessionStorage.getItem('admin-api-key')).toBe('secret-key')
  })

  it('restores the key from sessionStorage', async () => {
    // arrange
    sessionStorage.setItem('admin-api-key', 'stored-key')

    // act
    const Admin = await importAdmin()
    render(Admin)

    // assert
    expect(screen.getByPlaceholderText('API key')).toHaveValue('stored-key')
  })

  it('passes the entered key when rebuilding', async () => {
    // arrange
    const Admin = await importAdmin()
    const { api } = await import('@/shared/api')
    render(Admin)
    await userEvent.type(screen.getByPlaceholderText('API key'), 'secret-key')

    // act
    await userEvent.click(screen.getByRole('button'))

    // assert
    expect(api.parse.rebuild).toHaveBeenCalledWith('secret-key')
  })
})
