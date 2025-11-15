import { ArrowRight, MessageCircle, Sparkles, ShoppingCart } from "lucide-react";

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
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-primary/5 via-transparent to-transparent" />

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
                {/*
                  IMAGE/ANIMATION PLACEHOLDER:
                  Step 1: Screenshot/GIF of chat interface with user typing
                  Step 2: Animated visualization of AI searching/analyzing
                  Step 3: Product comparison view with prices

                  Recommended: Optimized PNG or short GIF animations
                */}
                <div className={`${!isEven ? "lg:col-start-1 lg:row-start-1" : ""}`}>
                  <div className={`relative aspect-square rounded-3xl bg-linear-to-br ${step.color} opacity-10 overflow-hidden border border-border`}>
                    <div className="absolute inset-0 flex items-center justify-center">
                      <div className="text-center space-y-4">
                        <div className="text-8xl">{step.image}</div>
                        <p className="text-sm text-muted-foreground px-8">
                          [IMAGE: {step.title}]
                          <br />
                          <span className="text-xs">
                            Screenshot or animation of this step
                          </span>
                        </p>
                      </div>
                    </div>

                    {/* Decorative elements */}
                    <div className="absolute top-4 right-4 w-32 h-32 rounded-full bg-primary/10 blur-3xl" />
                    <div className="absolute bottom-4 left-4 w-32 h-32 rounded-full bg-secondary/10 blur-3xl" />
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
