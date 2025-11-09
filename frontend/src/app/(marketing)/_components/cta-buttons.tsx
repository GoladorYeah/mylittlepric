"use client";

import { useRouter } from "next/navigation";

export function CTAButtons() {
  const router = useRouter();

  return (
    <div className="flex gap-4 justify-center pt-8">
      <button
        onClick={() => router.push("/chat")}
        className="px-8 py-4 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-opacity"
      >
        Start Shopping
      </button>
      <button
        onClick={() => router.push("/chat")}
        className="px-8 py-4 bg-secondary text-secondary-foreground rounded-full font-semibold hover:bg-secondary/80 transition-colors"
      >
        Learn More
      </button>
    </div>
  );
}
