"use client";

import { LockKeyhole, Sparkles } from "lucide-react";
import { useState } from "react";
import AuthDialog from "./AuthDialog";

interface AuthPromptProps {
  searchesUsed: number;
  searchesLimit: number;
  onDismiss?: () => void;
}

export function AuthPrompt({ searchesUsed, searchesLimit, onDismiss }: AuthPromptProps) {
  const [showAuthDialog, setShowAuthDialog] = useState(false);

  return (
    <>
      <div className="flex items-center justify-center min-h-[60vh] px-4">
        <div className="max-w-lg w-full bg-card border border-border rounded-2xl shadow-xl p-8 space-y-6">
          {/* Icon */}
          <div className="flex justify-center">
            <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
              <LockKeyhole className="w-8 h-8 text-primary" />
            </div>
          </div>

          {/* Title */}
          <div className="text-center space-y-2">
            <h2 className="text-2xl font-bold text-foreground">
              Free Searches Used
            </h2>
            <p className="text-muted-foreground">
              You've completed {searchesUsed} out of {searchesLimit} free searches
            </p>
          </div>

          {/* Benefits */}
          <div className="bg-secondary/30 rounded-xl p-4 space-y-3">
            <div className="flex items-start gap-3">
              <Sparkles className="w-5 h-5 text-primary shrink-0 mt-0.5" />
              <div className="text-sm text-foreground">
                <strong>Unlimited searches</strong> with a free account
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Sparkles className="w-5 h-5 text-primary shrink-0 mt-0.5" />
              <div className="text-sm text-foreground">
                <strong>Save your search history</strong> and access it anytime
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Sparkles className="w-5 h-5 text-primary shrink-0 mt-0.5" />
              <div className="text-sm text-foreground">
                <strong>Personalized recommendations</strong> based on your preferences
              </div>
            </div>
          </div>

          {/* Buttons */}
          <div className="space-y-3">
            <button
              onClick={() => setShowAuthDialog(true)}
              className="w-full px-6 py-3 bg-primary hover:bg-primary/90 text-primary-foreground rounded-lg font-semibold transition-colors"
            >
              Sign Up / Log In
            </button>
            {onDismiss && (
              <button
                onClick={onDismiss}
                className="w-full px-6 py-3 bg-secondary hover:bg-secondary/80 text-secondary-foreground rounded-lg font-semibold transition-colors"
              >
                Maybe Later
              </button>
            )}
          </div>

          {/* Footer */}
          <p className="text-center text-xs text-muted-foreground">
            No credit card required • Free forever
          </p>
        </div>
      </div>

      <AuthDialog isOpen={showAuthDialog} onClose={() => setShowAuthDialog(false)} />
    </>
  );
}
