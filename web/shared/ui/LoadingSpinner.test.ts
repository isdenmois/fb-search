import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/vue'
import LoadingSpinner from './LoadingSpinner.vue'

describe('LoadingSpinner', () => {
  it('renders with role progressbar', () => {
    // arrange
    render(LoadingSpinner)

    // act
    const spinner = screen.getByRole('progressbar')

    // assert
    expect(spinner).toBeTruthy()
  })
})
