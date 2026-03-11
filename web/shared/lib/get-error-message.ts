import type { WretchError } from 'wretch'

function isWretchError(error: unknown): error is WretchError {
  return error instanceof Error && 'status' in error && typeof (error as WretchError).status === 'number'
}

export const getErrorMessage = (error: unknown): string => {
  if (error && typeof error === 'object') {
    if ('json' in error && isWretchError(error)) {
      return getErrorMessage(error.json)
    }

    if ('message' in error) {
      return String(error.message)
    }
  }

  return String(error)
}
