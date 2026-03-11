import { render, screen } from '@testing-library/vue'
import { describe, expect, it } from 'vitest'
import SearchIcon from './SearchIcon.vue'

describe('SearchIcon', () => {
  it('renders SVG element', () => {
    // arrange
    render(SearchIcon)

    // act
    const svg = screen.getByRole('img')

    // assert
    expect(svg).toBeTruthy()
  })
})
