import { HeroSection } from "./_components/hero-section-v2";
import { StatsSection } from "./_components/stats-section";
import { FeaturesGrid } from "./_components/features-grid";
import { DemoSection } from "./_components/demo-section";
import { HowItWorksVisual } from "./_components/how-it-works-visual";
import { KickstarterSection } from "./_components/kickstarter-section";
import { FinalCTA } from "./_components/final-cta";

export default function HomePage() {
  return (
    <div className="min-h-screen">
      <HeroSection />
      <StatsSection />
      <FeaturesGrid />
      <DemoSection />
      <HowItWorksVisual />
      <KickstarterSection />
      <FinalCTA />
    </div>
  );
}
