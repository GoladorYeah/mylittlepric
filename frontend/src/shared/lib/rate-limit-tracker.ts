/**
 * RateLimitInfo - Information about current rate limit status
 */
export interface RateLimitInfo {
  limit: number;
  remaining: number;
  reset: Date;
  percentage: number; // 0-100
}

/**
 * RateLimitTracker - Tracks rate limit information from HTTP headers
 * Monitors X-RateLimit-* headers to provide proactive warnings
 */
export class RateLimitTracker {
  private info: RateLimitInfo | null = null;
  private listeners: Set<(info: RateLimitInfo | null) => void> = new Set();

  /**
   * Update rate limit info from HTTP response headers
   * Parses X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset
   */
  updateFromHeaders(headers: Headers) {
    const limit = headers.get("X-RateLimit-Limit");
    const remaining = headers.get("X-RateLimit-Remaining");
    const reset = headers.get("X-RateLimit-Reset");

    if (limit && remaining && reset) {
      const limitNum = parseInt(limit, 10);
      const remainingNum = parseInt(remaining, 10);
      const resetNum = parseInt(reset, 10);

      this.info = {
        limit: limitNum,
        remaining: remainingNum,
        reset: new Date(resetNum * 1000), // Convert Unix timestamp to Date
        percentage: (remainingNum / limitNum) * 100,
      };

      this.notifyListeners();
    }
  }

  /**
   * Subscribe to rate limit info changes
   * Returns unsubscribe function
   */
  subscribe(listener: (info: RateLimitInfo | null) => void): () => void {
    this.listeners.add(listener);
    // Immediately call with current info
    listener(this.info);
    return () => this.listeners.delete(listener);
  }

  /**
   * Notify all listeners of rate limit info change
   */
  private notifyListeners() {
    this.listeners.forEach((listener) => listener(this.info));
  }

  /**
   * Get current rate limit info
   */
  getInfo(): RateLimitInfo | null {
    return this.info;
  }

  /**
   * Check if we're near the rate limit
   * @param threshold - Percentage threshold (default 10%)
   * @returns true if remaining percentage is below threshold
   */
  isNearLimit(threshold = 10): boolean {
    return this.info ? this.info.percentage < threshold : false;
  }

  /**
   * Reset rate limit info
   */
  reset() {
    this.info = null;
    this.notifyListeners();
  }
}

// Create singleton instance
export const rateLimitTracker = new RateLimitTracker();
