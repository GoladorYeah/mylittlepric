"use client";

import { RotateCcw, Wifi, WifiOff } from "lucide-react";
import { ThemeToggle } from "../ThemeToggle";
import { Logo } from "../Logo";
import UserMenu from "../UserMenu";

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
  return (
    <header className="border-b border-border bg-background sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Logo className="h-6 md:h-8" width={84.24} height={32} />
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
          <ThemeToggle />
          <UserMenu />
        </div>
      </div>
    </header>
  );
}
