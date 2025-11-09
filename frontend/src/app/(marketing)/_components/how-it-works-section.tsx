export function HowItWorksSection() {
  return (
    <div className="pt-20 space-y-4">
      <h2 className="text-3xl font-bold">How It Works</h2>
      <div className="grid md:grid-cols-3 gap-8 pt-8 text-left">
        <div className="space-y-2">
          <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
            1
          </div>
          <h4 className="font-semibold">Tell us what you need</h4>
          <p className="text-sm text-muted-foreground">
            Describe the product you're looking for in natural language
          </p>
        </div>

        <div className="space-y-2">
          <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
            2
          </div>
          <h4 className="font-semibold">AI finds the best matches</h4>
          <p className="text-sm text-muted-foreground">
            Our AI searches thousands of products to find perfect matches
          </p>
        </div>

        <div className="space-y-2">
          <div className="w-8 h-8 rounded-full bg-primary text-primary-foreground flex items-center justify-center font-bold">
            3
          </div>
          <h4 className="font-semibold">Compare and buy</h4>
          <p className="text-sm text-muted-foreground">
            Review options, compare prices, and make your choice
          </p>
        </div>
      </div>
    </div>
  );
}
