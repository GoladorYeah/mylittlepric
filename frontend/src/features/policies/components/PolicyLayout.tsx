"use client";

import Link from "next/link";
import { Home, MessageCircle } from "lucide-react";
import { Logo } from "@/shared/components/ui";
import { ReactNode } from "react";

interface PolicyLayoutProps {
  title: string;
  lastUpdated: string;
  children: ReactNode;
}

export function PolicyLayout({ title, lastUpdated, children }: PolicyLayoutProps) {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
      <header className="fixed top-0 left-0 right-0 z-50 bg-background border-b border-border">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <Link
                href="/"
                className="flex items-center gap-2 px-3 py-2 text-muted-foreground hover:text-foreground hover:bg-muted rounded-lg transition-colors"
              >
                <Home className="w-5 h-5" />
                <span className="hidden sm:inline">Home</span>
              </Link>
              <Link
                href="/chat"
                className="flex items-center gap-2 px-3 py-2 text-muted-foreground hover:text-foreground hover:bg-muted rounded-lg transition-colors"
              >
                <MessageCircle className="w-5 h-5" />
                <span className="hidden sm:inline">Chat</span>
              </Link>
            </div>
            <div className="h-6 w-px bg-border hidden sm:block" />
            <Logo width={84.24} height={32} />
          </div>
        </div>
      </header>

      <main className="pt-24 pb-12">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto">
            <div className="text-center space-y-2 mb-8">
              <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                {title}
              </h1>
              <p className="text-sm text-muted-foreground">
                Last updated: {lastUpdated}
              </p>
            </div>

            <div className="bg-background/50 border border-border rounded-2xl shadow-xl p-6 md:p-10">
              <div className="prose prose-neutral dark:prose-invert max-w-none">
                {children}
              </div>
            </div>
          </div>
        </div>
      </main>

      <footer className="border-t border-border py-8 mt-12">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-muted-foreground">
            <p>&copy; 2025 MyLittlePrice. All rights reserved.</p>
            <div className="flex gap-6">
              <a href="/privacy-policy" className="hover:text-foreground transition-colors">
                Privacy Policy
              </a>
              <a href="/terms-of-use" className="hover:text-foreground transition-colors">
                Terms of Use
              </a>
              <a href="/cookie-policy" className="hover:text-foreground transition-colors">
                Cookie Policy
              </a>
              <a href="/advertising-policy" className="hover:text-foreground transition-colors">
                Advertising Policy
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
