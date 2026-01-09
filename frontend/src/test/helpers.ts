import { vi } from 'vitest'

/**
 * Poll condition until true or timeout
 * @param fn - Function that returns true when condition is met
 * @param timeout - Timeout in milliseconds (default 1000)
 * @returns Promise that resolves when condition is met or rejects on timeout
 */
export async function waitFor(fn: () => boolean, timeout = 1000): Promise<void> {
  const startTime = Date.now()
  const pollInterval = 50

  return new Promise((resolve, reject) => {
    const interval = setInterval(() => {
      if (fn()) {
        clearInterval(interval)
        resolve()
      } else if (Date.now() - startTime >= timeout) {
        clearInterval(interval)
        reject(new Error(`waitFor timeout after ${timeout}ms`))
      }
    }, pollInterval)
  })
}

/**
 * Wait for all pending promises to resolve
 * Uses setTimeout(0) to flush the microtask queue
 */
export async function flushPromises(): Promise<void> {
  return new Promise((resolve) => {
    setTimeout(resolve, 0)
  })
}

/**
 * Mock global fetch with predefined responses by URL pattern
 * @param responses - Record of URL patterns to response data
 */
export function mockFetch(responses: Record<string, any>): void {
  global.fetch = vi.fn((url: string | URL | Request) => {
    const urlString = typeof url === 'string' ? url : url.toString()

    // Find matching response by checking if URL contains the pattern
    for (const [pattern, responseData] of Object.entries(responses)) {
      if (urlString.includes(pattern)) {
        return Promise.resolve({
          ok: true,
          status: 200,
          json: async () => responseData,
          text: async () => JSON.stringify(responseData),
          headers: new Headers({ 'Content-Type': 'application/json' }),
        } as Response)
      }
    }

    // Default 404 response
    return Promise.resolve({
      ok: false,
      status: 404,
      json: async () => ({ error: 'Not found' }),
      text: async () => JSON.stringify({ error: 'Not found' }),
      headers: new Headers({ 'Content-Type': 'application/json' }),
    } as Response)
  }) as any
}

/**
 * Reset all vi.mock() calls
 */
export function resetMocks(): void {
  vi.clearAllMocks()
  vi.resetAllMocks()
}
