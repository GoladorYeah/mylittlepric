"use client";

import { useRef } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { ChatMessage as ChatMessageType } from "@/types";
import { ProductCard } from "./ProductCard";
import { AdPlaceholder } from "./AdPlaceholder";

interface ChatMessageProps {
  message: ChatMessageType;
  onQuickReply: (reply: string) => void;
}

// Parse quick reply to separate text and price
function parseQuickReply(reply: string): { text: string; price: string | null } {
  // Match various price patterns:
  // - "Option (≈CHF 100-200)"
  // - "Option (CHF 100–200)"
  // - "Option (≈$100)"
  // - "Option (CHF 500–1500+)"
  // - "Option (≈$100-200k)"
  // Support various dash types: - – — (hyphen, en-dash, em-dash)
  const priceMatch = reply.match(/\(([≈~]?[A-Z$€£¥]{1,4}[\s]?[\d,.\-–—]+[\+]?(?:[\s]?[kK]|[\s]?[\-–—][\s]?[\d,.\-–—]+[\+]?(?:[kK])?)?)\)$/);

  if (priceMatch) {
    const text = reply.substring(0, priceMatch.index).trim();
    const price = priceMatch[1];
    return { text, price };
  }

  return { text: reply, price: null };
}

export function ChatMessage({ message, onQuickReply }: ChatMessageProps) {
  const isUser = message.role === "user";
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollContainerRef.current) {
      const scrollAmount = 224; // Width of card (210px) + gap (14px)
      const newScrollLeft = direction === 'left'
        ? scrollContainerRef.current.scrollLeft - scrollAmount
        : scrollContainerRef.current.scrollLeft + scrollAmount;

      scrollContainerRef.current.scrollTo({
        left: newScrollLeft,
        behavior: 'smooth'
      });
    }
  };

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`${message.products && message.products.length > 0 ? 'w-full' : 'max-w-[80%]'} space-y-3 ${
          isUser ? "items-end" : "items-start"
        }`}
      >
        {message.content && message.content.trim() !== '' && (
          <div
            className={`rounded-2xl px-4 py-3 ${
              isUser
                ? "bg-secondary text-secondary-foreground"
                : "bg-card text-foreground border border-border"
            }`}
          >
            <p className="whitespace-pre-wrap">{message.content}</p>
          </div>
        )}

        {message.quick_replies && message.quick_replies.length > 0 && (
          <div className="flex flex-wrap gap-2.5 stagger-container">
            {message.quick_replies.map((reply, index) => {
              const { text, price } = parseQuickReply(reply);

              return (
                <button
                  key={index}
                  onClick={() => onQuickReply(reply)}
                  className="group relative px-4 py-2.5 rounded-xl bg-gradient-to-br from-secondary to-secondary/70 hover:from-secondary hover:to-secondary/90 text-sm transition-all duration-300 border border-border/50 hover:border-primary/30 hover:shadow-lg hover:shadow-primary/10 hover:-translate-y-1 hover:scale-[1.02] flex items-center gap-2.5 overflow-hidden elevation-transition"
                  style={{ animationDelay: `${index * 80}ms` }}
                >
                  {/* Animated gradient overlay on hover */}
                  <div className="absolute inset-0 bg-gradient-to-r from-transparent via-primary/5 to-transparent opacity-0 group-hover:opacity-100 group-hover:animate-shimmer transition-opacity duration-300" />

                  <span className="relative font-medium text-foreground group-hover:text-primary transition-colors duration-300">
                    {text}
                  </span>

                  {price && (
                    <span className="relative text-xs px-2.5 py-1 rounded-lg bg-gradient-to-br from-primary/15 to-primary/10 text-primary font-bold group-hover:from-primary/25 group-hover:to-primary/15 group-hover:scale-105 transition-all duration-300 shadow-sm border border-primary/20">
                      {price}
                    </span>
                  )}

                  {/* Subtle shine effect */}
                  <div className="absolute inset-0 bg-gradient-to-tr from-transparent via-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none" />
                </button>
              );
            })}
          </div>
        )}

        {message.products && message.products.length > 0 && (
          <div className="w-full">
            {/* Header */}
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="h-8 w-1 bg-gradient-to-b from-primary to-primary/50 rounded-full" />
                <div>
                  <h3 className="text-lg font-bold text-foreground">
                    {message.products.length} {message.products.length === 1 ? 'Product' : 'Products'} Found
                  </h3>
                  <p className="text-xs text-muted-foreground">Swipe to explore more options</p>
                </div>
              </div>

              {/* Navigation Arrows - Desktop only */}
              {message.products.length > 3 && (
                <div className="hidden md:flex items-center gap-2">
                  <button
                    onClick={() => scroll('left')}
                    className="group p-2 rounded-full bg-secondary hover:bg-primary/10 border border-border hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5"
                    aria-label="Scroll left"
                  >
                    <ChevronLeft className="w-5 h-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  </button>
                  <button
                    onClick={() => scroll('right')}
                    className="group p-2 rounded-full bg-secondary hover:bg-primary/10 border border-border hover:border-primary/30 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5"
                    aria-label="Scroll right"
                  >
                    <ChevronRight className="w-5 h-5 text-muted-foreground group-hover:text-primary transition-colors" />
                  </button>
                </div>
              )}
            </div>

            {/* Products Slider */}
            <div className="relative group/slider -mx-2">
              {/* Gradient Overlays */}
              <div className="absolute left-0 top-0 bottom-0 w-12 bg-gradient-to-r from-background to-transparent z-10 pointer-events-none opacity-0 group-hover/slider:opacity-100 transition-opacity" />
              <div className="absolute right-0 top-0 bottom-0 w-12 bg-gradient-to-l from-background to-transparent z-10 pointer-events-none opacity-0 group-hover/slider:opacity-100 transition-opacity" />

              {/* Slider Container */}
              <div
                ref={scrollContainerRef}
                className="flex gap-3.5 overflow-x-auto overflow-y-visible pb-6 pt-2 px-2 snap-x snap-mandatory hide-scrollbar"
                style={{ scrollbarWidth: 'none', msOverflowStyle: 'none' }}
              >
                {message.products?.map((product, index) => {
                  const showAdAfter = (index + 1) % 5 === 0 && index < (message.products?.length || 0) - 1;

                  return (
                    <div key={`product-${index}`} className="contents">
                      <div className="flex-none w-[210px] snap-start first:ml-1">
                        <ProductCard product={product} index={index + 1} />
                      </div>
                      {showAdAfter && (
                        <div className="flex-none w-[210px] snap-start">
                          <AdPlaceholder format="card" />
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>

              {/* Scroll Indicator */}
              <div className="flex justify-center gap-1.5 mt-3">
                {Array.from({ length: Math.min(message.products?.length || 0, 10) }).map((_, idx) => (
                  <div
                    key={idx}
                    className="h-1 rounded-full bg-muted transition-all duration-300"
                    style={{
                      width: idx < Math.min(3, message.products?.length || 0) ? '24px' : '8px',
                      opacity: idx < Math.min(3, message.products?.length || 0) ? 1 : 0.3
                    }}
                  />
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}