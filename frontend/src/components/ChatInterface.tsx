"use client";

import { useChat } from "@/hooks";
import { useChatStore } from "@/lib/store";
import { SearchHistory } from "./SearchHistory";
import { ChatMessages } from "./chat/chat-messages";
import { ChatInput } from "./chat/chat-input";

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

        {/* Sticky input field - always visible */}
        <div className="sticky bottom-0 bg-background">
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
