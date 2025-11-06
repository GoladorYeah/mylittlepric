"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/auth-store";
import { useChatStore } from "@/lib/store";
import { SearchHistoryAPI, type SearchHistoryItem as APISearchHistoryItem } from "@/lib/search-history-api";
import { Clock, Trash2, Search, Package, RefreshCw, LogIn, ChevronDown, ChevronUp } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { ru, uk, enUS } from "date-fns/locale";
import { SearchHistory as Sidebar } from "@/components/SearchHistory";
import { ProductCard } from "@/components/ProductCard";

const localeMap: Record<string, any> = {
  ru,
  uk,
  en: enUS,
};

export default function HistoryPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuthStore();
  const { language, isSidebarOpen } = useChatStore();

  const [history, setHistory] = useState<APISearchHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set());
  const limit = 50;

  const locale = localeMap[language] || enUS;

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, authLoading, router]);

  const loadHistory = async (resetOffset = false) => {
    if (!isAuthenticated) {
      setHistory([]);
      setLoading(false);
      return;
    }

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
    if (isAuthenticated) {
      loadHistory(true);
    }
  }, [isAuthenticated]);

  const toggleExpanded = (id: string) => {
    setExpandedItems((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      const deletedItem = history.find(item => item.id === id);
      await SearchHistoryAPI.deleteSearchHistory(id);
      setHistory((prev) => prev.filter((item) => item.id !== id));

      if (deletedItem) {
        const { sessionId, clearMessages } = useChatStore.getState();
        if (deletedItem.session_id === sessionId) {
          clearMessages();
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

  if (authLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <>
      <Sidebar onNewSearch={() => router.push('/chat')} />

      <div
        className={`min-h-screen bg-gradient-to-br from-background via-background to-background/95 transition-all duration-300 ${
          isSidebarOpen ? 'lg:pl-80' : 'lg:pl-16'
        }`}
      >
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <Clock className="w-8 h-8 text-primary" />
                <h1 className="text-3xl font-bold">Search History</h1>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => loadHistory(true)}
                  className="p-2 hover:bg-secondary rounded-lg transition-colors"
                  title="Refresh"
                  disabled={loading}
                >
                  <RefreshCw className={`w-5 h-5 text-muted-foreground ${loading ? 'animate-spin' : ''}`} />
                </button>
                {history.length > 0 && (
                  <button
                    onClick={handleClearAll}
                    className="px-4 py-2 bg-destructive/10 text-destructive hover:bg-destructive/20 rounded-lg transition-colors flex items-center gap-2"
                    title="Clear all"
                  >
                    <Trash2 className="w-4 h-4" />
                    <span className="text-sm font-medium">Clear All</span>
                  </button>
                )}
              </div>
            </div>
            <p className="text-muted-foreground">
              {history.length > 0
                ? `${history.length} search${history.length !== 1 ? 'es' : ''} found`
                : 'No search history yet'}
            </p>
          </div>

          {error && (
            <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-lg">
              {error}
            </div>
          )}

          {/* History List */}
          {!isAuthenticated ? (
            <div className="text-center py-16">
              <LogIn className="w-16 h-16 text-muted-foreground/50 mb-4 mx-auto" />
              <h2 className="text-xl font-semibold mb-2">Sign in to view history</h2>
              <p className="text-muted-foreground mb-6">
                Access your search history across all devices
              </p>
              <button
                onClick={() => router.push('/login')}
                className="px-6 py-3 bg-primary text-primary-foreground rounded-lg hover:opacity-90 transition-all"
              >
                Sign in
              </button>
            </div>
          ) : loading && history.length === 0 ? (
            <div className="flex items-center justify-center py-16">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
            </div>
          ) : history.length === 0 ? (
            <div className="text-center py-16">
              <Search className="w-16 h-16 text-muted-foreground/50 mb-4 mx-auto" />
              <h2 className="text-xl font-semibold mb-2">No search history yet</h2>
              <p className="text-muted-foreground mb-6">
                Start searching to build your history
              </p>
              <button
                onClick={() => router.push('/chat')}
                className="px-6 py-3 bg-primary text-primary-foreground rounded-lg hover:opacity-90 transition-all"
              >
                Start Searching
              </button>
            </div>
          ) : (
            <>
              <div className="space-y-4">
                {history.filter(item => item.session_id).map((item) => {
                  const isExpanded = expandedItems.has(item.id);
                  const hasProducts = item.products_found && item.products_found.length > 0;

                  return (
                    <div
                      key={item.id}
                      className="rounded-lg bg-card border border-border overflow-hidden"
                    >
                      {/* Header */}
                      <div className="relative group">
                        <button
                          onClick={() => hasProducts && toggleExpanded(item.id)}
                          className={`w-full text-left p-4 transition-colors ${
                            hasProducts ? 'hover:bg-secondary cursor-pointer' : 'cursor-default'
                          }`}
                        >
                          <div className="flex items-start gap-3">
                            <Search className="w-5 h-5 text-muted-foreground mt-0.5 shrink-0" />
                            <div className="flex-1 min-w-0 pr-16">
                              <p className="text-base font-medium mb-1 group-hover:text-primary transition-colors">
                                {item.search_query}
                              </p>
                              <div className="flex flex-wrap items-center gap-2">
                                <span className="text-sm text-muted-foreground">
                                  {getTimeAgo(item.created_at)}
                                </span>
                                {item.result_count > 0 && (
                                  <>
                                    <span className="text-sm text-muted-foreground">•</span>
                                    <span className="text-sm text-muted-foreground flex items-center gap-1">
                                      <Package className="w-4 h-4" />
                                      {item.result_count} product{item.result_count !== 1 ? 's' : ''}
                                    </span>
                                  </>
                                )}
                                {item.category && (
                                  <>
                                    <span className="text-sm text-muted-foreground">•</span>
                                    <span className="text-sm text-muted-foreground capitalize">
                                      {item.category}
                                    </span>
                                  </>
                                )}
                              </div>
                            </div>
                            <div className="absolute right-2 top-2 flex items-center gap-1">
                              {hasProducts && (
                                <button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    toggleExpanded(item.id);
                                  }}
                                  className="p-2 hover:bg-secondary/50 rounded-lg transition-all"
                                  title={isExpanded ? "Hide products" : "Show products"}
                                >
                                  {isExpanded ? (
                                    <ChevronUp className="w-4 h-4 text-muted-foreground" />
                                  ) : (
                                    <ChevronDown className="w-4 h-4 text-muted-foreground" />
                                  )}
                                </button>
                              )}
                              <button
                                onClick={(e) => handleDelete(item.id, e)}
                                className="opacity-0 group-hover:opacity-100 p-2 hover:bg-destructive/10 rounded-lg transition-all"
                                title="Delete"
                              >
                                <Trash2 className="w-4 h-4 text-destructive" />
                              </button>
                            </div>
                          </div>
                        </button>
                      </div>

                      {/* Products Grid */}
                      {isExpanded && hasProducts && (
                        <div className="border-t border-border bg-muted/20 p-4">
                          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                            {item.products_found!.map((product, idx) => (
                              <div key={`${item.id}-product-${idx}`} className="w-full">
                                <ProductCard product={product} index={idx + 1} />
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>

              {hasMore && (
                <div className="mt-6 text-center">
                  <button
                    onClick={() => loadHistory(false)}
                    disabled={loading}
                    className="px-6 py-3 bg-secondary hover:bg-secondary/80 text-foreground rounded-lg transition-colors disabled:opacity-50"
                  >
                    {loading ? 'Loading...' : 'Load More'}
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </>
  );
}
