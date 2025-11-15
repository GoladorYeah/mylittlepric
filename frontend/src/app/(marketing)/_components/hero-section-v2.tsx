"use client";

import { useRouter } from "next/navigation";
import { ArrowRight, Sparkles, ExternalLink } from "lucide-react";
import { BetaBadgeLarge } from "@/shared/components/ui/BetaBadge";

export function HeroSection() {
  const router = useRouter();

  return (
    <section className="relative min-h-[85vh] flex items-center overflow-hidden py-12">
      {/* Background gradient */}
      <div className="absolute inset-0 bg-linear-to-br from-primary/10 via-background to-secondary/20" />

      {/* Animated grid pattern */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#8882_1px,transparent_1px),linear-gradient(to_bottom,#8882_1px,transparent_1px)] bg-size-[64px_64px] opacity-20" />

      <div className="container mx-auto px-4 relative z-10">
        <div className="grid lg:grid-cols-2 gap-8 lg:gap-12 items-center">
          {/* Left side - Content */}
          <div className="space-y-6">
            {/* Beta Badge Highlight */}
            <div className="flex items-center gap-3 flex-wrap">
              <BetaBadgeLarge />
              <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-primary/10 border border-primary/20 text-xs">
                <Sparkles className="w-3.5 h-3.5 text-primary" />
                <span className="text-primary font-medium">AI-Powered Assistant</span>
              </div>
            </div>

            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold leading-tight">
              Find Your Perfect
              <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Product in Seconds
              </span>
            </h1>

            <p className="text-lg md:text-xl text-muted-foreground max-w-xl">
              Join our beta testing community! Help us build the future of AI-powered shopping. Never overpay again — get instant product recommendations with real-time price comparison.
            </p>

            {/* Compact Kickstarter Banner */}
            <div className="p-3 rounded-xl bg-linear-to-r from-primary/5 to-secondary/5 border border-primary/20 max-w-xl">
              <div className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-2 flex-1 min-w-0">
                  <span className="text-xl shrink-0">🚀</span>
                  <div className="min-w-0">
                    <p className="text-xs font-semibold truncate">Live on Kickstarter!</p>
                    <p className="text-xs text-muted-foreground truncate">Get exclusive early access rewards</p>
                  </div>
                </div>
                <a
                  href="https://www.kickstarter.com/projects/mylittleprice/mylittleprice"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1.5 px-4 py-1.5 bg-primary text-primary-foreground rounded-full font-medium hover:opacity-90 transition-all text-xs whitespace-nowrap shrink-0"
                >
                  View
                  <ExternalLink className="w-3 h-3" />
                </a>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-3">
              <button
                onClick={() => router.push("/chat")}
                className="group px-8 py-3.5 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-all flex items-center justify-center gap-2 shadow-lg"
              >
                Try Beta Version Free
                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </button>
              <button
                onClick={() => {
                  document.getElementById('demo')?.scrollIntoView({ behavior: 'smooth' });
                }}
                className="px-8 py-3.5 bg-secondary text-secondary-foreground rounded-full font-semibold hover:bg-secondary/80 transition-colors"
              >
                Watch Demo
              </button>
            </div>

            {/* Beta indicators */}
            <div className="flex items-center gap-4 flex-wrap text-sm">
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-primary rounded-full animate-pulse" />
                <span className="text-muted-foreground">Beta Access Available</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-primary rounded-full animate-pulse" />
                <span className="text-muted-foreground">Free Forever (Basic)</span>
              </div>
            </div>
          </div>

          {/* Right side - Visual placeholder */}
          <div className="relative lg:order-last">
            {/*
              VIDEO/IMAGE PLACEHOLDER:
              Replace this with an animated video showing:
              - User typing a query
              - AI analyzing the request
              - Product results appearing
              - User selecting a product

              Recommended: Short looping video (5-10 seconds), MP4 format, optimized for web
            */}
            <div className="relative aspect-square rounded-3xl bg-linear-to-br from-primary/20 to-secondary/20 border border-border overflow-hidden">
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="text-center space-y-4 p-8">
                  <div className="w-32 h-32 mx-auto rounded-2xl bg-linear-to-br from-primary/30 to-secondary/30 flex items-center justify-center">
                    <Sparkles className="w-16 h-16 text-primary" />
                  </div>
                  <p className="text-sm text-muted-foreground">
                    [VIDEO: AI Shopping Demo]
                    <br />
                    <span className="text-xs">
                      Looping animation of product search
                    </span>
                  </p>
                </div>
              </div>

              {/* Floating elements animation */}
              <div className="absolute top-6 right-6 p-3 rounded-xl bg-background/80 backdrop-blur-sm border border-border shadow-lg animate-float">
                <div className="text-xs font-medium">💰 Best Price</div>
                <div className="text-sm font-bold text-primary">$299.99</div>
              </div>

              <div className="absolute bottom-6 left-6 p-3 rounded-xl bg-background/80 backdrop-blur-sm border border-border shadow-lg animate-float" style={{ animationDelay: "1s" }}>
                <div className="text-xs font-medium">⚡ AI Match</div>
                <div className="text-sm font-bold text-primary">98%</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
