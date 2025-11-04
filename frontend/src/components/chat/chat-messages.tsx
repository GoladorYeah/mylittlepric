"use client";

import { useEffect, useRef } from "react";
import { ChatMessage as ChatMessageComponent } from "../ChatMessage";
import { LoadingDots } from "../ui/loading-dots";
import { ChatMessage } from "@/types";
import { ChatEmptyState } from "./chat-empty-state";
import { useChatStore } from "@/lib/store";

interface ChatMessagesProps {
  messages: ChatMessage[];
  isLoading: boolean;
  onQuickReply: (reply: string) => void;
}

export function ChatMessages({
  messages,
  isLoading,
  onQuickReply,
}: ChatMessagesProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { sessionId } = useChatStore();

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  return (
    <div className="flex-1 overflow-y-auto">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {messages.length === 0 ? (
          <ChatEmptyState />
        ) : (
          <div key={sessionId} className="space-y-6">
            {messages.map((message) => (
              <ChatMessageComponent
                key={message.id}
                message={message}
                onQuickReply={onQuickReply}
              />
            ))}
            {isLoading && <LoadingDots />}
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>
    </div>
  );
}
