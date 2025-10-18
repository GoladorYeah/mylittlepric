"use client";

import { useState, useRef, useEffect } from "react";
import { useAuthStore } from "@/lib/auth-store";
import { authAPI } from "@/lib/auth-api";
import AuthDialog from "./AuthDialog";

export default function UserMenu() {
  const { user, isAuthenticated, clearAuth } = useAuthStore();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isAuthDialogOpen, setIsAuthDialogOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = async () => {
    try {
      await authAPI.logout();
      setIsMenuOpen(false);
    } catch (error) {
      console.error("Logout failed:", error);
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
    return (
      <>
        <button
          onClick={() => setIsAuthDialogOpen(true)}
          className="rounded-full bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 transition-all"
        >
          Sign in
        </button>

        <AuthDialog
          isOpen={isAuthDialogOpen}
          onClose={() => setIsAuthDialogOpen(false)}
        />
      </>
    );
  }

  return (
    <div className="relative" ref={menuRef}>
      <button
        onClick={() => setIsMenuOpen(!isMenuOpen)}
        className="flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 transition-all"
      >
        {getUserInitial()}
      </button>

      {isMenuOpen && (
        <div className="absolute right-0 mt-2 w-64 rounded-lg bg-background border border-border shadow-xl">
          <div className="p-4">
            <p className="text-sm font-medium text-foreground">
              {user?.full_name || "User"}
            </p>
            <p className="text-sm text-muted-foreground truncate">
              {user?.email}
            </p>
          </div>

          <div className="border-t border-border">
            <button
              onClick={handleLogout}
              className="w-full px-4 py-3 text-left text-sm text-red-600 dark:text-red-400 hover:bg-secondary transition-colors"
            >
              Sign out
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
