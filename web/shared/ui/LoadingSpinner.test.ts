import { render, screen } from '@testing-library/vue'
import { describe, expect, it } from 'vitest'
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
