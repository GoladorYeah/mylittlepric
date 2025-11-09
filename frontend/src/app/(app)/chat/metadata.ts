import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Chat - AI Shopping Assistant",
  description:
    "Chat with our AI shopping assistant to find the perfect products. Get personalized recommendations and compare prices instantly.",
  openGraph: {
    title: "Chat - AI Shopping Assistant | MyLittlePrice",
    description:
      "Chat with our AI shopping assistant to find the perfect products",
  },
  robots: {
    index: false, // Don't index chat pages
    follow: true,
  },
};
