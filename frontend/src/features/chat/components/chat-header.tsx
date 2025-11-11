"use client";

import { RotateCcw, Wifi, WifiOff, Coins, PanelLeft } from "lucide-react";

import { Logo } from "@/shared/components/ui";
import { UserMenu } from "@/features/auth";
import { useChatStore } from "@/shared/lib";



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
  const { currency, country, isSidebarOpen, toggleSidebar } = useChatStore();

  return (
    <header className="border-b border-border bg-background sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center gap-3">
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
          <Logo className="h-6 md:h-8" width={84.24} height={32} />
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-1.5">
              {isConnected ? (
                <>
                  <Wifi className="w-4 h-4 text-green-500" />
                  <span className="text-xs text-green-500 font-medium">
                    {connectionStatus}
                  </span>
                </>
              ) : (
                <>
                  <WifiOff className="w-4 h-4 text-red-500" />
                  <span className="text-xs text-red-500 font-medium">
                    {connectionStatus}
                  </span>
                </>
              )}
            </div>
            {currency && (
              <div className="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-primary/10 border border-primary/20">
                <Coins className="w-3.5 h-3.5 text-primary" />
                <span className="text-xs font-semibold text-primary">
                  {currency}
                </span>
                {country && (
                  <span className="text-xs text-muted-foreground">
                    ({country})
                  </span>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={onNewSearch}
            disabled={!isConnected}
            className="px-4 py-2 rounded-full bg-secondary hover:bg-secondary/80 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <RotateCcw className="w-4 h-4" />
            <span className="hidden sm:inline">New Search</span>
          </button>
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
