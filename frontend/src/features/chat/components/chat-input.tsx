"use client";

import { useState } from "react";
import { Send } from "lucide-react";
import { CountrySelector } from "@/shared/components/ui";

interface ChatInputProps {
  onSend: (message: string) => void;
  isLoading: boolean;
  isConnected: boolean;
  connectionStatus: string;
}

export function ChatInput({
  onSend,
  isLoading,
  isConnected,
  connectionStatus,
}: ChatInputProps) {
  const [input, setInput] = useState("");

  const handleSend = () => {
    const trimmedInput = input.trim();
    if (trimmedInput) {
      onSend(trimmedInput);
      setInput("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto px-4 py-4 pb-2">
      <div className="flex items-center gap-2 px-4 py-2 rounded-xl bg-secondary border border-border focus-within:border-primary transition-colors">
        <CountrySelector />
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={
            isConnected ? "Type your message..." : `${connectionStatus}...`
          }
          disabled={isLoading || !isConnected}
          className="flex-1 px-2 py-3 bg-transparent focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
        />
        <button
          onClick={handleSend}
          disabled={!input.trim() || isLoading || !isConnected}
          className="w-10 h-10 rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer flex items-center justify-center shrink-0"
        >
          <Send className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
}
