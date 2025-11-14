"use client";

import { useChatStore } from "@/shared/lib";
import { Clock, Plus, PanelLeft, PanelLeftClose } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";

import { Logo } from "@/shared/components/ui";
import { UserMenu } from "@/features/auth";

interface SearchHistoryProps {
  isConnected?: boolean;
  connectionStatus?: string;
  onNewSearch?: () => void;
}

export function SearchHistory({ isConnected = true, connectionStatus = "Connected", onNewSearch }: SearchHistoryProps) {
  const {
    isSidebarOpen,
    toggleSidebar,
  } = useChatStore();

  const router = useRouter();
  const pathname = usePathname();

  const handleHistoryClick = () => {
    router.push('/history');
    if (window.innerWidth < 1024) {
      toggleSidebar();
    }
  };

  const handleNewSearchClick = () => {
    if (onNewSearch) {
      onNewSearch();
    }
    // Close sidebar on mobile
    if (window.innerWidth < 1024) {
      toggleSidebar();
    }
  };

  return (
    <>
      {/* Sidebar Panel */}
      <div
        className={`fixed left-0 top-0 bottom-0 border-r border-border transform transition-all duration-300 ease-in-out z-60 will-change-transform ${
          isSidebarOpen
            ? "w-80 translate-x-0 bg-card"
            : "w-80 -translate-x-full bg-card lg:w-16 lg:translate-x-0 lg:bg-background"
        }`}
      >
        {/* Main scrollable content area */}
        <div className="flex flex-col h-full relative overflow-hidden pb-[73px]">
          {/* Header - Desktop only */}
          <div className="to-transparent items-center justify-between gap-2 hidden lg:flex p-4">
            {/* Toggle button - always in the same position */}
            <button
              onClick={toggleSidebar}
              className="p-2 hover:bg-secondary rounded-lg transition-colors shrink-0 group cursor-pointer"
              title={isSidebarOpen ? "Collapse sidebar" : "Expand sidebar"}
            >
              {isSidebarOpen ? (
                <PanelLeftClose className="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors" />
              ) : (
                <PanelLeft className="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors" />
              )}
            </button>
            {/* Logo - fades in/out */}
            <div className={`transition-opacity duration-300 ${isSidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'}`}>
              <Logo className="h-8 w-auto" />
            </div>
          </div>

          {/* Mobile Header - Logo and close button */}
          <div className="px-4 py-2 flex items-center justify-between lg:hidden">
            <button
              onClick={toggleSidebar}
              className="p-2 rounded-lg bg-primary/10 cursor-pointer shrink-0"
              aria-label="Close sidebar"
            >
              <PanelLeftClose className="w-5 h-5" />
            </button>
            <Logo className="h-8 w-auto" />
          </div>

          {/* Main Content Area */}
          <div className="flex-1 relative overflow-hidden">
            {/* Expanded sidebar content */}
            <div className={`absolute inset-0 transition-opacity duration-300 ${isSidebarOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}>
              <div className="p-4 space-y-2">
                {/* New Search Button */}
                {onNewSearch && (
                  <button
                    onClick={handleNewSearchClick}
                    disabled={!isConnected}
                    className="w-full p-4 rounded-lg flex items-center gap-3 transition-colors bg-primary/10 hover:bg-primary/20 text-primary font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                  >
                    <Plus className="w-5 h-5" />
                    <span className="text-sm font-semibold">New Search</span>
                  </button>
                )}

                {/* Search History Button */}
                <button
                  onClick={handleHistoryClick}
                  className={`w-full p-4 rounded-lg flex items-center gap-3 transition-colors cursor-pointer ${
                    pathname === '/history'
                      ? 'bg-secondary/80 text-foreground'
                      : 'hover:bg-secondary/50 text-muted-foreground hover:text-foreground'
                  }`}
                >
                  <Clock className="w-5 h-5" />
                  <span className="text-sm font-semibold">Search History</span>
                </button>
              </div>
            </div>

            {/* Collapsed sidebar - only icons with tooltips (desktop only) */}
            <div className={`absolute inset-0 hidden lg:flex flex-col items-center gap-2 py-4 transition-opacity duration-300 ${!isSidebarOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'}`}>
              {/* New Search Icon */}
              {onNewSearch && (
                <button
                  onClick={handleNewSearchClick}
                  disabled={!isConnected}
                  className="p-3 rounded-lg bg-primary/10 hover:bg-primary/20 text-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer relative group"
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
                onClick={handleHistoryClick}
                className={`p-3 rounded-lg transition-colors relative group cursor-pointer ${
                  pathname === '/history'
                    ? 'bg-secondary/80 text-foreground'
                    : 'hover:bg-secondary/50 text-muted-foreground hover:text-foreground'
                }`}
                title="Search History"
              >
                <Clock className="w-5 h-5" />
                <div className="absolute left-full ml-2 top-1/2 -translate-y-1/2 px-2 py-1 bg-popover text-popover-foreground text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50">
                  Search History
                </div>
              </button>
            </div>
          </div>

        </div>

        {/* Bottom Controls - outside overflow context */}
        <div className="absolute bottom-0 left-0 right-0">
          {/* Expanded - User Menu with name */}
          <div className={`p-2 flex items-start justify-start transition-opacity duration-300 ${isSidebarOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none absolute inset-0'}`}>
            <UserMenu showName={true} />
          </div>
          {/* Collapsed - Icons only (desktop only) */}
          <div className={`hidden lg:flex flex-col items-start pl-4 py-4 transition-opacity duration-300 ${!isSidebarOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none absolute inset-0'}`}>
            <UserMenu showName={false} />
          </div>
        </div>
      </div>

      {/* Overlay for mobile - only show when sidebar is open */}
      <div
        className={`fixed inset-0 bg-black/50 z-50 lg:hidden transition-opacity duration-300 ${
          isSidebarOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
        }`}
        onClick={toggleSidebar}
      />
    </>
  );
}
