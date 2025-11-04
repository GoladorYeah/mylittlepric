"use client";

import { useChat } from "@/hooks";
import { useChatStore } from "@/lib/store";
import { SearchHistory } from "./SearchHistory";
import { ChatMessages } from "./chat/chat-messages";
import { ChatInput } from "./chat/chat-input";
import { AdPlaceholder } from "./AdPlaceholder";

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
        className={`flex flex-col h-screen bg-gradient-to-br from-background via-background to-background/95 transition-[padding] duration-300 ease-in-out will-change-[padding] ${
          isSidebarOpen ? 'lg:pl-80' : 'lg:pl-16'
        }`}
      >
        {/* Top Banner Ad - Only show when there are messages */}
        {messages.length > 0 && (
          <div className="flex justify-center py-3 px-4 border-b border-border/30 bg-muted/20 backdrop-blur-sm">
            <AdPlaceholder format="banner" />
          </div>
        )}

        <ChatMessages
          messages={messages}
          isLoading={isLoading}
          onQuickReply={handleQuickReply}
        />

        <ChatInput
          onSend={sendMessage}
          isLoading={isLoading}
          isConnected={isConnected}
          connectionStatus={connectionStatus}
        />
      </div>
    </>
  );
}
