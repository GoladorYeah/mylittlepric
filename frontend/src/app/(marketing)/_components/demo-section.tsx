"use client";

import { MessageSquare, Play } from "lucide-react";
import { useState } from "react";

export function DemoSection() {
  const [isPlaying, setIsPlaying] = useState(false);

  return (
    <section id="demo" className="py-24 bg-secondary/20">
      <div className="container mx-auto px-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16 space-y-4">
            <h2 className="text-4xl md:text-5xl font-bold">
              See It In
              <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Action
              </span>
            </h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Watch how MyLittlePrice helps you find the perfect product in seconds
            </p>
          </div>

          {/*
            VIDEO PLACEHOLDER:
            This should be a demo video showing:
            1. User opening chat
            2. Typing: "I need a laptop for video editing under $1500"
            3. AI analyzing and searching
            4. Results appearing with price comparison
            5. User selecting and viewing product details

            Recommended format: MP4, 1080p, 30-60 seconds
            Alternative: Use an embedded YouTube/Vimeo video
          */}
          <div className="relative aspect-video rounded-3xl bg-linear-to-br from-primary/10 via-background to-secondary/10 border-2 border-border overflow-hidden shadow-2xl">
            {!isPlaying ? (
              <>
                {/* Video thumbnail preview - first frame */}
                <div className="absolute inset-0 bg-linear-to-br from-primary/5 via-background to-secondary/10">
                  <video
                    className="w-full h-full object-cover"
                    muted
                    playsInline
                    preload="metadata"
                  >
                    <source src="/presentation.mp4#t=0.1" type="video/mp4" />
                  </video>
                </div>

                {/* Overlay gradient */}
                <div className="absolute inset-0 bg-linear-to-t from-black/60 via-black/20 to-black/40" />

                {/* Play button overlay */}
                <button
                  onClick={() => setIsPlaying(true)}
                  className="absolute inset-0 w-full h-full group"
                >
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-24 h-24 rounded-full bg-primary text-primary-foreground flex items-center justify-center group-hover:scale-110 transition-all shadow-2xl group-hover:shadow-primary/50">
                    <Play className="w-10 h-10 ml-1" fill="currentColor" />
                  </div>
                  {/* Pulse ring animation */}
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-24 h-24 rounded-full bg-primary/30 animate-pulse-scale" />
                </button>

                {/* Video info overlay */}
                <div className="absolute bottom-4 md:bottom-8 left-4 md:left-8 right-4 md:right-8 space-y-2 md:space-y-3 pointer-events-none">
                  <div className="flex items-center gap-2 md:gap-3 mb-2 md:mb-4 flex-wrap">
                    <div className="px-2 md:px-3 py-1 rounded-full bg-primary/90 text-primary-foreground text-[10px] md:text-xs font-bold backdrop-blur-sm">
                      DEMO VIDEO
                    </div>
                    <div className="px-2 md:px-3 py-1 rounded-full bg-background/80 backdrop-blur-sm text-[10px] md:text-xs font-medium">
                      See MyLittlePrice in Action
                    </div>
                  </div>

                  <div className="flex justify-start animate-fade-in-up">
                    <div className="max-w-[200px] md:max-w-xs p-2 md:p-3 rounded-xl md:rounded-2xl rounded-bl-sm bg-background/95 backdrop-blur-sm border border-border shadow-lg">
                      <p className="text-xs md:text-sm">I need a laptop for video editing under $1500</p>
                    </div>
                  </div>
                  <div className="flex justify-end animate-fade-in-up" style={{ animationDelay: '0.2s' }}>
                    <div className="max-w-[220px] md:max-w-sm p-2 md:p-3 rounded-xl md:rounded-2xl rounded-br-sm bg-primary/95 text-primary-foreground backdrop-blur-sm shadow-lg">
                      <p className="text-xs md:text-sm">I found 12 great options! Here are the top 3...</p>
                    </div>
                  </div>
                </div>
              </>
            ) : (
              <div className="absolute inset-0 bg-background">
                <video
                  controls
                  autoPlay
                  className="w-full h-full object-contain"
                  onEnded={() => setIsPlaying(false)}
                >
                  <source src="/presentation.mp4" type="video/mp4" />
                  Your browser does not support the video tag.
                </video>
              </div>
            )}
          </div>

          {/* Key features showcase below video */}
          <div className="grid md:grid-cols-3 gap-6 mt-12">
            <div className="text-center p-6 rounded-xl bg-background border border-border">
              <div className="text-3xl font-bold text-primary mb-2">&lt;3s</div>
              <div className="text-sm text-muted-foreground">Average Response Time</div>
            </div>
            <div className="text-center p-6 rounded-xl bg-background border border-border">
              <div className="text-3xl font-bold text-primary mb-2">10+</div>
              <div className="text-sm text-muted-foreground">Retailers Compared</div>
            </div>
            <div className="text-center p-6 rounded-xl bg-background border border-border">
              <div className="text-3xl font-bold text-primary mb-2">24/7</div>
              <div className="text-sm text-muted-foreground">Always Available</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
