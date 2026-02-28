import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/vue'
import userEvent from '@testing-library/user-event'
import InputField from './InputField.vue'

describe.skip('InputField', () => {
  it('renders with correct value', () => {
    // arrange
    render(InputField, { props: { modelValue: 'test' } })

    // act
    const input = screen.getByRole('textbox')

    // assert
    expect(input).toHaveValue('test')
  })

  it('emits update event on input', async () => {
    // arrange
    const user = userEvent.setup()
    render(InputField, { props: { modelValue: '' } })

    // act
    const input = screen.getByRole('textbox')
    await user.type(input, 'hello')

    // assert
    expect(screen.getByRole('textbox')).toHaveValue('hello')
  })

  it('shows clear button when modelValue is non-empty', () => {
    // arrange
    render(InputField, { props: { modelValue: 'test' } })

    // act
    const clearButton = screen.queryByTestId('clear')

    // assert
    expect(clearButton).toBeTruthy()
  })

  it('hides clear button when modelValue is empty', () => {
    // arrange
    render(InputField, { props: { modelValue: '' } })

    // act
    const clearButton = screen.queryByTestId('clear')

    // assert
    expect(clearButton).toBeFalsy()
  })

  it('hides clear button when disabled', () => {
    // arrange
    render(InputField, { props: { modelValue: 'test', disabled: true } })

    // act
    const clearButton = screen.queryByTestId('clear')

    // assert
    expect(clearButton).toBeFalsy()
  })

  it('clear button emits empty string', async () => {
    // arrange
    const user = userEvent.setup()
    render(InputField, { props: { modelValue: 'test' } })

    // act
    const clearButton = screen.getByTestId('clear')
    await user.click(clearButton)

    // assert
    expect(screen.getByRole('textbox')).toHaveValue('')
  })

  it('applies disabled attribute when disabled prop is true', () => {
    // arrange
    render(InputField, { props: { modelValue: 'test', disabled: true } })

    // act
    const input = screen.getByRole('textbox')

    // assert
    expect(input).toBeDisabled()
  })
})
