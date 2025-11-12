"use client";

import { useEffect, useState } from "react";
import { useChatStore } from "@/shared/lib";

/**
 * RateLimitNotification component
 * Displays a notification when user exceeds rate limits
 * Shows countdown timer until they can send messages again
 */
export function RateLimitNotification() {
  const { rateLimitState } = useChatStore();
  const [timeRemaining, setTimeRemaining] = useState<number>(0);

  useEffect(() => {
    if (!rateLimitState.isLimited || !rateLimitState.expiresAt) {
      return;
    }

    const interval = setInterval(() => {
      const now = Date.now();
      const remaining = Math.max(
        0,
        Math.floor((rateLimitState.expiresAt.getTime() - now) / 1000)
      );
      setTimeRemaining(remaining);

      if (remaining === 0) {
        clearInterval(interval);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [rateLimitState.isLimited, rateLimitState.expiresAt]);

  if (!rateLimitState.isLimited) {
    return null;
  }

  return (
    <div className="fixed top-4 right-4 z-50 max-w-md p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg shadow-lg animate-in fade-in slide-in-from-top-2 duration-300">
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0">
          <svg
            className="w-5 h-5 text-yellow-600 dark:text-yellow-500"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fillRule="evenodd"
              d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
              clipRule="evenodd"
            />
          </svg>
        </div>
        <div className="flex-1">
          <h3 className="text-sm font-medium text-yellow-800 dark:text-yellow-200">
            Rate Limit Exceeded
          </h3>
          <p className="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
            {rateLimitState.reason ||
              "You've sent too many messages. Please wait before sending more."}
          </p>
          {timeRemaining > 0 && (
            <p className="mt-2 text-sm font-semibold text-yellow-800 dark:text-yellow-200">
              Retry in {timeRemaining} second{timeRemaining !== 1 ? "s" : ""}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
