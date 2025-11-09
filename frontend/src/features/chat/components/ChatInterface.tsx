"use client";

import { useChat } from "@/shared/hooks";
import { useChatStore } from "@/shared/lib";


import { SearchHistory } from "@/features/search";
import { ChatMessages } from "./chat-messages";
import { ChatInput } from "./chat-input";


interface ChatInterfaceProps {
  initialQuery?: string;
}

export function ChatInterface({ initialQuery }: ChatInterfaceProps) {
  const { messages, isLoading } = useChatStore();

  const { sendMessage, handleNewSearch, connectionStatus, isConnected } =
    useChat({ initialQuery });

  const handleQuickReply = (reply: string) => {
    sendMessage(reply);
  };

  const { isSidebarOpen } = useChatStore();

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
        className={`min-h-screen w-full bg-gradient-to-br from-background via-background to-background/95 transition-[padding] duration-300 ease-in-out will-change-[padding] ${
          isSidebarOpen ? 'lg:pl-80' : 'lg:pl-16'
        }`}
      >
        {/* Messages - scrolls naturally with page */}
        <ChatMessages
          messages={messages}
          isLoading={isLoading}
          onQuickReply={handleQuickReply}
        />

        {/* Fixed input field - always visible at bottom */}
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
      </div>
    </>
  );
}
