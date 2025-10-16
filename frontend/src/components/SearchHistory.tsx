"use client";

import { useChatStore, SearchHistoryItem } from "@/lib/store";
import { Clock, ChevronLeft, ChevronRight, Trash2, Search } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { ru, uk, enUS } from "date-fns/locale";

const localeMap: Record<string, any> = {
  ru,
  uk,
  en: enUS,
};

export function SearchHistory() {
  const {
    searchHistory,
    isSidebarOpen,
    toggleSidebar,
    loadSearchFromHistory,
    clearSearchHistory,
    language,
  } = useChatStore();

  const locale = localeMap[language] || enUS;

  const handleHistoryClick = (item: SearchHistoryItem) => {
    loadSearchFromHistory(item);
  };

  const getTimeAgo = (timestamp: number) => {
    try {
      return formatDistanceToNow(timestamp, { addSuffix: true, locale });
    } catch {
      return new Date(timestamp).toLocaleDateString();
    }
  };

  return (
    <>
      {/* Toggle Button */}
      <button
        onClick={toggleSidebar}
        className="fixed left-0 top-20 z-50 bg-primary text-primary-foreground p-2 rounded-r-lg shadow-lg hover:opacity-90 transition-opacity"
        aria-label={isSidebarOpen ? "Close sidebar" : "Open sidebar"}
      >
        {isSidebarOpen ? (
          <ChevronLeft className="w-5 h-5" />
        ) : (
          <ChevronRight className="w-5 h-5" />
        )}
      </button>

      {/* Sidebar Panel */}
      <div
        className={`fixed left-0 top-16 bottom-0 w-80 bg-background border-r border-border shadow-lg transform transition-transform duration-300 ease-in-out z-40 ${
          isSidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="p-4 border-b border-border">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-2">
                <Clock className="w-5 h-5 text-primary" />
                <h2 className="text-lg font-semibold">Search History</h2>
              </div>
              {searchHistory.length > 0 && (
                <button
                  onClick={clearSearchHistory}
                  className="p-1.5 hover:bg-secondary rounded-md transition-colors"
                  title="Clear history"
                >
                  <Trash2 className="w-4 h-4 text-muted-foreground" />
                </button>
              )}
            </div>
            <p className="text-sm text-muted-foreground">
              {searchHistory.length} {searchHistory.length === 1 ? "search" : "searches"}
            </p>
          </div>

          {/* History List */}
          <div className="flex-1 overflow-y-auto">
            {searchHistory.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full p-8 text-center">
                <Search className="w-12 h-12 text-muted-foreground/50 mb-3" />
                <p className="text-sm text-muted-foreground">
                  No search history yet
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  Your searches will appear here
                </p>
              </div>
            ) : (
              <div className="p-2 space-y-1">
                {searchHistory.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => handleHistoryClick(item)}
                    className="w-full text-left p-3 rounded-lg hover:bg-secondary transition-colors group"
                  >
                    <div className="flex items-start gap-2">
                      <Search className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate group-hover:text-primary transition-colors">
                          {item.query}
                        </p>
                        <div className="flex items-center gap-2 mt-1">
                          <span className="text-xs text-muted-foreground">
                            {getTimeAgo(item.timestamp)}
                          </span>
                          {item.category && (
                            <>
                              <span className="text-xs text-muted-foreground">•</span>
                              <span className="text-xs text-primary/70 capitalize">
                                {item.category}
                              </span>
                            </>
                          )}
                          {item.productsCount !== undefined && item.productsCount > 0 && (
                            <>
                              <span className="text-xs text-muted-foreground">•</span>
                              <span className="text-xs text-muted-foreground">
                                {item.productsCount} products
                              </span>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Overlay for mobile */}
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-black/20 z-30 lg:hidden"
          onClick={toggleSidebar}
        />
      )}
    </>
  );
}
