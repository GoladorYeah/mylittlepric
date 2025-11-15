export function StatsSection() {
  const stats = [
    { value: "🚀 Beta", label: "Version Available Now", icon: "🚀" },
    { value: "AI-Powered", label: "Smart Product Search", icon: "🤖" },
    { value: "Real-Time", label: "Price Comparison", icon: "💰" },
    { value: "24/7", label: "Always Available", icon: "⚡" },
  ];

  return (
    <section className="py-20 bg-secondary/20 border-y border-border">
      <div className="container mx-auto px-4">
        <div className="grid md:grid-cols-4 gap-8">
          {stats.map((stat, index) => (
            <div key={index} className="text-center space-y-2">
              <div className="text-4xl mb-2">{stat.icon}</div>
              <div className="text-3xl md:text-4xl font-bold bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                {stat.value}
              </div>
              <div className="text-sm text-muted-foreground">{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
