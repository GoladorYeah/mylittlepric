"use client";

import { useEffect } from "react";
import { ThemeProvider as NextThemesProvider } from "next-themes";
import { useAuthStore, useChatStore } from "@/shared/lib";

function PreferencesSync() {
  const { isAuthenticated, _hasHydrated } = useAuthStore();
  const syncPreferencesFromServer = useChatStore((state) => state.syncPreferencesFromServer);

  useEffect(() => {
    // Only sync if user is authenticated and hydration is complete
    if (_hasHydrated && isAuthenticated) {
      console.log("🔄 Syncing preferences from server on app load...");
      syncPreferencesFromServer();
    }
  }, [isAuthenticated, _hasHydrated, syncPreferencesFromServer]);

  return null;
}

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme="system"
      enableSystem={true}
      enableColorScheme={true}
      storageKey="mylittleprice-theme"
      disableTransitionOnChange={false}
    >
      <PreferencesSync />
      {children}
    </NextThemesProvider>
  );
}