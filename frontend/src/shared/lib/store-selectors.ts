/**
 * Optimized Zustand selectors to prevent unnecessary re-renders
 * Import these instead of using the full store in components
 */

import { useChatStore } from "./store";
import { useShallow } from "zustand/react/shallow";

/**
 * Selector for chat messages only
 * Use this in ChatMessages component
 */
export const useMessages = () =>
  useChatStore((state) => state.messages);

/**
 * Selector for loading state only
 * Use this in components that only need to show loading state
 */
export const useLoadingState = () =>
  useChatStore((state) => state.isLoading);

/**
 * Selector for session ID only
 */
export const useSessionId = () =>
  useChatStore((state) => state.sessionId);

/**
 * Selector for preferences (country, language, currency)
 * Use this in settings components
 */
export const usePreferences = () =>
  useChatStore(
    useShallow((state) => ({
      country: state.country,
      language: state.language,
      currency: state.currency,
    }))
  );

/**
 * Selector for preference actions only
 * Use this to avoid re-renders when preferences change
 */
export const usePreferenceActions = () =>
  useChatStore(
    useShallow((state) => ({
      setCountry: state.setCountry,
      setLanguage: state.setLanguage,
      setCurrency: state.setCurrency,
      syncPreferencesToServer: state.syncPreferencesToServer,
    }))
  );

/**
 * Selector for sidebar state
 */
export const useSidebarState = () =>
  useChatStore(
    useShallow((state) => ({
      isSidebarOpen: state.isSidebarOpen,
      toggleSidebar: state.toggleSidebar,
      setSidebarOpen: state.setSidebarOpen,
    }))
  );

/**
 * Selector for rate limit state
 */
export const useRateLimitState = () =>
  useChatStore((state) => state.rateLimitState);

/**
 * Selector for message actions
 * Use this to avoid re-renders when messages change
 */
export const useMessageActions = () =>
  useChatStore(
    useShallow((state) => ({
      addMessage: state.addMessage,
      setMessages: state.setMessages,
      clearMessages: state.clearMessages,
      removeMessage: state.removeMessage,
      updateMessageStatus: state.updateMessageStatus,
    }))
  );

/**
 * Selector for session actions
 */
export const useSessionActions = () =>
  useChatStore(
    useShallow((state) => ({
      setSessionId: state.setSessionId,
      newSearch: state.newSearch,
      loadSessionMessages: state.loadSessionMessages,
      saveCurrentSearch: state.saveCurrentSearch,
      restoreSavedSearch: state.restoreSavedSearch,
      clearSavedSearch: state.clearSavedSearch,
    }))
  );

/**
 * Selector for saved search prompt state
 */
export const useSavedSearchPrompt = () =>
  useChatStore(
    useShallow((state) => ({
      showSavedSearchPrompt: state.showSavedSearchPrompt,
      setShowSavedSearchPrompt: state.setShowSavedSearchPrompt,
      checkSavedSearchPrompt: state.checkSavedSearchPrompt,
    }))
  );
