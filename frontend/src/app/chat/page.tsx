"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { ChatInterface } from "@/components/ChatInterface";

function ChatContent() {
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get("q");

  return <ChatInterface initialQuery={initialQuery || undefined} />;
}

export default function ChatPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background" />}>
      <ChatContent />
    </Suspense>
  );
}