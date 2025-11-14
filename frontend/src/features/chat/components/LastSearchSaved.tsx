"use client";

import { useState } from "react";
import { useChatStore } from "@/shared/lib";
import { RotateCcw, X, ChevronDown } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";
import { formatDistanceToNow } from "date-fns";
import { ru, uk, enUS } from "date-fns/locale";

const localeMap: Record<string, any> = {
  ru,
  uk,
  en: enUS,
};

export function LastSearchSaved() {
  const [isExpanded, setIsExpanded] = useState(false);

  const {
    savedSearch,
    restoreSavedSearch,
    clearSavedSearch,
    language,
  } = useChatStore();

  const router = useRouter();
  const pathname = usePathname();
  const locale = localeMap[language] || enUS;

  const handleRestoreSearch = () => {
    restoreSavedSearch();
    // Navigate to chat if we're not already there
    if (pathname !== '/chat') {
      router.push('/chat');
    }
  };

  const getSearchPreview = () => {
    if (!savedSearch || savedSearch.messages.length === 0) return "";
    // Get the last user message as preview
    const userMessages = savedSearch.messages.filter(m => m.role === "user");
    if (userMessages.length === 0) return "";
    const lastUserMsg = userMessages[userMessages.length - 1];
    return lastUserMsg.content.length > 50
      ? lastUserMsg.content.substring(0, 50) + "..."
      : lastUserMsg.content;
  };

  if (!savedSearch) {
    return null;
  }

  return (
    <div className="w-full max-w-md rounded-lg bg-amber-500/10 border border-amber-500/30">
      {/* Header - Always visible */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full px-3 py-2 flex items-center justify-between gap-2 hover:bg-amber-500/20 rounded-lg transition-colors"
      >
        <span className="text-xs font-medium text-amber-700 dark:text-amber-400">
          Last search saved
        </span>
        <ChevronDown
          className={`w-4 h-4 text-amber-700 dark:text-amber-400 transition-transform duration-200 ${
            isExpanded ? 'rotate-180' : ''
          }`}
        />
      </button>

      {/* Expanded Content */}
      {isExpanded && (
        <div className="px-3 pb-2.5 space-y-2">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1 min-w-0">
              <p className="text-xs text-muted-foreground truncate">
                {getSearchPreview()}
              </p>
              <p className="text-xs text-muted-foreground/70 mt-1">
                {formatDistanceToNow(new Date(savedSearch.timestamp), { addSuffix: true, locale })}
              </p>
            </div>
            <button
              onClick={(e) => {
                e.stopPropagation();
                clearSavedSearch();
              }}
              className="p-1 hover:bg-amber-500/20 rounded transition-colors shrink-0"
              title="Dismiss"
            >
              <X className="w-3.5 h-3.5 text-amber-700 dark:text-amber-400" />
            </button>
          </div>
          <button
            onClick={handleRestoreSearch}
            className="w-full px-3 py-1.5 bg-amber-500/20 hover:bg-amber-500/30 text-amber-700 dark:text-amber-400 rounded text-xs font-medium transition-colors flex items-center justify-center gap-2"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            Restore Search
          </button>
        </div>
      )}
    </div>
  );
}
