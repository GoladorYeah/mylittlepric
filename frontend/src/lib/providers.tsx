"use client";

import { ThemeProvider as NextThemesProvider } from "next-themes";

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
      {children}
    </NextThemesProvider>
  );
}