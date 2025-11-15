"use client";

import { useRouter } from "next/navigation";
import { ArrowRight, CheckCircle } from "lucide-react";

export function FinalCTA() {
  const router = useRouter();

  const benefits = [
    "No credit card required",
    "Free forever for basic features",
    "Instant product recommendations",
    "Compare prices from 100+ retailers",
  ];

  return (
    <section className="py-24 relative overflow-hidden">
      {/* Background with gradient */}
      <div className="absolute inset-0 bg-linear-to-br from-primary/10 via-background to-secondary/10" />

      {/* Animated background pattern */}
      <div className="absolute inset-0">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-pulse" />
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-secondary/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: "1s" }} />
      </div>

      <div className="container mx-auto px-4 relative">
        <div className="max-w-4xl mx-auto">
          <div className="text-center space-y-8 mb-12">
            <h2 className="text-4xl md:text-6xl font-bold leading-tight">
              Ready to Shop
              <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Smarter?
              </span>
            </h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Join our beta community and help shape the future of AI-powered shopping
            </p>

            {/* Benefits list */}
            <div className="grid md:grid-cols-2 gap-4 max-w-2xl mx-auto pt-8">
              {benefits.map((benefit, index) => (
                <div key={index} className="flex items-center gap-3 text-left">
                  <CheckCircle className="w-5 h-5 text-primary flex-shrink-0" />
                  <span className="text-sm">{benefit}</span>
                </div>
              ))}
            </div>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <button
              onClick={() => router.push("/chat")}
              className="group px-10 py-5 bg-primary text-primary-foreground rounded-full font-bold text-lg hover:opacity-90 transition-all shadow-lg hover:shadow-xl flex items-center gap-3"
            >
              Get Started Free
              <ArrowRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
            </button>
            <button
              onClick={() => router.push("/about")}
              className="px-10 py-5 bg-secondary text-secondary-foreground rounded-full font-bold text-lg hover:bg-secondary/80 transition-colors"
            >
              Learn More
            </button>
          </div>

          {/* Beta Community Badge */}
          <div className="text-center pt-12">
            <p className="text-sm text-muted-foreground">
              Join our <span className="font-semibold text-primary">Beta Testing Community</span>
            </p>
            <div className="flex justify-center gap-4 mt-4 flex-wrap">
              <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20">
                <span className="text-primary">🚀</span>
                <span className="text-sm font-medium">Beta Version</span>
              </div>
              <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20">
                <span className="text-primary">💎</span>
                <span className="text-sm font-medium">Kickstarter Live</span>
              </div>
              <div className="flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20">
                <span className="text-primary">🌍</span>
                <span className="text-sm font-medium">Global Launch</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
