"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Search, Zap, Shield, TrendingUp } from "lucide-react";
import { Logo } from "@/components/Logo";

export default function HomePage() {
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
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
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
        </div>
      </header>

      <main className="pt-32 pb-20">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center space-y-8">
            <div className="space-y-4">
              <h1 className="text-5xl md:text-6xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Find the Perfect Product
              </h1>
              <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
                AI-powered shopping assistant that helps you discover the best
                products at the best prices
              </p>
            </div>

            <div className="flex gap-4 justify-center pt-8">
              <button
                onClick={() => router.push("/chat")}
                className="px-8 py-4 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-opacity"
              >
                Start Shopping
              </button>
              <button
                onClick={() => router.push("/chat")}
                className="px-8 py-4 bg-secondary text-secondary-foreground rounded-full font-semibold hover:bg-secondary/80 transition-colors"
              >
                Learn More
              </button>
            </div>

            <div className="grid md:grid-cols-3 gap-8 pt-20">
              <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
                  <Zap className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-semibold mb-2">Lightning Fast</h3>
                <p className="text-muted-foreground text-sm">
                  Get instant product recommendations powered by advanced AI
                </p>
              </div>

              <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
                  <Shield className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-semibold mb-2">Smart & Secure</h3>
                <p className="text-muted-foreground text-sm">
                  Verified products from trusted merchants across Europe
                </p>
              </div>

              <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
                  <TrendingUp className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-semibold mb-2">Best Prices</h3>
                <p className="text-muted-foreground text-sm">
                  Compare prices and find the best deals automatically
                </p>
              </div>
            </div>

            <div className="pt-20 space-y-4">
              <h2 className="text-3xl font-bold">How It Works</h2>
              <div className="grid md:grid-cols-3 gap-8 pt-8 text-left">
                <div className="space-y-2">
                  <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                    1
                  </div>
                  <h4 className="font-semibold">Tell us what you need</h4>
                  <p className="text-sm text-muted-foreground">
                    Describe the product you're looking for in natural language
                  </p>
                </div>

                <div className="space-y-2">
                  <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                    2
                  </div>
                  <h4 className="font-semibold">AI finds the best matches</h4>
                  <p className="text-sm text-muted-foreground">
                    Our AI searches thousands of products to find perfect
                    matches
                  </p>
                </div>

                <div className="space-y-2">
                  <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
                    3
                  </div>
                  <h4 className="font-semibold">Compare and buy</h4>
                  <p className="text-sm text-muted-foreground">
                    Review options, compare prices, and make your choice
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      <footer className="border-t border-border py-8">
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