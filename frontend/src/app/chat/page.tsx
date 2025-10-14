"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { ChatInterface } from "@/components/ChatInterface";
import { useChatStore } from "@/lib/store";
import { generateId } from "@/lib/utils";

function ChatContent() {
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get("q");
  const { sessionId, setSessionId, addMessage } = useChatStore();
  const [hasInitialized, setHasInitialized] = useState(false);

  useEffect(() => {
    if (!sessionId) {
      setSessionId(generateId());
    }
  }, [sessionId, setSessionId]);

  useEffect(() => {
    if (initialQuery && !hasInitialized && sessionId) {
      addMessage({
        id: generateId(),
        role: "user",
        content: initialQuery,
        timestamp: Date.now(),
      });
      setHasInitialized(true);
    }
  }, [initialQuery, hasInitialized, sessionId, addMessage]);

  return <ChatInterface initialQuery={initialQuery || undefined} />;
}

export default function ChatPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background" />}>
      <ChatContent />
    </Suspense>
  );
}