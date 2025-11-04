"use client";

import { useChatStore } from "@/lib/store";
import { Clock, Plus, PanelLeft, PanelLeftClose } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";
import { ThemeToggle } from "./ThemeToggle";
import { Logo } from "./Logo";
import UserMenu from "./UserMenu";

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
                <button
                  onClick={toggleSidebar}
                  className="p-2 hover:bg-secondary rounded-lg transition-colors shrink-0 group"
                  title="Collapse sidebar"
                >
                  <PanelLeftClose className="w-5 h-5 text-muted-foreground group-hover:text-foreground transition-colors" />
                </button>
                <Logo className="h-8" width={105.3} height={40} />
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
            <button
              onClick={toggleSidebar}
              className="p-2 hover:bg-secondary rounded-lg transition-colors"
              aria-label="Close sidebar"
            >
              <PanelLeftClose className="w-5 h-5 text-muted-foreground" />
            </button>
            <Logo className="h-8" width={105.3} height={40} />
          </div>

          {/* New Search Button - Right after header */}
          {isSidebarOpen && onNewSearch && (
            <div className="p-4">
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
                {/* Search History Button */}
                <div className="p-4">
                  <button
                    onClick={handleHistoryClick}
                    className={`w-full p-4 rounded-lg flex items-center gap-3 transition-colors ${
                      pathname === '/history'
                        ? 'bg-primary text-primary-foreground'
                        : 'hover:bg-secondary/50'
                    }`}
                  >
                    <Clock className="w-5 h-5" />
                    <span className="text-sm font-semibold">Search History</span>
                  </button>
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
                  onClick={handleHistoryClick}
                  className={`p-3 rounded-lg transition-colors relative group ${
                    pathname === '/history'
                      ? 'bg-primary text-primary-foreground'
                      : 'hover:bg-secondary'
                  }`}
                  title="Search History"
                >
                  <Clock className="w-5 h-5" />
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
    </>
  );
}
