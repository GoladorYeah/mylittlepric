"use client";

import { useEffect, useState } from "react";
import { useChatStore } from "@/lib/store";
import { useAuthStore } from "@/lib/auth-store";
import { SearchHistoryAPI, type SearchHistoryItem as APISearchHistoryItem } from "@/lib/search-history-api";
import { Clock, Trash2, Search, Package, RefreshCw, LogIn, Plus, ChevronDown, ChevronRight, PanelLeft, PanelLeftClose } from "lucide-react";
import { formatDistanceToNow } from "date-fns";
import { ru, uk, enUS } from "date-fns/locale";
import AuthDialog from "./AuthDialog";
import { ThemeToggle } from "./ThemeToggle";
import { Logo } from "./Logo";
import UserMenu from "./UserMenu";

const localeMap: Record<string, any> = {
  ru,
  uk,
  en: enUS,
};

interface SearchHistoryProps {
  isConnected?: boolean;
  connectionStatus?: string;
  onNewSearch?: () => void;
}

export function SearchHistory({ isConnected = true, connectionStatus = "Connected", onNewSearch }: SearchHistoryProps) {
  const {
    isSidebarOpen,
    toggleSidebar,
    language,
  } = useChatStore();

  const { isAuthenticated } = useAuthStore();

  const [history, setHistory] = useState<APISearchHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const [isAuthDialogOpen, setIsAuthDialogOpen] = useState(false);
  const [isHistoryExpanded, setIsHistoryExpanded] = useState(false);
  const limit = 20;

  const locale = localeMap[language] || enUS;

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
    if (isSidebarOpen && isHistoryExpanded) {
      loadHistory(true);
    }
  }, [isHistoryExpanded, isAuthenticated]);

  const handleHistoryClick = async (item: APISearchHistoryItem) => {
    if (!item.session_id) {
      console.warn("History item has no session_id:", item);
      setError("This search has no session data");
      return;
    }

    const localHistoryItem = {
      id: item.id,
      query: item.search_query,
      timestamp: new Date(item.created_at).getTime(),
      category: item.category,
      productsCount: item.result_count,
      sessionId: item.session_id,
    };

    const { loadSearchFromHistory } = useChatStore.getState();
    await loadSearchFromHistory(localHistoryItem);

    if (window.innerWidth < 1024) {
      toggleSidebar();
    }
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

  return (
    <>
      {/* Mobile Toggle Button - Only when sidebar closed on mobile */}
      {!isSidebarOpen && (
        <button
          onClick={toggleSidebar}
          className="fixed left-4 top-4 z-50 bg-primary text-primary-foreground p-2 rounded-lg shadow-lg hover:opacity-90 transition-all lg:hidden"
          aria-label="Open sidebar"
        >
          <PanelLeft className="w-5 h-5" />
        </button>
      )}

      {/* Sidebar Panel */}
      <div
        className={`fixed left-0 top-0 bottom-0 backdrop-blur-xl border-r border-border shadow-2xl transform transition-all duration-300 ease-in-out z-40 overflow-hidden ${
          isSidebarOpen
            ? "w-80 translate-x-0 bg-card/95"
            : "w-16 -translate-x-full lg:translate-x-0 lg:bg-background"
        }`}
      >
        <div className="flex flex-col h-full overflow-hidden">
          {/* Header - Desktop only */}
          <div className="border-b border-border bg-linear-to-b from-background/50 to-transparent items-center justify-between gap-2 hidden lg:flex p-4">
            {isSidebarOpen ? (
              <>
                <Logo className="h-8" width={105.3} height={40} />
                <button
                  onClick={toggleSidebar}
                  className="p-2 hover:bg-secondary rounded-lg transition-colors shrink-0 group"
                  title="Collapse sidebar"
                >
                  <PanelLeftClose className="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors" />
                </button>
              </>
            ) : (
              <button
                onClick={toggleSidebar}
                className="p-2 hover:bg-secondary rounded-lg transition-colors mx-auto group"
                title="Expand sidebar"
              >
                <PanelLeft className="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors" />
              </button>
            )}
          </div>

          {/* Mobile Header - Logo and close button */}
          <div className="p-4 border-b border-border flex items-center justify-between lg:hidden">
            <Logo className="h-8" width={105.3} height={40} />
            <button
              onClick={toggleSidebar}
              className="p-2 hover:bg-secondary rounded-lg transition-colors"
              aria-label="Close sidebar"
            >
              <PanelLeftClose className="w-5 h-5 text-muted-foreground" />
            </button>
          </div>

          {/* New Search Button - Right after header */}
          {isSidebarOpen && onNewSearch && (
            <div className="p-4 border-b border-border">
              <button
                onClick={onNewSearch}
                disabled={!isConnected}
                className="w-full px-4 py-3 rounded-lg bg-primary hover:bg-primary/90 text-primary-foreground font-medium transition-all duration-300 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed hover:shadow-lg hover:shadow-primary/30"
              >
                <Plus className="w-5 h-5" />
                <span>New Search</span>
              </button>
            </div>
          )}

          {/* Main Content Area */}
          <div className="flex-1">
            {isSidebarOpen ? (
              // Expanded sidebar content
              <>
                {/* Search History Submenu */}
                <div className="border-b border-border">
                  <button
                    onClick={() => setIsHistoryExpanded(!isHistoryExpanded)}
                    className="w-full p-4 flex items-center justify-between hover:bg-secondary/50 transition-colors"
                  >
                    <div className="flex items-center gap-2">
                      <Clock className="w-5 h-5 text-primary" />
                      <span className="text-sm font-semibold">Search History</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs text-muted-foreground">
                        {history.length}
                      </span>
                      {isHistoryExpanded ? (
                        <ChevronDown className="w-4 h-4 text-muted-foreground" />
                      ) : (
                        <ChevronRight className="w-4 h-4 text-muted-foreground" />
                      )}
                    </div>
                  </button>

                  {/* History Content (when expanded) */}
                  {isHistoryExpanded && (
                    <div className="border-t border-border">
                      {/* History Controls */}
                      <div className="p-3 border-b border-border/50 flex items-center justify-between bg-secondary/20">
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
                              title="Clear all"
                            >
                              <Trash2 className="w-4 h-4 text-muted-foreground" />
                            </button>
                          )}
                        </div>
                      </div>

                      {error && (
                        <div className="mx-2 mt-2 p-2 bg-destructive/10 text-destructive text-xs rounded">
                          {error}
                        </div>
                      )}

                      {/* History List */}
                      <div className="max-h-96 overflow-y-auto">
                        {!isAuthenticated ? (
                          <div className="p-6 text-center">
                            <LogIn className="w-10 h-10 text-muted-foreground/50 mb-2 mx-auto" />
                            <p className="text-sm font-medium text-foreground mb-1">
                              Sign in to save history
                            </p>
                            <p className="text-xs text-muted-foreground mb-4">
                              Access your searches across devices
                            </p>
                            <button
                              onClick={() => setIsAuthDialogOpen(true)}
                              className="rounded-lg bg-primary px-4 py-2 text-xs font-medium text-primary-foreground hover:opacity-90 transition-all"
                            >
                              Sign in
                            </button>
                          </div>
                        ) : loading && history.length === 0 ? (
                          <div className="flex items-center justify-center py-8">
                            <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
                          </div>
                        ) : history.length === 0 ? (
                          <div className="p-6 text-center">
                            <Search className="w-10 h-10 text-muted-foreground/50 mb-2 mx-auto" />
                            <p className="text-xs text-muted-foreground">
                              No search history yet
                            </p>
                          </div>
                        ) : (
                          <>
                            <div className="p-2 space-y-1">
                              {history.filter(item => item.session_id).map((item) => (
                                <button
                                  key={item.id}
                                  onClick={() => handleHistoryClick(item)}
                                  className="w-full text-left p-2 rounded-lg hover:bg-secondary transition-colors group relative"
                                >
                                  <div className="flex items-start gap-2">
                                    <Search className="w-3.5 h-3.5 text-muted-foreground mt-0.5 shrink-0" />
                                    <div className="flex-1 min-w-0 pr-6">
                                      <p className="text-xs font-medium truncate group-hover:text-primary transition-colors">
                                        {item.search_query}
                                      </p>
                                      <div className="flex flex-wrap items-center gap-1.5 mt-0.5">
                                        <span className="text-[10px] text-muted-foreground">
                                          {getTimeAgo(item.created_at)}
                                        </span>
                                        {item.result_count > 0 && (
                                          <>
                                            <span className="text-[10px] text-muted-foreground">•</span>
                                            <span className="text-[10px] text-muted-foreground flex items-center gap-0.5">
                                              <Package className="w-2.5 h-2.5" />
                                              {item.result_count}
                                            </span>
                                          </>
                                        )}
                                      </div>
                                    </div>
                                    <button
                                      onClick={(e) => handleDelete(item.id, e)}
                                      className="absolute right-1 top-1 opacity-0 group-hover:opacity-100 p-1 hover:bg-destructive/10 rounded transition-opacity"
                                      title="Delete"
                                    >
                                      <Trash2 className="w-3 h-3 text-destructive" />
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
                                  className="w-full py-1.5 text-xs text-muted-foreground hover:text-foreground hover:bg-secondary rounded-lg transition-colors disabled:opacity-50"
                                >
                                  {loading ? 'Loading...' : 'Load more'}
                                </button>
                              </div>
                            )}
                          </>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </>
            ) : (
              // Collapsed sidebar - only icons with tooltips (desktop only)
              <div className="hidden lg:flex flex-col items-center gap-2 py-4">
                {/* New Search Icon */}
                {onNewSearch && (
                  <button
                    onClick={onNewSearch}
                    disabled={!isConnected}
                    className="p-3 rounded-lg bg-primary hover:bg-primary/90 text-primary-foreground transition-colors disabled:opacity-50 disabled:cursor-not-allowed relative group"
                    title="New Search"
                  >
                    <Plus className="w-5 h-5" />
                    <div className="absolute left-full ml-2 top-1/2 -translate-y-1/2 px-2 py-1 bg-popover text-popover-foreground text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50">
                      New Search
                    </div>
                  </button>
                )}

                {/* History Icon */}
                <button
                  onClick={() => {
                    toggleSidebar();
                    setIsHistoryExpanded(true);
                  }}
                  className="p-3 hover:bg-secondary rounded-lg transition-colors relative group"
                  title="Search History"
                >
                  <Clock className="w-5 h-5 text-muted-foreground" />
                  {history.length > 0 && (
                    <span className="absolute -top-1 -right-1 bg-primary text-primary-foreground text-[10px] font-bold rounded-full w-4 h-4 flex items-center justify-center">
                      {history.length > 9 ? '9+' : history.length}
                    </span>
                  )}
                  <div className="absolute left-full ml-2 top-1/2 -translate-y-1/2 px-2 py-1 bg-popover text-popover-foreground text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50">
                    Search History
                  </div>
                </button>
              </div>
            )}
          </div>

          {/* Bottom Controls */}
          <div className="mt-auto border-t border-border bg-linear-to-t from-background/50 to-transparent">
            {isSidebarOpen ? (
              // Expanded - Theme & User
              <div className="p-4 flex items-center justify-between">
                <ThemeToggle />
                <UserMenu />
              </div>
            ) : (
              // Collapsed - Icons only (desktop only)
              <div className="hidden lg:flex flex-col items-center gap-4 py-4">
                <div className="relative group">
                  <ThemeToggle />
                </div>
                <UserMenu />
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Overlay for mobile - only show when sidebar is open */}
      {isSidebarOpen && (
        <div
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-30 lg:hidden transition-opacity duration-300"
          onClick={toggleSidebar}
        />
      )}

      {/* Auth Dialog */}
      <AuthDialog
        isOpen={isAuthDialogOpen}
        onClose={() => setIsAuthDialogOpen(false)}
      />
    </>
  );
}
