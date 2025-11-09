import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "MyLittlePrice - Smart Shopping Assistant",
  description: "Find the best products at the best prices with AI-powered shopping assistant",
  keywords: ["shopping", "AI", "price comparison", "products", "deals"],
  openGraph: {
    title: "MyLittlePrice - Smart Shopping Assistant",
    description: "Find the best products at the best prices",
    type: "website",
    locale: "en_US",
  },
  twitter: {
    card: "summary_large_image",
    title: "MyLittlePrice - Smart Shopping Assistant",
    description: "Find the best products at the best prices",
  },
};

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
