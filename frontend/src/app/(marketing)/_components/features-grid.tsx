import { Bot, Search, Shield, Sparkles, TrendingUp, Zap } from "lucide-react";

export function FeaturesGrid() {
  const features = [
    {
      icon: Bot,
      title: "AI-Powered Intelligence",
      description: "Advanced Gemini AI understands your needs in natural language and finds exactly what you're looking for.",
      gradient: "from-purple-500/20 to-pink-500/20",
    },
    {
      icon: Search,
      title: "Real-Time Product Search",
      description: "Search across thousands of retailers simultaneously with instant results and up-to-date pricing.",
      gradient: "from-blue-500/20 to-cyan-500/20",
    },
    {
      icon: TrendingUp,
      title: "Smart Price Comparison",
      description: "Automatically compare prices, features, and reviews to find the best deal every time.",
      gradient: "from-green-500/20 to-emerald-500/20",
    },
    {
      icon: Sparkles,
      title: "Personalized Recommendations",
      description: "Get tailored product suggestions based on your preferences, budget, and shopping history.",
      gradient: "from-orange-500/20 to-red-500/20",
    },
    {
      icon: Zap,
      title: "Lightning Fast Responses",
      description: "WebSocket-powered real-time chat for instant answers to all your product questions.",
      gradient: "from-yellow-500/20 to-orange-500/20",
    },
    {
      icon: Shield,
      title: "Trusted & Secure",
      description: "Shop with confidence from verified retailers. Your data is encrypted and never shared.",
      gradient: "from-indigo-500/20 to-purple-500/20",
    },
  ];

  return (
    <section className="py-24 bg-background">
      <div className="container mx-auto px-4">
        <div className="text-center max-w-3xl mx-auto mb-16 space-y-4">
          <h2 className="text-4xl md:text-5xl font-bold">
            Why Choose
            <span className="block bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              MyLittlePrice?
            </span>
          </h2>
          <p className="text-xl text-muted-foreground">
            Powerful features that make shopping smarter, faster, and more enjoyable
          </p>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <div
                key={index}
                className="group relative p-8 rounded-2xl bg-secondary/30 border border-border hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:shadow-primary/10"
              >
                {/* Background gradient */}
                <div className={`absolute inset-0 rounded-2xl bg-linear-to-br ${feature.gradient} opacity-0 group-hover:opacity-100 transition-opacity duration-300`} />

                <div className="relative space-y-4">
                  <div className="w-14 h-14 rounded-xl bg-primary/10 flex items-center justify-center group-hover:scale-110 transition-transform">
                    <Icon className="w-7 h-7 text-primary" />
                  </div>

                  <h3 className="text-xl font-semibold">{feature.title}</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed">
                    {feature.description}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
