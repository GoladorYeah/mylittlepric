"use client";

import { useCallback } from "react";
import { useChat } from "@/shared/hooks";
import {
  useMessages,
  useLoadingState,
  useSidebarState,
  useSavedSearchPrompt,
  useSessionActions,
  useMessageActions,
  generateId,
  useChatStore
} from "@/shared/lib";
import { useAuthStore } from "@/shared/lib";

import { SearchHistory } from "@/features/search";
import { ChatMessages } from "./chat-messages";
import { ChatInput } from "./chat-input";
import { ChatHeader } from "./chat-header";
import { SavedSearchPrompt } from "./SavedSearchPrompt";
import { RateLimitNotification } from "./RateLimitNotification";
import { AuthPrompt } from "@/features/auth/components/AuthPrompt";


interface ChatInterfaceProps {
  initialQuery?: string;
  sessionId?: string;
}

export function ChatInterface({ initialQuery, sessionId }: ChatInterfaceProps) {
  // Use optimized selectors
  const messages = useMessages();
  const isLoading = useLoadingState();
  const { isSidebarOpen } = useSidebarState();
  const { showSavedSearchPrompt, setShowSavedSearchPrompt } = useSavedSearchPrompt();
  const { setSessionId, newSearch, restoreSavedSearch, clearSavedSearch } = useSessionActions();
  const { removeMessage } = useMessageActions();
  const searchState = useChatStore((state) => state.searchState);
  const { isAuthenticated } = useAuthStore();

  const { sendMessage, handleNewSearch, connectionStatus, isConnected } =
    useChat({ initialQuery, sessionId });

  // Check if user needs to authenticate (anonymous user hit search limit)
  const showAuthPrompt = !isAuthenticated &&
    searchState?.requires_authentication === true;

  const handleQuickReply = useCallback((reply: string) => {
    sendMessage(reply);
  }, [sendMessage]);

  const handleRetry = useCallback((messageId: string) => {
    // Find the failed message
    const message = messages.find((m) => m.id === messageId);
    if (message && message.content) {
      // Remove the failed message from store
      removeMessage(messageId);
      // Resend the message
      sendMessage(message.content);
    }
  }, [messages, removeMessage, sendMessage]);

  const handleContinueSearch = useCallback(() => {
    restoreSavedSearch();
    setShowSavedSearchPrompt(false);
  }, [restoreSavedSearch, setShowSavedSearchPrompt]);

  const handleStartNewSearch = useCallback(() => {
    clearSavedSearch();
    newSearch();
    const newSessionId = generateId();
    setSessionId(newSessionId);
    localStorage.setItem("chat_session_id", newSessionId);
    setShowSavedSearchPrompt(false);
  }, [clearSavedSearch, newSearch, setSessionId, setShowSavedSearchPrompt]);

  return (
    <>
      {/* Rate Limit Notification - Fixed position notification */}
      <RateLimitNotification />

      {/* Sidebar with all controls */}
      <SearchHistory
        isConnected={isConnected}
        connectionStatus={connectionStatus}
        onNewSearch={handleNewSearch}
      />

      {/* Main content area - pushed by sidebar on desktop */}
      <div
        className={`min-h-screen w-full bg-background transition-[padding] duration-300 ease-in-out will-change-[padding] ${
          isSidebarOpen ? 'lg:pl-80' : 'lg:pl-16'
        }`}
      >
        {/* Mobile Header */}
        <ChatHeader
          isConnected={isConnected}
          connectionStatus={connectionStatus}
          onNewSearch={handleNewSearch}
        />

        {/* Saved Search Prompt - shown when user returns to app with incomplete search */}
        {showSavedSearchPrompt ? (
          <SavedSearchPrompt
            onContinue={handleContinueSearch}
            onNewSearch={handleStartNewSearch}
          />
        ) : showAuthPrompt ? (
          /* Auth Prompt - shown when anonymous user hits search limit */
          <div className="container mx-auto flex min-h-[calc(100vh-200px)] max-w-2xl items-center justify-center px-4 py-8">
            <AuthPrompt
              searchesUsed={searchState?.anonymous_search_used || 0}
              searchesLimit={searchState?.anonymous_search_limit || 3}
            />
          </div>
        ) : (
          <>
            {/* Messages - scrolls naturally with page */}
            <ChatMessages
              messages={messages}
              isLoading={isLoading}
              onQuickReply={handleQuickReply}
              onRetry={handleRetry}
            />
          </>
        )}

        {/* Fixed input field - always visible at bottom (hidden when showing prompts) */}
        {!showSavedSearchPrompt && !showAuthPrompt && (
          <div
            className={`fixed bottom-0 left-0 right-0 z-30 transition-[padding] duration-300 ease-in-out ${
              isSidebarOpen ? 'lg:pl-80' : 'lg:pl-16'
            }`}
          >
            <ChatInput
              onSend={sendMessage}
              isLoading={isLoading}
              isConnected={isConnected}
              connectionStatus={connectionStatus}
            />
          </div>
        )}
      </div>
    </>
  );
}
