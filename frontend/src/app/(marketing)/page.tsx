import { HeroSection } from "./_components/hero-section";
import { CTAButtons } from "./_components/cta-buttons";
import { FeaturesSection } from "./_components/features-section";
import { HowItWorksSection } from "./_components/how-it-works-section";
import { Footer } from "./_components/footer";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-linear-to-b from-background to-muted/20">
      <HeroSection />

      <main className="pt-32 pb-20">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center space-y-8">
            <div className="space-y-4">
              <h1 className="text-5xl md:text-6xl font-bold bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Find the Perfect Product
              </h1>
              <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
                AI-powered shopping assistant that helps you discover the best
                products at the best prices
              </p>
            </div>

            <CTAButtons />
            <FeaturesSection />
            <HowItWorksSection />
          </div>
        </div>
      </main>

      <Footer />
    </div>
  );
}
