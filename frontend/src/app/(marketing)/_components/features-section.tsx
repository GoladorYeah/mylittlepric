import { Zap, Shield, TrendingUp } from "lucide-react";

export function FeaturesSection() {
  return (
    <div className="grid md:grid-cols-3 gap-8 pt-20">
      <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
        <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
          <Zap className="w-6 h-6 text-primary" />
        </div>
        <h3 className="text-lg font-semibold mb-2">Lightning Fast</h3>
        <p className="text-muted-foreground text-sm">
          Get instant product recommendations powered by advanced AI
        </p>
      </div>

      <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
        <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
          <Shield className="w-6 h-6 text-primary" />
        </div>
        <h3 className="text-lg font-semibold mb-2">Smart & Secure</h3>
        <p className="text-muted-foreground text-sm">
          Verified products from trusted merchants across Europe
        </p>
      </div>

      <div className="p-6 rounded-2xl bg-secondary/50 border border-border">
        <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 mx-auto">
          <TrendingUp className="w-6 h-6 text-primary" />
        </div>
        <h3 className="text-lg font-semibold mb-2">Best Prices</h3>
        <p className="text-muted-foreground text-sm">
          Compare prices and find the best deals automatically
        </p>
      </div>
    </div>
  );
}
