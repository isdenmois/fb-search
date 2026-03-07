import { screen } from '@testing-library/vue'

export const home = {
  get input(): HTMLInputElement {
    return screen.getByRole('textbox')
  },

  get spinner(): HTMLProgressElement {
    return screen.getByRole('progressbar')
  },
}
