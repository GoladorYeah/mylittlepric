"use client";

import { useState } from "react";

interface AdPlaceholderProps {
  /**
   * Ad format type
   * - "card": 300x250 (Medium Rectangle) - fits in product slider
   * - "banner": 728x90 (Leaderboard) - horizontal banner
   * - "wide-banner": 970x90 (Large Leaderboard) - wider horizontal
   * - "skyscraper": 300x600 (Half Page) - vertical sidebar
   */
  format?: "card" | "banner" | "wide-banner" | "skyscraper";
  className?: string;
}

export function AdPlaceholder({ format = "card", className = "" }: AdPlaceholderProps) {
  const [isVisible, setIsVisible] = useState(true);

  if (!isVisible) return null;

  const dimensions = {
    card: { width: 300, height: 250, label: "300 × 250" },
    banner: { width: 728, height: 90, label: "728 × 90" },
    "wide-banner": { width: 970, height: 90, label: "970 × 90" },
    skyscraper: { width: 300, height: 600, label: "300 × 600" },
  };

  const { width, height, label } = dimensions[format];

  return (
    <div
      className={`relative flex items-center justify-center bg-gradient-to-br from-muted to-muted/70 border border-border/50 rounded-lg overflow-hidden ${className}`}
      style={{
        width: format === "card" ? "100%" : `${width}px`,
        height: `${height}px`,
        minWidth: format === "card" ? "280px" : undefined,
      }}
    >
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0" style={{
          backgroundImage: `repeating-linear-gradient(45deg, currentColor 0, currentColor 1px, transparent 0, transparent 50%)`,
          backgroundSize: '10px 10px'
        }} />
      </div>

      {/* Content */}
      <div className="relative z-10 text-center space-y-2 p-4">
        <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-background/80 backdrop-blur-sm rounded-full border border-border/50">
          <svg
            className="w-4 h-4 text-muted-foreground"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
            />
          </svg>
          <span className="text-xs font-medium text-muted-foreground">Advertisement</span>
        </div>
        <p className="text-xs text-muted-foreground/60 font-mono">{label}</p>
      </div>

      {/* Close button (optional) */}
      <button
        onClick={() => setIsVisible(false)}
        className="absolute top-2 right-2 w-6 h-6 rounded-full bg-background/80 backdrop-blur-sm border border-border/50 flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-background transition-all opacity-0 hover:opacity-100 group-hover:opacity-100"
        aria-label="Close ad"
      >
        <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      {/* Google Ads Script Placeholder Comment */}
      {/*
        Replace this component with actual Google AdSense code:
        <ins className="adsbygoogle"
             style={{display:'inline-block',width:'300px',height:'250px'}}
             data-ad-client="ca-pub-XXXXXXXXXXXXXXXX"
             data-ad-slot="XXXXXXXXXX"></ins>
      */}
    </div>
  );
}
