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
                {/* Thumbnail/Placeholder */}
                <div className="absolute inset-0 flex items-center justify-center bg-gradient-to-br from-primary/20 to-secondary/20">
                  <div className="text-center space-y-6">
                    <div className="w-24 h-24 mx-auto rounded-full bg-primary/20 backdrop-blur-sm border border-primary/30 flex items-center justify-center">
                      <MessageSquare className="w-12 h-12 text-primary" />
                    </div>
                    <div className="space-y-2">
                      <p className="text-lg font-medium">Interactive Demo Video</p>
                      <p className="text-sm text-muted-foreground">
                        [VIDEO: Full shopping experience walkthrough]
                      </p>
                    </div>
                  </div>
                </div>

                {/* Play button overlay */}
                <button
                  onClick={() => setIsPlaying(true)}
                  className="absolute inset-0 w-full h-full group"
                >
                  <div className="absolute inset-0 bg-black/20 group-hover:bg-black/30 transition-colors" />
                  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-20 h-20 rounded-full bg-primary text-primary-foreground flex items-center justify-center group-hover:scale-110 transition-transform shadow-xl">
                    <Play className="w-8 h-8 ml-1" fill="currentColor" />
                  </div>
                </button>

                {/* Demo chat messages preview */}
                <div className="absolute bottom-8 left-8 right-8 space-y-3">
                  <div className="flex justify-start">
                    <div className="max-w-xs p-3 rounded-2xl rounded-bl-sm bg-background/90 backdrop-blur-sm border border-border shadow-lg">
                      <p className="text-sm">I need a laptop for video editing under $1500</p>
                    </div>
                  </div>
                  <div className="flex justify-end">
                    <div className="max-w-sm p-3 rounded-2xl rounded-br-sm bg-primary/90 text-primary-foreground backdrop-blur-sm shadow-lg">
                      <p className="text-sm">I found 12 great options! Here are the top 3...</p>
                    </div>
                  </div>
                </div>
              </>
            ) : (
              <div className="absolute inset-0 flex items-center justify-center bg-background">
                <p className="text-muted-foreground">
                  [Embed your demo video here using video tag or iframe]
                </p>
                {/*
                  Example:
                  <video controls autoPlay className="w-full h-full">
                    <source src="/videos/demo.mp4" type="video/mp4" />
                  </video>
                */}
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
