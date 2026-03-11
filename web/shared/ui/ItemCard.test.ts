import { render, screen } from '@testing-library/vue'
import { describe, expect, it } from 'vitest'
import ItemCard from './ItemCard.vue'

describe('ItemCard', () => {
  it('renders with required title', () => {
    // arrange
    render(ItemCard, { props: { title: 'Test Title' } })

    // act
    const title = screen.getByText('Test Title')

    // assert
    expect(title).toBeTruthy()
  })

  it('renders suptitle when provided', () => {
    // arrange
    render(ItemCard, { props: { title: 'Title', suptitle: 'Subtitle' } })

    // act
    const suptitle = screen.getByText('Subtitle')

    // assert
    expect(suptitle).toBeTruthy()
  })

  it('hides suptitle when not provided', () => {
    // arrange
    const { queryByText } = render(ItemCard, { props: { title: 'Title' } })

    // act
    const suptitle = queryByText('Subtitle')

    // assert
    expect(suptitle).toBeFalsy()
  })

  it('renders subtitle when provided', () => {
    // arrange
    render(ItemCard, { props: { title: 'Title', subtitle: 'Description' } })

    // act
    const subtitle = screen.getByText('Description')

    // assert
    expect(subtitle).toBeTruthy()
  })

  it('hides subtitle when not provided', () => {
    // arrange
    render(ItemCard, { props: { title: 'Title' } })

    // act
    const subtitle = screen.queryByText('Description')

    // assert
    expect(subtitle).toBeFalsy()
  })

  it('renders icon slot content', () => {
    // arrange
    render(ItemCard, {
      props: { title: 'Title' },
      slots: { icon: '<span>Icon</span>' },
    })

    // act
    const icon = screen.getByText('Icon')

    // assert
    expect(icon).toBeTruthy()
  })
})
