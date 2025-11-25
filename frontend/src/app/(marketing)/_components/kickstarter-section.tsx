"use client";

import { ExternalLink, Users, Trophy, Target, Clock, TrendingUp, Heart, Rocket } from "lucide-react";
import { useState, useEffect } from "react";

export function KickstarterSection() {
  const [animationPhase, setAnimationPhase] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setAnimationPhase((prev) => (prev + 1) % 3);
    }, 4000);
    return () => clearInterval(timer);
  }, []);
  const rewards = [
    {
      tier: "$10",
      name: "Launch Supporter",
      benefits: ["Name on Launch Wall"],
      badge: "🎯",
    },
    {
      tier: "$25",
      name: "Founding Member",
      benefits: ["Founding Member Badge", "Private Discord Access"],
      badge: "🏅",
      popular: true,
    },
    {
      tier: "$49",
      name: "Early Access",
      benefits: ["12-month Pro Subscription", "Private Discord", "Beta Voting"],
      badge: "⭐",
      popular: true,
    },
    {
      tier: "$99",
      name: "Pro 24",
      benefits: ["24-month Pro Subscription", "Private Discord", "2× Beta Voting"],
      badge: "💎",
    },
  ];

  const features = [
    {
      icon: Users,
      title: "Join the Community",
      description: "Be part of the founding members shaping the future of AI shopping",
    },
    {
      icon: Trophy,
      title: "Exclusive Rewards",
      description: "Get Pro subscriptions, founding badges, and special perks",
    },
    {
      icon: Target,
      title: "Support Fair Pricing",
      description: "Help us build transparent shopping where AI works for users",
    },
  ];

  return (
    <section className="py-24 relative overflow-hidden bg-linear-to-br from-primary/5 via-background to-secondary/5">
      {/* Background decoration */}
      <div className="absolute inset-0">
        <div className="absolute top-20 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl" />
        <div className="absolute bottom-20 right-1/4 w-96 h-96 bg-secondary/10 rounded-full blur-3xl" />
      </div>

      <div className="container mx-auto px-4 relative">
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16 space-y-6">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20 text-sm">
            <Clock className="w-4 h-4 text-primary" />
            <span className="text-primary font-medium">Campaign Live Now</span>
          </div>

          <h2 className="text-4xl md:text-6xl font-bold">
            Support Our
            <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              Kickstarter Campaign
            </span>
          </h2>

          <p className="text-xl text-muted-foreground">
            We're scaling globally — connecting new markets, adding Pro features, and building our community of shoppers who believe in fair, transparent pricing.
          </p>

          <a
            href="https://www.kickstarter.com/projects/mylittleprice/mylittleprice-ai-find-ideal-solution-at-the-right-price"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-3 px-10 py-5 bg-primary text-primary-foreground rounded-full font-bold text-lg hover:opacity-90 transition-all shadow-lg hover:shadow-xl"
          >
            <span className="text-2xl">🚀</span>
            Back Us on Kickstarter
            <ExternalLink className="w-6 h-6" />
          </a>
        </div>

        {/* Why Support Section */}
        <div className="grid md:grid-cols-3 gap-6 mb-16">
          {features.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <div
                key={index}
                className="p-6 rounded-2xl bg-background/50 backdrop-blur-sm border border-border hover:border-primary/50 transition-all"
              >
                <div className="w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center mb-4">
                  <Icon className="w-6 h-6 text-primary" />
                </div>
                <h3 className="text-lg font-bold mb-2">{feature.title}</h3>
                <p className="text-sm text-muted-foreground">{feature.description}</p>
              </div>
            );
          })}
        </div>

        {/* Kickstarter Campaign Animation Banner */}
        <div className="relative mb-16">
          <div className="aspect-video md:aspect-[21/9] rounded-2xl md:rounded-3xl bg-linear-to-br from-primary/10 via-background to-secondary/10 border-2 border-primary/30 overflow-hidden relative shadow-2xl">
            {/* Animation Phase 0: Community Growth */}
            {animationPhase === 0 && (
              <div className="absolute inset-0 p-6 md:p-12 flex flex-col md:flex-row items-center justify-center md:justify-between animate-fade-in-up">
                <div className="flex-1 space-y-3 md:space-y-4 text-center md:text-left">
                  <div className="inline-flex items-center gap-2 px-3 py-1.5 md:px-4 md:py-2 rounded-full bg-primary/20 border border-primary/30">
                    <Users className="w-4 h-4 md:w-5 md:h-5 text-primary" />
                    <span className="text-xs md:text-sm font-bold text-primary">Join Our Community</span>
                  </div>
                  <h3 className="text-2xl md:text-3xl lg:text-5xl font-bold">
                    <span className="bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                      1,000+ Backers
                    </span>
                  </h3>
                  <p className="text-sm md:text-lg text-muted-foreground max-w-md mx-auto md:mx-0">
                    Be part of the founding members shaping the future of AI-powered shopping
                  </p>
                </div>
                <div className="hidden md:flex flex-1 items-center justify-center relative">
                  <div className="relative">
                    <Users className="w-24 h-24 lg:w-32 lg:h-32 text-primary/30 animate-pulse-scale" />
                    {/* Animated user circles */}
                    {[...Array(8)].map((_, i) => (
                      <div
                        key={i}
                        className="absolute top-1/2 left-1/2 w-6 h-6 lg:w-8 lg:h-8 -ml-3 lg:-ml-4 -mt-3 lg:-mt-4 rounded-full bg-primary/40 border-2 border-background animate-float-particle"
                        style={{
                          animationDelay: `${i * 0.3}s`,
                          transform: `rotate(${i * 45}deg) translateY(-50px) translateY(-10px)`
                        }}
                      >
                        <div className="w-full h-full rounded-full bg-primary/60 flex items-center justify-center text-xs">
                          👤
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {/* Animation Phase 1: Campaign Progress */}
            {animationPhase === 1 && (
              <div className="absolute inset-0 p-6 md:p-8 lg:p-12 flex flex-col justify-center animate-fade-in-up">
                <div className="space-y-4 md:space-y-6 max-w-3xl mx-auto w-full">
                  <div className="text-center space-y-2">
                    <div className="inline-flex items-center gap-2 px-3 py-1.5 md:px-4 md:py-2 rounded-full bg-green-500/20 border border-green-500/30 mb-2 md:mb-4">
                      <TrendingUp className="w-4 h-4 md:w-5 md:h-5 text-green-500" />
                      <span className="text-xs md:text-sm font-bold text-green-500">Campaign Progress</span>
                    </div>
                    <h3 className="text-2xl md:text-3xl lg:text-5xl font-bold">
                      <span className="bg-linear-to-r from-green-500 to-primary bg-clip-text text-transparent">
                        $50,000
                      </span>
                      <span className="text-lg md:text-2xl text-muted-foreground ml-1 md:ml-2">/ $100,000</span>
                    </h3>
                    <p className="text-sm md:text-lg text-muted-foreground">
                      50% funded in just 2 weeks!
                    </p>
                  </div>

                  <div className="relative h-4 md:h-6 bg-background/50 rounded-full overflow-hidden border border-border">
                    <div className="absolute inset-0 bg-linear-to-r from-primary/20 to-green-500/20" />
                    <div className="h-full bg-linear-to-r from-primary to-green-500 rounded-full animate-progress" style={{ width: '50%' }}>
                      <div className="absolute inset-0 bg-linear-to-r from-transparent via-white/20 to-transparent animate-shimmer" />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-2 md:gap-4 mt-4 md:mt-8">
                    <div className="text-center p-2 md:p-4 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border">
                      <div className="text-lg md:text-2xl font-bold text-primary mb-0.5 md:mb-1">850+</div>
                      <div className="text-[10px] md:text-xs text-muted-foreground">Backers</div>
                    </div>
                    <div className="text-center p-2 md:p-4 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border">
                      <div className="text-lg md:text-2xl font-bold text-primary mb-0.5 md:mb-1">21</div>
                      <div className="text-[10px] md:text-xs text-muted-foreground">Days Left</div>
                    </div>
                    <div className="text-center p-2 md:p-4 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border">
                      <div className="text-lg md:text-2xl font-bold text-green-500 mb-0.5 md:mb-1">$59</div>
                      <div className="text-[10px] md:text-xs text-muted-foreground">Avg Pledge</div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Animation Phase 2: Rewards & Benefits */}
            {animationPhase === 2 && (
              <div className="absolute inset-0 p-6 md:p-8 lg:p-12 flex flex-col md:flex-row items-center justify-center md:justify-between animate-fade-in-up">
                <div className="flex-1 space-y-4 md:space-y-6 text-center md:text-left">
                  <div className="inline-flex items-center gap-2 px-3 py-1.5 md:px-4 md:py-2 rounded-full bg-yellow-500/20 border border-yellow-500/30">
                    <Trophy className="w-4 h-4 md:w-5 md:h-5 text-yellow-500" />
                    <span className="text-xs md:text-sm font-bold text-yellow-500">Exclusive Rewards</span>
                  </div>
                  <h3 className="text-2xl md:text-3xl lg:text-5xl font-bold">
                    <span className="bg-linear-to-r from-yellow-500 to-primary bg-clip-text text-transparent">
                      Early Access
                    </span>
                  </h3>
                  <div className="space-y-2 md:space-y-3 max-w-md mx-auto md:mx-0">
                    <div className="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border animate-bounce-in">
                      <div className="text-xl md:text-2xl">🏅</div>
                      <div className="text-left">
                        <div className="font-semibold text-xs md:text-sm">Founding Member Badge</div>
                        <div className="text-[10px] md:text-xs text-muted-foreground">Lifetime recognition</div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border animate-bounce-in" style={{ animationDelay: '0.1s' }}>
                      <div className="text-xl md:text-2xl">⭐</div>
                      <div className="text-left">
                        <div className="font-semibold text-xs md:text-sm">12-Month Pro Access</div>
                        <div className="text-[10px] md:text-xs text-muted-foreground">$120 value</div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2 md:gap-3 p-2 md:p-3 rounded-lg md:rounded-xl bg-background/50 backdrop-blur-sm border border-border animate-bounce-in" style={{ animationDelay: '0.2s' }}>
                      <div className="text-xl md:text-2xl">💬</div>
                      <div className="text-left">
                        <div className="font-semibold text-xs md:text-sm">Private Discord Community</div>
                        <div className="text-[10px] md:text-xs text-muted-foreground">Direct dev access</div>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="hidden md:flex flex-1 items-center justify-center">
                  <div className="relative">
                    <Rocket className="w-32 h-32 lg:w-40 lg:h-40 text-primary/30 animate-float" />
                    <div className="absolute inset-0 bg-primary/10 rounded-full blur-3xl animate-pulse" />
                  </div>
                </div>
              </div>
            )}

            {/* Decorative elements */}
            <div className="absolute top-0 left-0 w-16 h-16 md:w-32 md:h-32 bg-primary/10 rounded-br-full opacity-50" />
            <div className="absolute bottom-0 right-0 w-16 h-16 md:w-32 md:h-32 bg-secondary/10 rounded-tl-full opacity-50" />

            {/* Phase indicators */}
            <div className="absolute bottom-2 md:bottom-4 left-1/2 -translate-x-1/2 flex gap-1.5 md:gap-2 z-10">
              {[0, 1, 2].map((i) => (
                <div
                  key={i}
                  className={`h-1 md:h-1.5 rounded-full transition-all ${
                    animationPhase === i ? 'bg-primary w-6 md:w-8' : 'bg-primary/30 w-1 md:w-1.5'
                  }`}
                />
              ))}
            </div>
          </div>
        </div>

        {/* Reward Tiers */}
        <div className="max-w-5xl mx-auto">
          <h3 className="text-3xl font-bold text-center mb-8">
            Backer Rewards
          </h3>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {rewards.map((reward, index) => (
              <div
                key={index}
                className={`relative p-6 rounded-2xl border-2 transition-all ${
                  reward.popular
                    ? "border-primary bg-primary/5 shadow-lg shadow-primary/20"
                    : "border-border bg-background hover:border-primary/50"
                }`}
              >
                {reward.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-primary text-primary-foreground text-xs font-bold">
                    POPULAR
                  </div>
                )}

                <div className="text-center space-y-4">
                  <div className="text-4xl">{reward.badge}</div>
                  <div>
                    <div className="text-3xl font-bold text-primary mb-1">
                      {reward.tier}
                    </div>
                    <div className="text-sm font-semibold">{reward.name}</div>
                  </div>

                  <ul className="space-y-2 text-sm text-muted-foreground">
                    {reward.benefits.map((benefit, i) => (
                      <li key={i} className="flex items-start gap-2">
                        <span className="text-primary mt-0.5">✓</span>
                        <span>{benefit}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            ))}
          </div>

          <div className="text-center mt-8">
            <a
              href="https://www.kickstarter.com/projects/mylittleprice/mylittleprice-ai-find-ideal-solution-at-the-right-price"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 text-primary hover:underline font-medium"
            >
              View All Reward Tiers & Partner Options
              <ExternalLink className="w-4 h-4" />
            </a>
          </div>
        </div>

        {/* Campaign Goals */}
        <div className="mt-16 p-8 rounded-3xl bg-background/50 backdrop-blur-sm border border-border max-w-4xl mx-auto">
          <h3 className="text-2xl font-bold mb-6 text-center">How We'll Use Your Support</h3>
          <div className="space-y-4">
            <div className="flex items-start gap-4">
              <div className="w-12 h-2 bg-primary/40 rounded-full mt-2 flex-shrink-0">
                <div className="w-[40%] h-full bg-primary rounded-full" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between mb-1">
                  <span className="font-semibold">AI Optimization & Infrastructure</span>
                  <span className="text-sm text-muted-foreground">40%</span>
                </div>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-12 h-2 bg-primary/40 rounded-full mt-2 flex-shrink-0">
                <div className="w-[25%] h-full bg-primary rounded-full" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between mb-1">
                  <span className="font-semibold">Market Expansion & API Integrations</span>
                  <span className="text-sm text-muted-foreground">25%</span>
                </div>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-12 h-2 bg-primary/40 rounded-full mt-2 flex-shrink-0">
                <div className="w-[20%] h-full bg-primary rounded-full" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between mb-1">
                  <span className="font-semibold">Marketing & Partnerships</span>
                  <span className="text-sm text-muted-foreground">20%</span>
                </div>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-12 h-2 bg-primary/40 rounded-full mt-2 flex-shrink-0">
                <div className="w-[15%] h-full bg-primary rounded-full" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between mb-1">
                  <span className="font-semibold">UX/UI, Legal & Platform Fees</span>
                  <span className="text-sm text-muted-foreground">15%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
