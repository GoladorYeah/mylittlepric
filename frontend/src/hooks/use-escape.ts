import { useEffect } from "react";

/**
 * Custom hook to handle Escape key press
 * @param handler - Callback function to execute when Escape is pressed
 * @param enabled - Whether the hook should be active (default: true)
 */
export function useEscape(
  handler: (event: KeyboardEvent) => void,
  enabled: boolean = true
) {
  useEffect(() => {
    if (!enabled) return;

    const listener = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        handler(event);
      }
    };

    document.addEventListener("keydown", listener);

    return () => {
      document.removeEventListener("keydown", listener);
    };
  }, [handler, enabled]);
}
