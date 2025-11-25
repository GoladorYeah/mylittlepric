"use client";

import { useEffect, useState } from "react";
import { PanelLeft, LogIn } from "lucide-react";
import { useChatStore, useAuthStore } from "@/shared/lib";
import { RateLimitIndicator } from "./RateLimitIndicator";
import { LastSearchSaved } from "./LastSearchSaved";
import { BugReportButton } from "@/features/bug-report";
import { AuthDialog } from "@/features/auth/components";

interface ChatHeaderProps {
  isConnected: boolean;
  connectionStatus: string;
  onNewSearch: () => void;
}

export function ChatHeader({
  isConnected,
  connectionStatus,
  onNewSearch,
}: ChatHeaderProps) {
  const { isSidebarOpen, toggleSidebar, isLoading, rateLimitState, searchState } = useChatStore();
  const { isAuthenticated } = useAuthStore();
  const [isSyncing, setIsSyncing] = useState(false);
  const [showAuthDialog, setShowAuthDialog] = useState(false);

  // Detect syncing state (when connected but loading after reconnect)
  useEffect(() => {
    if (isConnected && isLoading) {
      // Could be syncing missed messages
      setIsSyncing(true);
      const timer = setTimeout(() => setIsSyncing(false), 5000); // Clear after 5s
      return () => clearTimeout(timer);
    } else {
      setIsSyncing(false);
    }
  }, [isConnected, isLoading]);

  const getStatusColor = () => {
    if (rateLimitState.isLimited) return "bg-yellow-500";
    if (isConnected) return "bg-green-500";
    if (connectionStatus === "Connecting") return "bg-yellow-500";
    return "bg-red-500";
  };

  const getStatusText = () => {
    if (rateLimitState.isLimited && rateLimitState.expiresAt) {
      const remaining = Math.ceil((rateLimitState.expiresAt.getTime() - Date.now()) / 1000);
      return `Rate limited (${remaining}s)`;
    }
    if (isSyncing) {
      return "Syncing...";
    }
    return connectionStatus;
  };

  return (
    <header className="bg-background sticky top-0 z-30">
      <div className="from-background via-background via-65% to-background-100/0 pointer-events-none absolute inset-0 -bottom-5 -z-1 bg-linear-to-b blur-sm"></div>
      <div className="mx-auto px-4 py-2 flex items-center justify-between gap-3">
        {/* Left section: Mobile toggle + Last Search Saved */}
        <div className="flex items-center gap-3">
          {/* Mobile Sidebar Toggle Button */}
          <button
            onClick={toggleSidebar}
            className={`p-2 rounded-lg bg-primary/10 transition-opacity lg:hidden shrink-0 ${
              isSidebarOpen ? 'opacity-0 pointer-events-none' : 'opacity-100'
            }`}
            aria-label="Open sidebar"
          >
            <PanelLeft className="w-5 h-5" />
          </button>

          {/* Last Search Saved - Desktop and Mobile */}
          <LastSearchSaved />
        </div>

        {/* Connection Status and Rate Limit Indicator */}
        <div className="flex items-center gap-3 shrink-0">
          {/* Anonymous User - Show Search Count and Sign In */}
          {!isAuthenticated && (
            <div className="flex items-center gap-2">
              {searchState && (
                <span className="text-sm text-muted-foreground">
                  {searchState.anonymous_search_used}/{searchState.anonymous_search_limit} free searches
                </span>
              )}
              <button
                onClick={() => setShowAuthDialog(true)}
                className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
              >
                <LogIn className="w-4 h-4" />
                <span className="hidden sm:inline">Sign In</span>
              </button>
            </div>
          )}

          {/* Rate Limit Indicator */}
          <RateLimitIndicator />

          {/* Connection Status */}
          <div className="hidden sm:flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${getStatusColor()} ${isConnected ? 'animate-pulse' : ''}`} />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {getStatusText()}
            </span>
          </div>

          {/* Bug Report Button */}
          <BugReportButton variant="header" />
        </div>
      </div>

      {/* Auth Dialog */}
      <AuthDialog isOpen={showAuthDialog} onClose={() => setShowAuthDialog(false)} />
    </header>
  );
}
