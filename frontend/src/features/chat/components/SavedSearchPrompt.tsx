"use client";

import { RotateCcw, Plus } from "lucide-react";
import { useChatStore } from "@/shared/lib";

interface SavedSearchPromptProps {
  onContinue: () => void;
  onNewSearch: () => void;
}

export function SavedSearchPrompt({ onContinue, onNewSearch }: SavedSearchPromptProps) {
  const { savedSearch } = useChatStore();

  if (!savedSearch) return null;

  // Get preview of last user message
  const userMessages = savedSearch.messages.filter(m => m.role === "user");
  const lastUserMessage = userMessages[userMessages.length - 1];
  const preview = lastUserMessage?.content?.substring(0, 60) || "";

  return (
    <div className="flex items-center justify-center min-h-[60vh] px-4">
      <div className="max-w-lg w-full bg-card border border-border rounded-2xl shadow-xl p-8 space-y-6">
        {/* Icon */}
        <div className="flex justify-center">
          <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
            <RotateCcw className="w-8 h-8 text-primary" />
          </div>
        </div>

        {/* Title */}
        <div className="text-center space-y-2">
          <h2 className="text-2xl font-bold text-foreground">
            У вас есть незавершенный поиск
          </h2>
          <p className="text-muted-foreground">
            {preview && (
              <span className="block mt-2 text-sm italic">
                "{preview}{preview.length >= 60 ? '...' : ''}"
              </span>
            )}
          </p>
        </div>

        {/* Buttons */}
        <div className="flex flex-col sm:flex-row gap-3">
          <button
            onClick={onContinue}
            className="flex-1 px-6 py-3 bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg font-semibold transition-colors flex items-center justify-center gap-2 cursor-pointer"
          >
            <RotateCcw className="w-5 h-5" />
            Продолжить поиск
          </button>
          <button
            onClick={onNewSearch}
            className="flex-1 px-6 py-3 bg-secondary hover:bg-secondary/80 text-secondary-foreground rounded-lg font-semibold transition-colors flex items-center justify-center gap-2 cursor-pointer"
          >
            <Plus className="w-5 h-5" />
            Начать новый
          </button>
        </div>

        {/* Footer */}
        <p className="text-center text-xs text-muted-foreground">
          Вы можете продолжить незавершенный поиск или начать новый
        </p>
      </div>
    </div>
  );
}
