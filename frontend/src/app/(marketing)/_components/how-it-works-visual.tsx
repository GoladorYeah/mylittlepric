"use client";

import { ArrowRight, MessageCircle, Sparkles, ShoppingCart, Search, Zap, Send, TrendingDown, ExternalLink } from "lucide-react";

// Step 1: Chat Interface Animation
function ChatInterfaceAnimation() {
  const messages = [
    { type: 'user', text: 'Looking for gaming laptop under $1000', delay: 0 },
    { type: 'ai', text: 'I\'ll help you find the perfect gaming laptop!', delay: 0.8 },
    { type: 'ai', text: 'Searching for options...', delay: 1.6 },
  ];

  return (
    <div className="relative w-full h-full bg-linear-to-br from-blue-500/10 to-cyan-500/10 rounded-3xl p-6 overflow-hidden flex flex-col">
      {/* Chat header */}
      <div className="flex items-center gap-3 pb-4 border-b border-border/30">
        <div className="w-10 h-10 rounded-full bg-primary/20 flex items-center justify-center">
          <Sparkles className="w-5 h-5 text-primary" />
        </div>
        <div>
          <h4 className="text-sm font-semibold">AI Shopping Assistant</h4>
          <p className="text-xs text-muted-foreground">Online</p>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 py-4 space-y-3 overflow-hidden">
        {messages.map((msg, index) => (
          <div
            key={index}
            className={`flex ${msg.type === 'user' ? 'justify-end' : 'justify-start'} animate-fade-in-up`}
            style={{ animationDelay: `${msg.delay}s` }}
          >
            <div className={`max-w-[80%] px-4 py-2.5 rounded-2xl ${
              msg.type === 'user'
                ? 'bg-primary text-primary-foreground rounded-br-sm'
                : 'bg-secondary text-secondary-foreground rounded-bl-sm'
            }`}>
              <p className="text-xs leading-relaxed">{msg.text}</p>
            </div>
          </div>
        ))}

        {/* Typing indicator */}
        <div className="flex justify-start animate-fade-in-up" style={{ animationDelay: '2.4s' }}>
          <div className="bg-secondary px-4 py-3 rounded-2xl rounded-bl-sm">
            <div className="flex gap-1">
              <div className="w-2 h-2 bg-primary/60 rounded-full animate-bounce" style={{ animationDelay: '0s' }} />
              <div className="w-2 h-2 bg-primary/60 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
              <div className="w-2 h-2 bg-primary/60 rounded-full animate-bounce" style={{ animationDelay: '0.4s' }} />
            </div>
          </div>
        </div>
      </div>

      {/* Input field */}
      <div className="pt-4 border-t border-border/30">
        <div className="flex items-center gap-2 px-4 py-2.5 bg-background/80 backdrop-blur-sm rounded-full border border-border">
          <MessageCircle className="w-4 h-4 text-muted-foreground" />
          <div className="flex-1 text-xs text-muted-foreground">Type your message...</div>
          <Send className="w-4 h-4 text-primary" />
        </div>
      </div>

      {/* Floating chat bubbles decoration */}
      <div className="absolute -top-4 -right-4 w-20 h-20 bg-primary/10 rounded-full blur-2xl animate-pulse" />
      <div className="absolute -bottom-4 -left-4 w-20 h-20 bg-cyan-500/10 rounded-full blur-2xl animate-pulse" style={{ animationDelay: '1s' }} />
    </div>
  );
}

