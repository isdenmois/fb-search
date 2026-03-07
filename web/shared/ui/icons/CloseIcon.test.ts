import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/vue'
import CloseIcon from './CloseIcon.vue'

describe('CloseIcon', () => {
  it('renders SVG element', () => {
    // arrange
    render(CloseIcon)

    // act
    const svg = screen.getByRole('img')

    // assert
    expect(svg).toBeTruthy()
  })

  it('has correct viewBox', () => {
    // arrange
    render(CloseIcon)

    // act
    const svg = screen.getByRole('img')

    // assert
    expect(svg).toHaveAttribute('viewBox', '0 0 320 512')
  })
})
