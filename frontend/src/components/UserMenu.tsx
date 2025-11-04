"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/lib/auth-store";
import { AuthAPI } from "@/lib/api/auth";
import { LogOut, User as UserIcon } from "lucide-react";

export default function UserMenu() {
  const router = useRouter();
  const { user, isAuthenticated, clearAuth, refreshToken } = useAuthStore();
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

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await AuthAPI.logout(refreshToken);
      }
      clearAuth();
      setIsMenuOpen(false);
      router.push('/login');
    } catch (error) {
      console.error("Logout failed:", error);
      // Clear auth anyway on error
      clearAuth();
      router.push('/login');
    }
  };

  const getUserInitial = () => {
    if (user?.full_name) {
      return user.full_name.charAt(0).toUpperCase();
    }
    if (user?.email) {
      return user.email.charAt(0).toUpperCase();
    }
    return "?";
  };

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setIsMenuOpen(!isMenuOpen)}
        className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 transition-all overflow-hidden"
      >
        {user?.picture ? (
          <img
            src={user.picture}
            alt={user.full_name || user.email}
            className="h-full w-full object-cover"
          />
        ) : (
          <span className="text-sm font-medium">{getUserInitial()}</span>
        )}
      </button>

      {isMenuOpen && (
        <div className="fixed md:absolute right-4 md:right-0 bottom-20 md:bottom-auto md:bottom-full md:mb-2 w-64 rounded-lg bg-background border border-border shadow-xl z-50">
          <div className="p-4">
            <div className="flex items-center gap-3 mb-2">
              {user?.picture ? (
                <img
                  src={user.picture}
                  alt={user.full_name || user.email}
                  className="h-10 w-10 rounded-full"
                />
              ) : (
                <div className="h-10 w-10 rounded-full bg-primary text-primary-foreground flex items-center justify-center">
                  <UserIcon className="h-5 w-5" />
                </div>
              )}
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
              onClick={handleLogout}
              className="w-full px-4 py-3 text-left text-sm text-red-600 dark:text-red-400 hover:bg-secondary transition-colors flex items-center gap-2"
            >
              <LogOut className="h-4 w-4" />
              Sign out
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
