"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Search } from "lucide-react";
import { Logo } from "@/shared/components/ui";

export function HeroSection() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");
  const [isAnimating, setIsAnimating] = useState(false);

  const handleSearch = () => {
    if (searchQuery.trim()) {
      setIsAnimating(true);
      setTimeout(() => {
        router.push(`/chat?q=${encodeURIComponent(searchQuery)}`);
      }, 500);
    } else {
      router.push("/chat");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-background border-b border-border">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Logo width={84.24} height={32} />
        </div>

        <div className="flex-1 max-w-2xl mx-8">
          <div
            className={`relative transition-all duration-500 ${
              isAnimating ? "scale-110 opacity-0" : "scale-100 opacity-100"
            }`}
          >
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={handleKeyDown}
              onClick={handleSearch}
              placeholder="Search for any product..."
              className="w-full px-6 py-3 pl-12 rounded-full bg-secondary border border-border focus:border-primary focus:outline-none transition-colors"
            />
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
          </div>
        </div>

        <div className="flex items-center gap-3">
          <button
            onClick={() => router.push("/login")}
            className="px-6 py-2 rounded-full font-medium text-sm hover:bg-secondary transition-colors"
          >
            Sign In
          </button>
        </div>
      </div>
    </header>
  );
}
