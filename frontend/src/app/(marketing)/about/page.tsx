import type { Metadata } from "next";
import { Sparkles, Target, Users, Zap } from "lucide-react";

export const metadata: Metadata = {
  title: "About Us - MyLittlePrice",
  description: "Learn about MyLittlePrice - Your AI-powered shopping assistant that helps you find the perfect products at the best prices.",
};

export default function AboutPage() {
  return (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto px-4 py-24">
        <div className="max-w-4xl mx-auto space-y-16">
          {/* Hero Section */}
          <div className="text-center space-y-4">
            <h1 className="text-5xl md:text-6xl font-bold bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              About MyLittlePrice
            </h1>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Revolutionizing online shopping with AI-powered product discovery
              and price intelligence
            </p>
          </div>

          {/* Mission Section */}
          <section className="space-y-6">
            <div className="flex items-center gap-3">
              <Target className="w-8 h-8 text-primary" />
              <h2 className="text-3xl font-bold">Our Mission</h2>
            </div>
            <p className="text-lg text-muted-foreground leading-relaxed">
              At MyLittlePrice, we believe that finding the perfect product
              shouldn't be overwhelming. Our mission is to simplify online
              shopping by leveraging cutting-edge AI technology to understand
              your needs, compare thousands of products, and present you with
              the best options tailored to your preferences and budget.
            </p>
          </section>

          {/* What We Do Section */}
          <section className="space-y-6">
            <div className="flex items-center gap-3">
              <Sparkles className="w-8 h-8 text-primary" />
              <h2 className="text-3xl font-bold">What We Do</h2>
            </div>
            <div className="space-y-4 text-lg text-muted-foreground leading-relaxed">
              <p>
                MyLittlePrice is an intelligent shopping assistant powered by
                advanced AI. Simply describe what you're looking for in natural
                language, and our system will:
              </p>
              <ul className="list-disc list-inside space-y-2 ml-4">
                <li>Understand your specific needs and preferences</li>
                <li>Search across multiple retailers and marketplaces</li>
                <li>Compare prices, features, and reviews</li>
                <li>Provide personalized product recommendations</li>
                <li>Answer questions about products in real-time</li>
              </ul>
              <p>
                We combine Google's Gemini AI for natural language understanding
                with comprehensive product search capabilities to deliver a
                shopping experience that feels like having a knowledgeable
                friend guide you through every purchase decision.
              </p>
            </div>
          </section>

          {/* Technology Section */}
          <section className="space-y-6">
            <div className="flex items-center gap-3">
              <Zap className="w-8 h-8 text-primary" />
              <h2 className="text-3xl font-bold">Our Technology</h2>
            </div>
            <div className="space-y-4 text-lg text-muted-foreground leading-relaxed">
              <p>
                Built with modern web technologies and powered by state-of-the-art
                AI, MyLittlePrice offers:
              </p>
              <ul className="list-disc list-inside space-y-2 ml-4">
                <li>
                  <strong>Real-time AI Conversations:</strong> Chat naturally
                  about what you need using WebSocket-powered instant messaging
                </li>
                <li>
                  <strong>Smart Product Search:</strong> Integration with Google
                  Shopping and multiple data sources for comprehensive results
                </li>
                <li>
                  <strong>Personalized Experience:</strong> Save your
                  preferences, search history, and favorite products
                </li>
                <li>
                  <strong>Multi-region Support:</strong> Shop in your preferred
                  language, currency, and country
                </li>
                <li>
                  <strong>Privacy-First:</strong> Your data is encrypted and
                  secure, with transparent privacy controls
                </li>
              </ul>
            </div>
          </section>

          {/* Values Section */}
          <section className="space-y-6">
            <div className="flex items-center gap-3">
              <Users className="w-8 h-8 text-primary" />
              <h2 className="text-3xl font-bold">Our Values</h2>
            </div>
            <div className="grid md:grid-cols-2 gap-6">
              <div className="p-6 rounded-lg border border-border bg-secondary/20">
                <h3 className="text-xl font-semibold mb-2">Transparency</h3>
                <p className="text-muted-foreground">
                  We provide clear, honest product information and never hide
                  affiliate relationships or sponsored content.
                </p>
              </div>
              <div className="p-6 rounded-lg border border-border bg-secondary/20">
                <h3 className="text-xl font-semibold mb-2">User-First</h3>
                <p className="text-muted-foreground">
                  Every feature is designed with your convenience and satisfaction
                  in mind, not just maximizing clicks or conversions.
                </p>
              </div>
              <div className="p-6 rounded-lg border border-border bg-secondary/20">
                <h3 className="text-xl font-semibold mb-2">Innovation</h3>
                <p className="text-muted-foreground">
                  We continuously improve our AI models and search algorithms to
                  provide better, faster, and more accurate recommendations.
                </p>
              </div>
              <div className="p-6 rounded-lg border border-border bg-secondary/20">
                <h3 className="text-xl font-semibold mb-2">Privacy</h3>
                <p className="text-muted-foreground">
                  Your shopping data is yours. We protect your privacy and give
                  you full control over your information.
                </p>
              </div>
            </div>
          </section>

          {/* Contact CTA */}
          <section className="text-center space-y-4 py-12">
            <h2 className="text-3xl font-bold">Get in Touch</h2>
            <p className="text-lg text-muted-foreground">
              Have questions or feedback? We'd love to hear from you.
            </p>
            <a
              href="/contact"
              className="inline-flex px-6 py-3 rounded-full text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              Contact Us
            </a>
          </section>
        </div>
      </main>
    </div>
  );
}