// AI Analysis Animation Component for Step 2
function AIAnalysisAnimation() {
  return (
    <div className="relative w-full h-full bg-linear-to-br from-purple-500/10 to-pink-500/10 rounded-3xl p-8 overflow-hidden">
      {/* Animated scanning lines */}
      <div className="absolute inset-0 opacity-30">
        <div className="absolute w-full h-0.5 bg-linear-to-r from-transparent via-primary to-transparent animate-scan" />
        <div className="absolute w-full h-0.5 bg-linear-to-r from-transparent via-primary to-transparent animate-scan-delayed" />
      </div>

      {/* Product cards being analyzed */}
      <div className="relative grid grid-cols-2 gap-4 h-full">
        {[1, 2, 3, 4].map((item, index) => (
          <div
            key={item}
            className="relative bg-background/80 backdrop-blur-sm rounded-xl p-4 border border-border shadow-lg animate-fade-in-up"
            style={{ animationDelay: `${index * 0.2}s` }}
          >
            {/* Product image placeholder */}
            <div className="aspect-square bg-linear-to-br from-primary/20 to-secondary/20 rounded-lg mb-3 relative overflow-hidden">
              <div className="absolute inset-0 flex items-center justify-center">
                <ShoppingCart className="w-8 h-8 text-primary/50" />
              </div>
              {/* Scanning effect */}
              <div className="absolute inset-0 bg-linear-to-b from-primary/30 to-transparent animate-scan-vertical"
                   style={{ animationDelay: `${index * 0.3}s` }} />
            </div>

            {/* Analysis indicators */}
            <div className="space-y-2">
              <div className="flex items-center gap-2 text-xs">
                <Search className="w-3 h-3 text-primary animate-pulse" style={{ animationDelay: `${index * 0.15}s` }} />
                <div className="h-2 bg-primary/30 rounded-full flex-1 overflow-hidden">
                  <div className="h-full bg-primary rounded-full animate-progress" style={{ animationDelay: `${index * 0.2}s` }} />
                </div>
              </div>
              <div className="flex items-center gap-2 text-xs">
                <Zap className="w-3 h-3 text-yellow-500 animate-pulse" style={{ animationDelay: `${index * 0.2}s` }} />
                <span className="text-muted-foreground text-[10px]">Analyzing...</span>
              </div>
            </div>

            {/* Match percentage badge */}
            <div className="absolute top-2 right-2 bg-primary/90 text-primary-foreground text-xs font-bold px-2 py-1 rounded-full animate-bounce-in"
                 style={{ animationDelay: `${index * 0.4 + 0.8}s` }}>
              {95 - index * 5}%
            </div>
          </div>
        ))}
      </div>

      {/* Central AI brain indicator */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none">
        <div className="relative">
          <Sparkles className="w-12 h-12 text-primary animate-pulse-scale" />
          <div className="absolute inset-0 bg-primary/20 rounded-full blur-xl animate-pulse" />
        </div>
      </div>

      {/* Floating particles */}
      {[...Array(6)].map((_, i) => (
        <div
          key={i}
          className="absolute w-1 h-1 bg-primary/60 rounded-full animate-float-particle"
          style={{
            left: `${20 + i * 15}%`,
            top: `${30 + (i % 3) * 20}%`,
            animationDelay: `${i * 0.3}s`,
            animationDuration: `${3 + i * 0.5}s`
          }}
        />
      ))}
    </div>
  );
}

// Step 3: Product Comparison Animation
function ProductComparisonAnimation() {
  const products = [
    { name: 'ASUS ROG Strix', price: '$899', rating: 4.8, badge: 'Best Value', delay: 0 },
    { name: 'MSI Katana 15', price: '$949', rating: 4.6, badge: 'Popular', delay: 0.3 },
    { name: 'Lenovo Legion 5', price: '$979', rating: 4.7, badge: 'Premium', delay: 0.6 },
  ];

  return (
    <div className="relative w-full h-full bg-linear-to-br from-green-500/10 to-emerald-500/10 rounded-3xl p-6 overflow-hidden">
      {/* Header */}
      <div className="mb-4 animate-fade-in-up">
        <h4 className="text-sm font-semibold mb-1">Top 3 Matches</h4>
        <p className="text-xs text-muted-foreground">Based on your preferences</p>
      </div>

      {/* Product cards */}
      <div className="space-y-3">
        {products.map((product, index) => (
          <div
            key={index}
            className="relative bg-background/90 backdrop-blur-sm rounded-xl p-3 border border-border shadow-lg animate-fade-in-up hover:border-primary/50 transition-colors cursor-pointer"
            style={{ animationDelay: `${product.delay}s` }}
          >
            <div className="flex gap-3">
              {/* Product image */}
              <div className="w-16 h-16 bg-linear-to-br from-primary/20 to-secondary/20 rounded-lg flex items-center justify-center shrink-0">
                <ShoppingCart className="w-6 h-6 text-primary/60" />
              </div>

              {/* Product info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2 mb-1">
                  <h5 className="text-xs font-semibold truncate">{product.name}</h5>
                  {index === 0 && (
                    <div className="bg-primary/20 text-primary text-[9px] font-bold px-2 py-0.5 rounded-full whitespace-nowrap">
                      {product.badge}
                    </div>
                  )}
                </div>

                {/* Rating */}
                <div className="flex items-center gap-1 mb-1.5">
                  <div className="flex gap-0.5">
                    {[...Array(5)].map((_, i) => (
                      <div
                        key={i}
                        className={`w-2 h-2 ${i < Math.floor(product.rating) ? 'text-yellow-500' : 'text-muted-foreground/30'}`}
                      >
                        ★
                      </div>
                    ))}
                  </div>
                  <span className="text-[10px] text-muted-foreground">{product.rating}</span>
                </div>

                {/* Price and action */}
                <div className="flex items-center justify-between">
                  <div className="flex items-baseline gap-1">
                    <TrendingDown className="w-3 h-3 text-green-500" />
                    <span className="text-sm font-bold text-primary">{product.price}</span>
                  </div>
                  <button className="flex items-center gap-1 text-[10px] text-primary hover:underline">
                    View
                    <ExternalLink className="w-2.5 h-2.5" />
                  </button>
                </div>
              </div>
            </div>

            {/* Highlight for best option */}
            {index === 0 && (
              <div className="absolute inset-0 bg-primary/5 rounded-xl pointer-events-none" />
            )}
          </div>
        ))}
      </div>

      {/* Action button */}
      <div className="mt-4 animate-fade-in-up" style={{ animationDelay: '0.9s' }}>
        <button className="w-full py-2.5 bg-primary text-primary-foreground rounded-lg text-xs font-semibold hover:opacity-90 transition-opacity">
          Compare All Options
        </button>
      </div>

      {/* Price savings indicator */}
      <div className="absolute top-4 right-4 bg-green-500/90 text-white text-[10px] font-bold px-2 py-1.5 rounded-full shadow-lg animate-bounce-in" style={{ animationDelay: '1.2s' }}>
        Save $100+
      </div>

      {/* Decorative elements */}
      <div className="absolute -top-4 -left-4 w-20 h-20 bg-green-500/10 rounded-full blur-2xl animate-pulse" />
      <div className="absolute -bottom-4 -right-4 w-20 h-20 bg-emerald-500/10 rounded-full blur-2xl animate-pulse" style={{ animationDelay: '1s' }} />
    </div>
  );
}

export function HowItWorksVisual() {
  const steps = [
    {
      number: 1,
      icon: MessageCircle,
      title: "Describe What You Need",
      description: "Simply chat with our AI in natural language. Tell us what you're looking for, your budget, and preferences.",
      image: "💬",
      color: "from-blue-500 to-cyan-500",
    },
    {
      number: 2,
      icon: Sparkles,
      title: "AI Analyzes & Searches",
      description: "Our advanced AI understands your needs and searches thousands of products across multiple retailers in real-time.",
      image: "🤖",
      color: "from-purple-500 to-pink-500",
    },
    {
      number: 3,
      icon: ShoppingCart,
      title: "Review & Purchase",
      description: "Get personalized recommendations with price comparisons. Ask questions, compare options, and buy with confidence.",
      image: "🛒",
      color: "from-green-500 to-emerald-500",
    },
  ];

  return (
    <section className="py-24 bg-background relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,var(--tw-gradient-stops))] from-primary/5 via-transparent to-transparent" />

      <div className="container mx-auto px-4 relative">
        <div className="text-center mb-16 space-y-4">
          <h2 className="text-4xl md:text-5xl font-bold">
            How It
            <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              Works
            </span>
          </h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Three simple steps to find your perfect product
          </p>
        </div>

        <div className="max-w-6xl mx-auto space-y-24">
          {steps.map((step, index) => {
            const Icon = step.icon;
            const isEven = index % 2 === 0;

            return (
              <div
                key={step.number}
                className={`grid lg:grid-cols-2 gap-12 items-center ${
                  !isEven ? "lg:grid-flow-dense" : ""
                }`}
              >
                {/* Content side */}
                <div className={`space-y-6 ${!isEven ? "lg:col-start-2" : ""}`}>
                  <div className="flex items-center gap-4">
                    <div className={`w-16 h-16 rounded-2xl bg-linear-to-br ${step.color} flex items-center justify-center text-white font-bold text-2xl shadow-lg`}>
                      {step.number}
                    </div>
                    <div className={`w-14 h-14 rounded-xl bg-linear-to-br ${step.color} opacity-20 flex items-center justify-center`}>
                      <Icon className="w-7 h-7 text-primary" />
                    </div>
                  </div>

                  <h3 className="text-3xl font-bold">{step.title}</h3>
                  <p className="text-lg text-muted-foreground leading-relaxed">
                    {step.description}
                  </p>

                  {index < steps.length - 1 && (
                    <div className="flex items-center gap-2 text-muted-foreground pt-4">
                      <span className="text-sm font-medium">Next step</span>
                      <ArrowRight className="w-4 h-4" />
                    </div>
                  )}
                </div>

                {/* Visual side */}
                <div className={`${!isEven ? "lg:col-start-1 lg:row-start-1" : ""}`}>
                  <div className="relative aspect-square rounded-3xl overflow-hidden border border-border">
                    {/* Render appropriate animation based on step number */}
                    {step.number === 1 ? (
                      <ChatInterfaceAnimation />
                    ) : step.number === 2 ? (
                      <AIAnalysisAnimation />
                    ) : step.number === 3 ? (
                      <ProductComparisonAnimation />
                    ) : null}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
