"use client";

import { useEffect, useRef } from "react";

import { ChatMessage } from "@/shared/types";

import { ChatMessage as ChatMessageComponent } from "./ChatMessage";
import { LoadingDots } from "@/shared/components/ui";

import { ChatEmptyState } from "./chat-empty-state";
import { useChatStore } from "@/shared/lib";

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
    <div className="w-full max-w-4xl mx-auto px-4 pt-8 pb-24">
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
  );
}
