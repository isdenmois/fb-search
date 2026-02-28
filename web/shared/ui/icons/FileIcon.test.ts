import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/vue'
import FileIcon from './FileIcon.vue'

describe('FileIcon', () => {
  it('renders SVG element', () => {
    // arrange
    render(FileIcon, { props: { text: 'fb2' } })

    // act
    const svg = screen.getByRole('img')

    // assert
    expect(svg).toBeTruthy()
  })

  it('renders text prop content', () => {
    // arrange
    render(FileIcon, { props: { text: 'PDF' } })

    // act
    const text = screen.getByText('PDF')

    // assert
    expect(text).toBeTruthy()
  })
})
