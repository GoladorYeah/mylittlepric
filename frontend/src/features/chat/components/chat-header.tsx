"use client";

import { useEffect, useState } from "react";
import { PanelLeft } from "lucide-react";
import { useChatStore } from "@/shared/lib";
import { RateLimitIndicator } from "./RateLimitIndicator";

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
  const { isSidebarOpen, toggleSidebar, isLoading, rateLimitState } = useChatStore();
  const [isSyncing, setIsSyncing] = useState(false);

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
      <div className="container mx-auto px-4 h-12 flex items-center justify-between">
        {/* Mobile Sidebar Toggle Button */}
        <button
          onClick={toggleSidebar}
          className={`p-2 rounded-lg bg-primary text-primary-foreground hover:opacity-90 transition-opacity lg:hidden ${
            isSidebarOpen ? 'opacity-0 pointer-events-none' : 'opacity-100'
          }`}
          aria-label="Open sidebar"
        >
          <PanelLeft className="w-5 h-5" />
        </button>

        {/* Connection Status and Rate Limit Indicator */}
        <div className="flex items-center gap-3 ml-auto">
          {/* Rate Limit Indicator */}
          <RateLimitIndicator />

          {/* Connection Status */}
          <div className="flex items-center gap-2">
            <div className={`w-2 h-2 rounded-full ${getStatusColor()} ${isConnected ? 'animate-pulse' : ''}`} />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {getStatusText()}
            </span>
          </div>
        </div>
      </div>
    </header>
  );
}
