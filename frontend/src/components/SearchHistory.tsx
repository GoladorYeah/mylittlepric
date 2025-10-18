"use client";

import { useEffect, useState } from "react";
import { useChatStore } from "@/lib/store";
import { SearchHistoryAPI, type SearchHistoryItem as APISearchHistoryItem } from "@/lib/search-history-api";
import { Clock, ChevronLeft, ChevronRight, Trash2, Search, Package, RefreshCw } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { ru, uk, enUS } from "date-fns/locale";

const localeMap: Record<string, any> = {
  ru,
  uk,
  en: enUS,
};

export function SearchHistory() {
  const {
    isSidebarOpen,
    toggleSidebar,
    language,
  } = useChatStore();

  const [history, setHistory] = useState<APISearchHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const limit = 20;

  const locale = localeMap[language] || enUS;

  const loadHistory = async (resetOffset = false) => {
    try {
      setLoading(true);
      setError(null);
      const currentOffset = resetOffset ? 0 : offset;
      const response = await SearchHistoryAPI.getSearchHistory(limit, currentOffset);

      if (resetOffset) {
        setHistory(response.items);
        setOffset(0);
      } else {
        setHistory((prev) => [...prev, ...response.items]);
      }

      setHasMore(response.has_more);
      if (!resetOffset) {
        setOffset(currentOffset + response.items.length);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load history');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isSidebarOpen) {
      loadHistory(true);
    }
  }, [isSidebarOpen]);

  const handleHistoryClick = async (item: APISearchHistoryItem) => {
    // Skip if no session_id (shouldn't happen, but be safe)
    if (!item.session_id) {
      console.warn("History item has no session_id:", item);
      setError("This search has no session data");
      return;
    }

    // Convert API history item to local history item format
    const localHistoryItem = {
      id: item.id,
      query: item.search_query,
      timestamp: new Date(item.created_at).getTime(),
      category: item.category,
      productsCount: item.result_count,
      sessionId: item.session_id,
    };

    // Use the store's loadSearchFromHistory method
    const { loadSearchFromHistory } = useChatStore.getState();
    await loadSearchFromHistory(localHistoryItem);

    // Close sidebar on mobile
    if (window.innerWidth < 1024) {
      toggleSidebar();
    }
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      // Find the item being deleted
      const deletedItem = history.find(item => item.id === id);

      await SearchHistoryAPI.deleteSearchHistory(id);
      setHistory((prev) => prev.filter((item) => item.id !== id));

      // If the deleted item matches the current session, clear messages
      if (deletedItem) {
        const { sessionId, clearMessages } = useChatStore.getState();
        if (deletedItem.session_id === sessionId) {
          clearMessages();
          // Also clear localStorage session
          localStorage.removeItem("chat_session_id");
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete');
    }
  };

  const handleClearAll = async () => {
    if (!confirm('Are you sure you want to delete all search history?')) {
      return;
    }

    try {
      await SearchHistoryAPI.deleteAllSearchHistory();
      setHistory([]);

      // Clear messages and session since all history is deleted
      const { clearMessages } = useChatStore.getState();
      clearMessages();
      localStorage.removeItem("chat_session_id");
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete all history');
    }
  };

  const getTimeAgo = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true, locale });
    } catch {
      return new Date(dateString).toLocaleDateString();
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
              <div className="flex items-center gap-1">
                <button
                  onClick={() => loadHistory(true)}
                  className="p-1.5 hover:bg-secondary rounded-md transition-colors"
                  title="Refresh"
                  disabled={loading}
                >
                  <RefreshCw className={`w-4 h-4 text-muted-foreground ${loading ? 'animate-spin' : ''}`} />
                </button>
                {history.length > 0 && (
                  <button
                    onClick={handleClearAll}
                    className="p-1.5 hover:bg-secondary rounded-md transition-colors"
                    title="Clear history"
                  >
                    <Trash2 className="w-4 h-4 text-muted-foreground" />
                  </button>
                )}
              </div>
            </div>
            <p className="text-sm text-muted-foreground">
              {history.length} {history.length === 1 ? "search" : "searches"}
            </p>
            {error && (
              <div className="mt-2 p-2 bg-destructive/10 text-destructive text-xs rounded">
                {error}
              </div>
            )}
          </div>

          {/* History List */}
          <div className="flex-1 overflow-y-auto">
            {loading && history.length === 0 ? (
              <div className="flex items-center justify-center h-full">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
              </div>
            ) : history.length === 0 ? (
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
              <>
                <div className="p-2 space-y-1">
                  {history.filter(item => item.session_id).map((item) => (
                    <button
                      key={item.id}
                      onClick={() => handleHistoryClick(item)}
                      className="w-full text-left p-3 rounded-lg hover:bg-secondary transition-colors group relative"
                    >
                      <div className="flex items-start gap-2">
                        <Search className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                        <div className="flex-1 min-w-0 pr-8">
                          <p className="text-sm font-medium truncate group-hover:text-primary transition-colors">
                            {item.search_query}
                          </p>
                          <div className="flex flex-wrap items-center gap-2 mt-1">
                            <span className="text-xs text-muted-foreground">
                              {getTimeAgo(item.created_at)}
                            </span>
                            {item.category && (
                              <>
                                <span className="text-xs text-muted-foreground">•</span>
                                <span className="text-xs text-primary/70 capitalize">
                                  {item.category}
                                </span>
                              </>
                            )}
                            {item.result_count > 0 && (
                              <>
                                <span className="text-xs text-muted-foreground">•</span>
                                <span className="text-xs text-muted-foreground flex items-center gap-1">
                                  <Package className="w-3 h-3" />
                                  {item.result_count}
                                </span>
                              </>
                            )}
                          </div>
                        </div>
                        <button
                          onClick={(e) => handleDelete(item.id, e)}
                          className="absolute right-2 top-2 opacity-0 group-hover:opacity-100 p-1 hover:bg-destructive/10 rounded transition-opacity"
                          title="Delete"
                        >
                          <Trash2 className="w-3.5 h-3.5 text-destructive" />
                        </button>
                      </div>
                    </button>
                  ))}
                </div>
                {hasMore && (
                  <div className="p-2">
                    <button
                      onClick={() => loadHistory(false)}
                      disabled={loading}
                      className="w-full py-2 text-sm text-muted-foreground hover:text-foreground hover:bg-secondary rounded-lg transition-colors disabled:opacity-50"
                    >
                      {loading ? 'Loading...' : 'Load more'}
                    </button>
                  </div>
                )}
              </>
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
