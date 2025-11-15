"use client";

import { ExternalLink, Users, Trophy, Target, Clock } from "lucide-react";

export function KickstarterSection() {
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
            href="https://www.kickstarter.com/projects/mylittleprice/mylittleprice"
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

        {/* Kickstarter Banner Placeholder */}
        <div className="relative mb-16">
          <div className="aspect-[21/9] rounded-3xl bg-linear-to-br from-primary/20 to-secondary/20 border-2 border-primary/30 overflow-hidden relative">
            {/*
              KICKSTARTER BANNER PLACEHOLDER:
              Replace this with your actual Kickstarter campaign banner/video
              - Option 1: Campaign video embed from Kickstarter
              - Option 2: High-quality campaign banner image
              - Option 3: Animated graphics showing campaign highlights

              Recommended size: 2100x900px (21:9 aspect ratio)
              Format: MP4 for video, WebP/PNG for images
            */}
            <div className="absolute inset-0 flex items-center justify-center p-8">
              <div className="text-center space-y-4">
                <div className="text-6xl mb-4">🎬</div>
                <div className="space-y-2">
                  <p className="text-lg font-bold">Kickstarter Campaign Banner</p>
                  <p className="text-sm text-muted-foreground max-w-md">
                    [Replace with campaign video or banner image]
                  </p>
                  <p className="text-xs text-muted-foreground">
                    Suggested: Embed Kickstarter video or upload campaign visual
                  </p>
                </div>
              </div>
            </div>

            {/* Decorative corner elements */}
            <div className="absolute top-0 left-0 w-32 h-32 bg-primary/20 rounded-br-full" />
            <div className="absolute bottom-0 right-0 w-32 h-32 bg-secondary/20 rounded-tl-full" />
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
              href="https://www.kickstarter.com/projects/mylittleprice/mylittleprice"
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
