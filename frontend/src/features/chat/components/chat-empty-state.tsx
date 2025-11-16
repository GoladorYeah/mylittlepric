import Image from "next/image";
import { BetaBadge } from "@/shared/components/ui/BetaBadge";
import { Sparkles, Search, TrendingDown } from "lucide-react";

export function ChatEmptyState() {
  const features = [
    {
      icon: Search,
      text: "Smart AI-powered product search"
    },
    {
      icon: TrendingDown,
      text: "Real-time price comparison"
    },
    {
      icon: Sparkles,
      text: "Personalized recommendations"
    }
  ];

  return (
    <div className="flex flex-col items-center justify-center h-full space-y-8 text-center pt-20 px-4">
      {/* Logo with Beta Badge */}
      <div className="flex items-center gap-2">
        <div className="relative inline-flex items-center">
          <Image
            src="/logo.svg"
            alt="MyLittlePrice Logo"
            width={160}
            height={60}
            className="h-16 w-auto relative z-10"
            priority
          />
        </div>
        <div className="relative z-0 opacity-70 -ml-4 mt-2">
          <BetaBadge />
        </div>
      </div>

      {/* Welcome Message */}
      <div className="space-y-3 max-w-2xl">
        <h2 className="text-3xl md:text-4xl font-bold">
          Welcome to Your AI Shopping Assistant
        </h2>
        <p className="text-lg text-muted-foreground">
          Tell me what you're looking for, and I'll help you find the perfect product at the best price
        </p>
      </div>

      {/* Features */}
      <div className="grid md:grid-cols-3 gap-4 max-w-2xl w-full">
        {features.map((feature, index) => (
          <div
            key={index}
            className="flex items-center gap-3 p-4 rounded-lg bg-secondary/30 border border-border/40"
          >
            <feature.icon className="w-5 h-5 text-primary shrink-0" />
            <span className="text-sm text-left">{feature.text}</span>
          </div>
        ))}
      </div>

      {/* Example Searches */}
      <div className="space-y-2 max-w-md">
        <p className="text-xs text-muted-foreground uppercase font-semibold">Try asking:</p>
        <div className="flex flex-wrap gap-2 justify-center">
          {["Gaming laptop under $1000", "Wireless headphones", "Running shoes"].map((example, index) => (
            <span
              key={index}
              className="px-3 py-1.5 rounded-full bg-primary/10 text-primary text-xs font-medium border border-primary/20"
            >
              {example}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}
