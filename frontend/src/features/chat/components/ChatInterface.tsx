"use client";

import { useChat } from "@/shared/hooks";
import { useChatStore } from "@/shared/lib";
import { generateId } from "@/shared/lib";

import { SearchHistory } from "@/features/search";
import { ChatMessages } from "./chat-messages";
import { ChatInput } from "./chat-input";
import { ChatHeader } from "./chat-header";
import { SavedSearchPrompt } from "./SavedSearchPrompt";


interface ChatInterfaceProps {
  initialQuery?: string;
  sessionId?: string;
}

export function ChatInterface({ initialQuery, sessionId }: ChatInterfaceProps) {
  const {
    messages,
    isLoading,
    isSidebarOpen,
    showSavedSearchPrompt,
    setShowSavedSearchPrompt,
    restoreSavedSearch,
    clearSavedSearch,
    setSessionId,
    newSearch,
  } = useChatStore();

  const { sendMessage, handleNewSearch, connectionStatus, isConnected } =
    useChat({ initialQuery, sessionId });

  const handleQuickReply = (reply: string) => {
    sendMessage(reply);
  };

  const handleContinueSearch = () => {
    restoreSavedSearch();
    setShowSavedSearchPrompt(false);
  };

  const handleStartNewSearch = () => {
    clearSavedSearch();
    newSearch();
    const newSessionId = generateId();
    setSessionId(newSessionId);
    localStorage.setItem("chat_session_id", newSessionId);
    setShowSavedSearchPrompt(false);
  };

  return (
    <>
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
        ) : (
          <>
            {/* Messages - scrolls naturally with page */}
            <ChatMessages
              messages={messages}
              isLoading={isLoading}
              onQuickReply={handleQuickReply}
            />
          </>
        )}

        {/* Fixed input field - always visible at bottom (hidden when showing prompt) */}
        {!showSavedSearchPrompt && (
          <div
            className={`fixed bottom-0 left-0 right-0 bg-background z-30 transition-[padding] duration-300 ease-in-out ${
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
