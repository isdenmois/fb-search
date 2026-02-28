import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/vue'
import Admin from './AdminPage.vue'

vi.mock('shared/api', () => ({
  api: {
    parse: {
      getProgress: vi.fn(),
      rebuild: vi.fn(),
    },
  },
}))

describe('Admin', () => {
  it('renders with red button', () => {
    render(Admin)

    const button = screen.getByRole('button') as HTMLButtonElement

    expect(button).toHaveClass('bg-red-500')
  })

  it('shows Rebuild Database text initially', () => {
    render(Admin)

    const buttons = screen.queryAllByRole('button')
    const button = buttons[0] as HTMLButtonElement

    expect(button).toHaveTextContent('Rebuild Database')
  })
})
