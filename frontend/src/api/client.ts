import { env } from '../utils/env';
import type { APIError } from '../types/api';

export class APIClientError extends Error {
  statusCode?: number;
  apiError?: APIError;

  constructor(message: string, statusCode?: number, apiError?: APIError) {
    super(message);
    this.name = 'APIClientError';
    this.statusCode = statusCode;
    this.apiError = apiError;
  }
}

interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>;
  timeout?: number;
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { params, timeout = 30000, ...fetchOptions } = options;

  // Build URL with query parameters
  let url = `${env.API_URL}${endpoint}`;
  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value));
      }
    });
    const queryString = searchParams.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
  }

  // Set up abort controller for timeout
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);

  try {
    // Log request in development
    if (env.IS_DEV) {
      console.log(`[API] ${options.method || 'GET'} ${url}`);
    }

    const response = await fetch(url, {
      ...fetchOptions,
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        ...fetchOptions.headers,
      },
    });

    clearTimeout(timeoutId);

    // Handle non-OK responses
    if (!response.ok) {
      let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
      let apiError: APIError | undefined;

      // Try to parse error response
      try {
        const errorData = await response.json();
        if (errorData.error || errorData.message) {
          apiError = errorData;
          errorMessage = errorData.message || errorData.error || errorMessage;
        }
      } catch {
        // If parsing fails, use the status text
      }

      throw new APIClientError(errorMessage, response.status, apiError);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return undefined as T;
    }

    // Handle 200 with potentially empty body (defensive fallback)
    if (response.status === 200) {
      const text = await response.text();

      // If body is empty, return undefined
      if (!text || text.trim() === '') {
        return undefined as T;
      }

      // Parse JSON response
      const data = JSON.parse(text);

      if (env.IS_DEV) {
        console.log(`[API] Response:`, data);
      }

      return data as T;
    }

    // For other status codes
    const data = await response.json();

    if (env.IS_DEV) {
      console.log(`[API] Response:`, data);
    }

    return data as T;
  } catch (error) {
    clearTimeout(timeoutId);

    // Handle abort error (timeout)
    if (error instanceof Error && error.name === 'AbortError') {
      throw new APIClientError('Request timeout', 0);
    }

    // Handle network errors
    if (error instanceof TypeError) {
      throw new APIClientError(
        'Network error: Unable to connect to the server. Please check your connection and try again.',
        0
      );
    }

    // Re-throw API errors
    if (error instanceof APIClientError) {
      throw error;
    }

    // Handle unknown errors
    throw new APIClientError(
      error instanceof Error ? error.message : 'An unknown error occurred',
      0
    );
  }
}

export const api = {
  get: <T>(url: string, params?: Record<string, any>) =>
    apiRequest<T>(url, { method: 'GET', params }),

  post: <T>(url: string, body?: any) =>
    apiRequest<T>(url, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    }),

  put: <T>(url: string, body?: any) =>
    apiRequest<T>(url, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    }),

  patch: <T>(url: string, body?: any) =>
    apiRequest<T>(url, {
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    }),

  delete: <T>(url: string) =>
    apiRequest<T>(url, { method: 'DELETE' }),
};
