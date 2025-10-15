"use client";

import { useEffect, useRef, useState } from "react";
import { Send, Sparkles, RotateCcw, Wifi, WifiOff } from "lucide-react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useChatStore } from "@/lib/store";
import { generateId } from "@/lib/utils";
import { ChatMessage as ChatMessageComponent } from "./ChatMessage";
import { ThemeToggle } from "./ThemeToggle";

const WS_URL = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").replace("http", "ws") + "/ws";

interface ChatInterfaceProps {
  initialQuery?: string;
}

export function ChatInterface({ initialQuery }: ChatInterfaceProps) {
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const initialQuerySentRef = useRef(false);
  const processedMessageIds = useRef<Set<string>>(new Set());

  const {
    messages,
    sessionId,
    isLoading,
    country,
    language,
    currentCategory,
    addMessage,
    setLoading,
    setSessionId,
    setSearchInProgress,
    setCurrentCategory,
    newSearch,
    initializeLocale,
  } = useChatStore();

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    WS_URL,
    {
      shouldReconnect: () => true,
      reconnectAttempts: 10,
      reconnectInterval: 3000,
    }
  );

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Connected",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Disconnected",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  const isConnected = readyState === ReadyState.OPEN;

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  useEffect(() => {
    initializeLocale();
  }, [initializeLocale]);

  useEffect(() => {
    if (!sessionId) {
      const savedSessionId = localStorage.getItem("chat_session_id");
      if (savedSessionId) {
        setSessionId(savedSessionId);
      } else {
        const newSessionId = generateId();
        setSessionId(newSessionId);
        localStorage.setItem("chat_session_id", newSessionId);
      }
    } else {
      localStorage.setItem("chat_session_id", sessionId);
    }
  }, [sessionId, setSessionId]);

  useEffect(() => {
    if (lastJsonMessage !== null) {
      const data: any = lastJsonMessage;
      
      if (data.type === "pong") {
        return;
      }

      const messageHash = JSON.stringify(data);
      
      if (processedMessageIds.current.has(messageHash)) {
        return;
      }

      processedMessageIds.current.add(messageHash);

      setLoading(false);

      if (data.type === "error") {
        const errorMessage = data.message || data.error || "An error occurred";
        addMessage({
          id: generateId(),
          role: "assistant",
          content: errorMessage,
          timestamp: Date.now(),
        });
        return;
      }

      if (data.session_id && data.session_id !== sessionId) {
        setSessionId(data.session_id);
        localStorage.setItem("chat_session_id", data.session_id);
      }

      const assistantMessage = {
        id: generateId(),
        role: "assistant" as const,
        content: data.output || "",
        timestamp: Date.now(),
        quick_replies: data.quick_replies,
        products: data.products,
        search_type: data.search_type,
      };

      addMessage(assistantMessage);

      if (data.search_state && data.search_state.category) {
        setCurrentCategory(data.search_state.category);
      }

      if (data.search_state) {
        setSearchInProgress(data.search_state.status === "completed");
      }
    }
  }, [lastJsonMessage, addMessage, setLoading, setSearchInProgress, setCurrentCategory, sessionId, setSessionId]);

  useEffect(() => {
    if (
      initialQuery &&
      !initialQuerySentRef.current &&
      sessionId &&
      isConnected
    ) {
      initialQuerySentRef.current = true;
      handleSend(initialQuery);
    }
  }, [initialQuery, sessionId, isConnected]);

  const handleSend = async (message?: string) => {
    const textToSend = message || input.trim();
    if (!textToSend || isLoading || !isConnected) return;

    const userMessage = {
      id: generateId(),
      role: "user" as const,
      content: textToSend,
      timestamp: Date.now(),
    };

    addMessage(userMessage);
    setInput("");
    setLoading(true);

    try {
      sendJsonMessage({
        type: "chat",
        session_id: sessionId,
        message: textToSend,
        country,
        language,
        new_search: false,
        current_category: currentCategory,
      });
    } catch (error) {
      console.error("Error sending message:", error);
      setLoading(false);
      addMessage({
        id: generateId(),
        role: "assistant",
        content: "Failed to send message. Please check your connection.",
        timestamp: Date.now(),
      });
    }
  };

  const handleNewSearch = () => {
    processedMessageIds.current.clear();
    initialQuerySentRef.current = false;
    newSearch();
    const newSessionId = generateId();
    setSessionId(newSessionId);
    localStorage.setItem("chat_session_id", newSessionId);
  };

  const handleQuickReply = (reply: string) => {
    handleSend(reply);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex flex-col h-screen bg-background">
      <header className="border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Sparkles className="w-6 h-6 text-primary" />
            <span className="text-xl font-bold">MyLittlePrice</span>
            <div className="flex items-center gap-1.5">
              {isConnected ? (
                <>
                  <Wifi className="w-4 h-4 text-green-500" />
                  <span className="text-xs text-green-500 font-medium">
                    {connectionStatus}
                  </span>
                </>
              ) : (
                <>
                  <WifiOff className="w-4 h-4 text-red-500" />
                  <span className="text-xs text-red-500 font-medium">
                    {connectionStatus}
                  </span>
                </>
              )}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={handleNewSearch}
              disabled={!isConnected}
              className="px-4 py-2 rounded-full bg-secondary hover:bg-secondary/80 transition-colors flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <RotateCcw className="w-4 h-4" />
              <span className="hidden sm:inline">New Search</span>
            </button>
            <ThemeToggle />
          </div>
        </div>
      </header>

      <div className="flex-1 overflow-y-auto">
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          {messages.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full space-y-4 text-center pt-20">
              <Sparkles className="w-16 h-16 text-primary/50" />
              <h2 className="text-2xl font-bold">What are you looking for?</h2>
              <p className="text-muted-foreground max-w-md">
                Tell me what product you need and I'll help you find the best
                options
              </p>
            </div>
          ) : (
            <div className="space-y-6">
              {messages.map((message) => (
                <ChatMessageComponent
                  key={message.id}
                  message={message}
                  onQuickReply={handleQuickReply}
                />
              ))}
              {isLoading && (
                <div className="flex justify-start">
                  <div className="bg-secondary rounded-2xl px-4 py-3">
                    <div className="flex gap-1">
                      <div className="w-2 h-2 rounded-full bg-primary animate-bounce" />
                      <div
                        className="w-2 h-2 rounded-full bg-primary animate-bounce"
                        style={{ animationDelay: "0.2s" }}
                      />
                      <div
                        className="w-2 h-2 rounded-full bg-primary animate-bounce"
                        style={{ animationDelay: "0.4s" }}
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>
      </div>

      <div className="border-t border-border bg-background">
        <div className="container mx-auto px-4 py-4 max-w-4xl">
          <div className="flex gap-2">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={
                isConnected
                  ? "Type your message..."
                  : `${connectionStatus}...`
              }
              disabled={isLoading || !isConnected}
              className="flex-1 px-4 py-3 rounded-full bg-secondary border border-border focus:border-primary focus:outline-none transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            />
            <button
              onClick={() => handleSend()}
              disabled={!input.trim() || isLoading || !isConnected}
              className="w-12 h-12 rounded-full bg-primary text-primary-foreground hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
            >
              <Send className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}