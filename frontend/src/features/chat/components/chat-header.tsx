"use client";

import { PanelLeft } from "lucide-react";
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
  const { isSidebarOpen, toggleSidebar } = useChatStore();

  return (
    <header className="border-b border-border bg-background sticky top-0 z-20">
      <div className="container mx-auto px-4 h-12 flex items-center">
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
      </div>
    </header>
  );
}
