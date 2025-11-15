"use client";

import { ReactNode } from "react";

interface PolicyLayoutProps {
  title: string;
  lastUpdated: string;
  children: ReactNode;
}

export function PolicyLayout({ title, lastUpdated, children }: PolicyLayoutProps) {
  return (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto px-4 py-24">
        <div className="max-w-4xl mx-auto space-y-8">
          <div className="text-center space-y-4">
            <h1 className="text-5xl md:text-6xl font-bold bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              {title}
            </h1>
            <p className="text-lg text-muted-foreground">
              Last updated: {lastUpdated}
            </p>
          </div>

          <div className="prose prose-neutral dark:prose-invert max-w-none space-y-8">
            {children}
          </div>
        </div>
      </main>
    </div>
  );
}
