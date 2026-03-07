import { WretchError } from 'wretch/resolver'

export const getErrorMessage = (error: unknown): string => {
  if (error && typeof error === 'object') {
    if ('json' in error && error instanceof WretchError) {
      return getErrorMessage(error.json)
    }

    if ('message' in error) {
      return String(error.message)
    }
  }

  return String(error)
}
