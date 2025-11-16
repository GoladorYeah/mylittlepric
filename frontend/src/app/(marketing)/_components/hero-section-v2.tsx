"use client";

import { useRouter } from "next/navigation";
import { ArrowRight, Sparkles, ExternalLink, Search, TrendingDown, Zap } from "lucide-react";
import { useState, useEffect } from "react";

export function HeroSection() {
  const router = useRouter();
  const [phase, setPhase] = useState(0); // 0: input, 1: processing, 2: searching, 3: results

  useEffect(() => {
    const timer = setInterval(() => {
      setPhase((prev) => (prev + 1) % 4);
    }, 3500); // Change phase every 3.5 seconds
    return () => clearInterval(timer);
  }, []);

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
            {/* AI-Powered Badge */}
            <div className="flex items-center gap-3 flex-wrap">
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

          {/* Right side - Interactive Animation */}
          <div className="relative lg:order-last">
            <div className="relative aspect-square rounded-3xl bg-linear-to-br from-primary/10 via-background to-secondary/10 border border-border overflow-hidden">
              {/* Phase 0: User Input */}
              {phase === 0 && (
                <div className="absolute inset-0 p-8 flex flex-col justify-center items-center">
                  <div className="w-full max-w-md space-y-6 animate-fade-in-up">
                    <div className="text-center space-y-2">
                      <Search className="w-12 h-12 text-primary mx-auto mb-4" />
                      <h3 className="text-lg font-semibold">Step 1: Tell us what you need</h3>
                      <p className="text-sm text-muted-foreground">Natural language search</p>
                    </div>
                    <div className="relative p-4 rounded-2xl bg-background/90 backdrop-blur-sm border border-primary/50 shadow-lg">
                      <div className="flex items-center gap-3">
                        <Sparkles className="w-5 h-5 text-primary animate-pulse-scale" />
                        <div className="flex-1">
                          <p className="text-sm font-medium animate-typing">
                            Gaming laptop under $1000
                          </p>
                        </div>
                      </div>
                      <div className="mt-3 flex items-center gap-2 text-xs text-muted-foreground">
                        <div className="w-2 h-2 bg-primary rounded-full animate-pulse" />
                        <span>Typing...</span>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Phase 1: AI Processing */}
              {phase === 1 && (
                <div className="absolute inset-0 p-8 flex flex-col justify-center items-center">
                  <div className="w-full max-w-md space-y-8 animate-fade-in-up">
                    <div className="text-center space-y-2">
                      <div className="relative w-24 h-24 mx-auto mb-4">
                        <Sparkles className="w-24 h-24 text-primary animate-pulse-scale" />
                        <div className="absolute inset-0 bg-primary/20 rounded-full blur-2xl animate-pulse" />
                      </div>
                      <h3 className="text-lg font-semibold">Step 2: AI Analyzes Request</h3>
                      <p className="text-sm text-muted-foreground">Understanding your needs</p>
                    </div>
                    <div className="space-y-4">
                      <div className="flex items-center gap-3">
                        <Zap className="w-5 h-5 text-yellow-500 animate-pulse" />
                        <div className="flex-1">
                          <div className="h-2 bg-primary/30 rounded-full overflow-hidden">
                            <div className="h-full bg-primary rounded-full animate-progress" />
                          </div>
                        </div>
                        <span className="text-xs text-muted-foreground">Analyzing...</span>
                      </div>
                      <div className="flex items-center gap-3">
                        <Search className="w-5 h-5 text-blue-500 animate-pulse" style={{ animationDelay: '0.2s' }} />
                        <div className="flex-1">
                          <div className="h-2 bg-primary/30 rounded-full overflow-hidden">
                            <div className="h-full bg-primary rounded-full animate-progress" style={{ animationDelay: '0.3s' }} />
                          </div>
                        </div>
                        <span className="text-xs text-muted-foreground">Processing...</span>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Phase 2: Searching */}
              {phase === 2 && (
                <div className="absolute inset-0 p-8 flex flex-col">
                  <div className="text-center mb-6 animate-fade-in-up">
                    <Search className="w-12 h-12 text-primary mx-auto mb-3 animate-pulse-scale" />
                    <h3 className="text-lg font-semibold">Step 3: Searching Products</h3>
                    <p className="text-sm text-muted-foreground">Scanning 1000+ products</p>
                  </div>
                  <div className="flex-1 grid grid-cols-2 gap-4">
                    {[1, 2, 3, 4].map((item, index) => (
                      <div
                        key={item}
                        className="relative bg-background/80 backdrop-blur-sm rounded-xl p-4 border border-border shadow-lg animate-fade-in-up"
                        style={{ animationDelay: `${index * 0.15}s` }}
                      >
                        <div className="aspect-square bg-linear-to-br from-primary/20 to-secondary/20 rounded-lg mb-3 relative overflow-hidden">
                          <div className="absolute inset-0 flex items-center justify-center">
                            <div className="text-3xl">🎮</div>
                          </div>
                          <div className="absolute inset-0 bg-linear-to-b from-primary/30 to-transparent animate-scan-vertical" style={{ animationDelay: `${index * 0.2}s` }} />
                        </div>
                        <div className="space-y-2">
                          <div className="flex items-center gap-2">
                            <Search className="w-3 h-3 text-primary animate-pulse" style={{ animationDelay: `${index * 0.1}s` }} />
                            <div className="h-2 bg-primary/30 rounded-full flex-1 overflow-hidden">
                              <div className="h-full bg-primary rounded-full animate-progress" style={{ animationDelay: `${index * 0.2}s` }} />
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Phase 3: Results & Comparison */}
              {phase === 3 && (
                <div className="absolute inset-0 p-8 flex flex-col">
                  <div className="text-center mb-6 animate-fade-in-up">
                    <TrendingDown className="w-12 h-12 text-green-500 mx-auto mb-3" />
                    <h3 className="text-lg font-semibold">Step 4: Best Matches Found!</h3>
                    <p className="text-sm text-muted-foreground">Comparing prices for you</p>
                  </div>
                  <div className="flex-1 grid grid-cols-2 gap-4">
                    {[1, 2, 3, 4].map((item, index) => (
                      <div
                        key={item}
                        className="relative bg-background/90 backdrop-blur-sm rounded-xl p-4 border border-border shadow-lg animate-fade-in-up hover:border-primary/50 transition-all cursor-pointer hover:scale-105"
                        style={{ animationDelay: `${index * 0.15}s` }}
                      >
                        <div className="aspect-square bg-linear-to-br from-primary/20 to-secondary/20 rounded-lg mb-3 relative overflow-hidden">
                          <div className="absolute inset-0 flex items-center justify-center">
                            <div className="text-4xl">🎮</div>
                          </div>
                        </div>
                        <div className="space-y-2">
                          <div className="h-2 bg-foreground/10 rounded-full" />
                          <div className="h-2 bg-foreground/10 rounded-full w-2/3" />
                          <div className="flex items-center justify-between mt-3">
                            <div className="flex items-center gap-1">
                              <TrendingDown className="w-3 h-3 text-green-500" />
                              <span className="text-xs text-green-500 font-semibold">Best</span>
                            </div>
                            <div className="text-sm font-bold text-primary animate-bounce-in" style={{ animationDelay: `${index * 0.2 + 0.3}s` }}>
                              ${749 + index * 50}
                            </div>
                          </div>
                        </div>
                        {index === 0 && (
                          <div className="absolute top-2 right-2 bg-primary text-primary-foreground text-[10px] font-bold px-2 py-1 rounded-full animate-bounce-in">
                            Best Match
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                  <div className="absolute top-6 right-6 p-3 rounded-xl bg-green-500/90 text-white backdrop-blur-sm shadow-lg animate-bounce-in" style={{ animationDelay: '0.6s' }}>
                    <div className="text-xs font-medium">💰 You Save</div>
                    <div className="text-sm font-bold">$250+</div>
                  </div>
                </div>
              )}

              {/* Phase indicator dots */}
              <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2">
                {[0, 1, 2, 3].map((i) => (
                  <div
                    key={i}
                    className={`w-2 h-2 rounded-full transition-all ${
                      phase === i ? 'bg-primary w-6' : 'bg-primary/30'
                    }`}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
