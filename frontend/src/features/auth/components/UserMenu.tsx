"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/shared/lib";
import { useChatStore } from "@/shared/lib";
import { AuthAPI } from "@/shared/lib/api/auth";
import { LogOut, Settings } from "lucide-react";

interface UserMenuProps {
  showName?: boolean;
}

export default function UserMenu({ showName = false }: UserMenuProps) {
  const router = useRouter();
  const { user, isAuthenticated, clearAuth, refreshToken } = useAuthStore();
  const { isSidebarOpen } = useChatStore();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close menu when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    }

    if (isMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
  }, [isMenuOpen]);

  const handleSettings = () => {
    setIsMenuOpen(false);
    router.push("/settings");
  };

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await AuthAPI.logout(refreshToken);
      }
      clearAuth();
      setIsMenuOpen(false);
    } catch (error) {
      console.error("Logout failed:", error);
      clearAuth();
    }
  };

  const getUserInitials = () => {
    if (user?.full_name) {
      const names = user.full_name.trim().split(' ');
      if (names.length >= 2) {
        return (names[0].charAt(0) + names[names.length - 1].charAt(0)).toUpperCase();
      }
      return user.full_name.substring(0, 2).toUpperCase();
    }
    if (user?.email) {
      return user.email.substring(0, 2).toUpperCase();
    }
    return "??";
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <>
      <div className="relative w-full" ref={menuRef}>
        <button
          onClick={() => setIsMenuOpen(!isMenuOpen)}
          className={`flex items-center gap-3 hover:opacity-90 transition-all cursor-pointer ${
            showName
              ? 'w-full p-2 rounded-lg hover:bg-secondary/50'
              : 'h-10 w-10 justify-center rounded-full bg-primary text-primary-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2'
          }`}
        >
          <div className={`flex items-center justify-center rounded-full bg-primary text-primary-foreground shrink-0 ${
            showName ? 'h-9 w-9' : 'h-10 w-10'
          }`}>
            <span className="text-sm font-semibold">{getUserInitials()}</span>
          </div>
          {showName && (
            <div className="flex-1 min-w-0 text-left">
              <p className="text-sm font-medium text-foreground truncate">
                {user?.full_name || "User"}
              </p>
              <p className="text-xs text-muted-foreground truncate">
                {user?.email}
              </p>
            </div>
          )}
        </button>

        {isMenuOpen && (
          <div className={`fixed w-64 rounded-lg bg-background border border-border shadow-xl z-[60] ${
            isSidebarOpen
              ? 'bottom-20 md:bottom-auto right-4 md:right-4 md:bottom-full md:mb-2'
              : 'bottom-20 md:bottom-4 left-4 md:left-4'
          }`}>
            <div className="p-4">
              <div className="flex items-center gap-3 mb-2">
                <div className="h-10 w-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center">
                  <span className="text-sm font-semibold">{getUserInitials()}</span>
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-foreground truncate">
                    {user?.full_name || "User"}
                  </p>
                  <p className="text-xs text-muted-foreground truncate">
                    {user?.email}
                  </p>
                </div>
              </div>
              {user?.provider && (
                <div className="mt-2 px-2 py-1 rounded-md bg-secondary/50 text-xs text-muted-foreground">
                  Signed in with {user.provider === 'google' ? 'Google' : user.provider}
                </div>
              )}
            </div>

            <div className="border-t border-border">
              <button
                onClick={handleSettings}
                className="w-full px-4 py-3 text-left text-sm text-foreground hover:bg-secondary transition-colors flex items-center gap-2 cursor-pointer"
              >
                <Settings className="h-4 w-4" />
                Settings
              </button>
              <button
                onClick={handleLogout}
                className="w-full px-4 py-3 text-left text-sm text-red-600 dark:text-red-400 hover:bg-secondary transition-colors flex items-center gap-2 cursor-pointer"
              >
                <LogOut className="h-4 w-4" />
                Sign out
              </button>
            </div>
          </div>
        )}
      </div>
    </>
  );
}
