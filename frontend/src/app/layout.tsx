import type { Metadata } from "next";
import { Geist, Geist_Mono, Space_Grotesk } from "next/font/google";
import { Providers } from "@/shared/lib/providers";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const spaceGrotesk = Space_Grotesk({
  variable: "--font-space-grotesk",
  subsets: ["latin"],
  weight: ["700"],
});

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000"),
  title: {
    default: "MyLittlePrice - Smart Shopping Assistant",
    template: "%s | MyLittlePrice",
  },
  description:
    "Find the best products at the best prices with AI-powered shopping assistant. Compare prices, discover deals, and shop smarter.",
  keywords: [
    "shopping",
    "price comparison",
    "AI assistant",
    "product search",
    "best deals",
    "online shopping",
    "smart shopping",
  ],
  authors: [{ name: "MyLittlePrice" }],
  creator: "MyLittlePrice",
  publisher: "MyLittlePrice",
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "/",
    siteName: "MyLittlePrice",
    title: "MyLittlePrice - Smart Shopping Assistant",
    description:
      "Find the best products at the best prices with AI-powered shopping assistant",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
        alt: "MyLittlePrice - Smart Shopping Assistant",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "MyLittlePrice - Smart Shopping Assistant",
    description:
      "Find the best products at the best prices with AI-powered shopping assistant",
    images: ["/og-image.png"],
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  verification: {
    // google: "your-google-verification-code",
    // yandex: "your-yandex-verification-code",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        {/* Theme color for mobile browsers - Light mode */}
        <meta name="theme-color" content="#fafafa" media="(prefers-color-scheme: light)" />
        {/* Theme color for mobile browsers - Dark mode */}
        <meta name="theme-color" content="#1a1d26" media="(prefers-color-scheme: dark)" />
        {/* Color scheme support for older browsers */}
        <meta name="color-scheme" content="light dark" />
        {/* iOS Safari specific */}
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${spaceGrotesk.variable} antialiased`}
      >
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
