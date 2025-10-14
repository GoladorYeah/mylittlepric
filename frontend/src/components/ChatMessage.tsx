"use client";

import { ChatMessage as ChatMessageType } from "@/types";
import { ProductCard } from "./ProductCard";

interface ChatMessageProps {
  message: ChatMessageType;
  onQuickReply: (reply: string) => void;
}

export function ChatMessage({ message, onQuickReply }: ChatMessageProps) {
  const isUser = message.role === "user";

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`max-w-[80%] space-y-3 ${
          isUser ? "items-end" : "items-start"
        }`}
      >
        <div
          className={`rounded-2xl px-4 py-3 ${
            isUser
              ? "bg-primary text-primary-foreground"
              : "bg-secondary text-secondary-foreground"
          }`}
        >
          <p className="whitespace-pre-wrap">{message.content}</p>
        </div>

        {message.quick_replies && message.quick_replies.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {message.quick_replies.map((reply, index) => (
              <button
                key={index}
                onClick={() => onQuickReply(reply)}
                className="px-4 py-2 rounded-full bg-secondary hover:bg-secondary/80 text-sm transition-colors border border-border"
              >
                {reply}
              </button>
            ))}
          </div>
        )}

        {message.products && message.products.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full">
            {message.products.map((product, index) => (
              <ProductCard key={index} product={product} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}