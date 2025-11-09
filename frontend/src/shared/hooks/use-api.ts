import { useState, useCallback } from "react";

export interface UseApiOptions<T> {
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
  retries?: number;
  retryDelay?: number;
}

export interface UseApiReturn<T, P extends unknown[]> {
  data: T | null;
  error: Error | null;
  loading: boolean;
  execute: (...args: P) => Promise<T | null>;
  reset: () => void;
}

/**
 * Custom hook for API calls with loading, error states, and retry logic
 * @param apiFunction - Async function to execute
 * @param options - Configuration options
 * @returns Object with data, error, loading, execute, and reset
 */
export function useApi<T, P extends unknown[] = []>(
  apiFunction: (...args: P) => Promise<T>,
  options: UseApiOptions<T> = {}
): UseApiReturn<T, P> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const [loading, setLoading] = useState(false);

  const { onSuccess, onError, retries = 0, retryDelay = 1000 } = options;

  const execute = useCallback(
    async (...args: P): Promise<T | null> => {
      setLoading(true);
      setError(null);

      let lastError: Error | null = null;
      let attemptCount = 0;

      while (attemptCount <= retries) {
        try {
          const result = await apiFunction(...args);
          setData(result);
          setLoading(false);
          onSuccess?.(result);
          return result;
        } catch (err) {
          lastError = err instanceof Error ? err : new Error(String(err));
          attemptCount++;

          if (attemptCount <= retries) {
            // Wait before retrying
            await new Promise((resolve) => setTimeout(resolve, retryDelay));
          }
        }
      }

      // All retries failed
      setError(lastError);
      setLoading(false);
      onError?.(lastError!);
      return null;
    },
    [apiFunction, onSuccess, onError, retries, retryDelay]
  );

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  return { data, error, loading, execute, reset };
}
