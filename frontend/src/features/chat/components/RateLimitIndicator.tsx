"use client";

import { useEffect, useState } from "react";
import { rateLimitTracker, RateLimitInfo } from "@/shared/lib/rate-limit-tracker";

/**
 * RateLimitIndicator component
 * Shows current rate limit status when approaching the limit
 * Only displays when less than 20% remaining
 */
export function RateLimitIndicator() {
  const [info, setInfo] = useState<RateLimitInfo | null>(null);

  useEffect(() => {
    return rateLimitTracker.subscribe(setInfo);
  }, []);

  // Only show when less than 20% remaining
  if (!info || info.percentage > 20) {
    return null;
  }

  const getColor = () => {
    if (info.percentage < 10) return "text-red-500 dark:text-red-400";
    return "text-yellow-600 dark:text-yellow-500";
  };

  const getIcon = () => {
    if (info.percentage < 10) {
      return (
        <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
            clipRule="evenodd"
          />
        </svg>
      );
    }
    return (
      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
        <path
          fillRule="evenodd"
          d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
          clipRule="evenodd"
        />
      </svg>
    );
  };

  return (
    <div className={`flex items-center gap-1 text-xs ${getColor()} px-2 py-1`}>
      {getIcon()}
      <span>
        {info.remaining}/{info.limit} requests remaining
      </span>
    </div>
  );
}
