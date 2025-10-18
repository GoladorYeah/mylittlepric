"use client";

import { useChat } from "@/hooks";
import { useChatStore } from "@/lib/store";
import { SearchHistory } from "./SearchHistory";
import { ChatHeader } from "./chat/chat-header";
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

  return (
    <div className="flex flex-col h-screen bg-background">
      <SearchHistory />

      <ChatHeader
        isConnected={isConnected}
        connectionStatus={connectionStatus}
        onNewSearch={handleNewSearch}
      />

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
  );
}
