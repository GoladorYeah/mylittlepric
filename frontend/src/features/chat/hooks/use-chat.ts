import { useEffect, useRef } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useChatStore } from "@/shared/lib";
import { useAuthStore } from "@/shared/lib";
import { generateId } from "@/shared/lib";

/**
 * Build WebSocket URL dynamically based on current page protocol
 */
function getWebSocketUrl(token?: string | null): string {
  if (typeof window === "undefined") return "";

  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  let baseUrl: string;
  if (apiUrl) {
    const url = new URL(apiUrl);
    baseUrl = `${protocol}//${url.host}/ws`;
  } else {
    baseUrl = `${protocol}//localhost:8080/ws`;
  }

  // Add token as query parameter if available
  if (token) {
    baseUrl += `?token=${encodeURIComponent(token)}`;
  }

  console.log("🔌 WebSocket URL:", baseUrl, "(Page protocol:", window.location.protocol + ")");
  return baseUrl;
}

export interface UseChatOptions {
  initialQuery?: string;
  sessionId?: string;
  reconnectAttempts?: number;
  reconnectInterval?: number;
}

export interface UseChatReturn {
  sendMessage: (message: string) => Promise<void>;
  handleNewSearch: () => void;
  connectionStatus: string;
  isConnected: boolean;
  readyState: ReadyState;
}

/**
 * Custom hook for managing WebSocket chat connection
 * Handles connection, message sending, and session management
 */
export function useChat(options: UseChatOptions = {}): UseChatReturn {
  const {
    initialQuery,
    sessionId: initialSessionId,
    reconnectAttempts = 10,
    reconnectInterval = 3000,
  } = options;

  const initialQuerySentRef = useRef(false);
  const processedMessageIds = useRef<Set<string>>(new Set());
  const sessionLoadedRef = useRef(false);

  const {
    messages,
    sessionId,
    country,
    language,
    currency,
    currentCategory,
    _hasInitialized,
    addMessage,
    setLoading,
    setSessionId,
    setSearchInProgress,
    setCurrentCategory,
    newSearch,
    initializeLocale,
    loadSessionMessages,
    addSearchToHistory,
    saveCurrentSearch,
  } = useChatStore();

  const { accessToken } = useAuthStore();

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    getWebSocketUrl(accessToken),
    {
      shouldReconnect: () => true,
      reconnectAttempts,
      reconnectInterval,
      onOpen: () => {
        console.log("✅ WebSocket connected");
      },
      onError: (event) => {
        console.error("❌ WebSocket error:", event);
      },
      onClose: (event) => {
        console.log("🔌 WebSocket closed:", event.code, event.reason);
      },
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

  // Initialize locale on mount
  useEffect(() => {
    initializeLocale();
  }, [initializeLocale]);

  // Initialize session on mount
  useEffect(() => {
    const initializeSession = async () => {
      const store = useChatStore.getState();
      if (store._hasInitialized && !initialSessionId) {
        return;
      }

      useChatStore.setState({ _hasInitialized: true });

      // If session_id is provided in URL, use it and load messages
      if (initialSessionId && !sessionLoadedRef.current) {
        sessionLoadedRef.current = true;
        setSessionId(initialSessionId);
        localStorage.setItem("chat_session_id", initialSessionId);

        try {
          await loadSessionMessages(initialSessionId);
        } catch (error) {
          console.error("Failed to load session from URL:", error);
        }
        return;
      }

      const savedSessionId = localStorage.getItem("chat_session_id");

      if (savedSessionId) {
        setSessionId(savedSessionId);

        try {
          await loadSessionMessages(savedSessionId);
        } catch (error) {
          console.error("Failed to load messages:", error);
        }
      } else {
        const newSessionId = generateId();
        setSessionId(newSessionId);
        localStorage.setItem("chat_session_id", newSessionId);
      }
    };

    initializeSession();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialSessionId]);

  // Sync sessionId to localStorage when it changes
  useEffect(() => {
    if (sessionId && _hasInitialized) {
      const currentStoredId = localStorage.getItem("chat_session_id");
      if (currentStoredId !== sessionId) {
        localStorage.setItem("chat_session_id", sessionId);
      }
    }
  }, [sessionId, _hasInitialized]);

  // Handle incoming WebSocket messages
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

      // Ignore messages from old sessions (e.g., after New Search)
      if (data.session_id && data.session_id !== sessionId) {
        console.warn(
          "⚠️ Ignoring message from old session:",
          data.session_id,
          "Current session:",
          sessionId
        );
        return;
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

      // Save to local history only for authenticated users
      // (Backend saves history for both authenticated and anonymous users)
      if (data.products && data.products.length > 0 && accessToken) {
        const userMessages = messages.filter((m) => m.role === "user");
        const lastUserMessage = userMessages[userMessages.length - 1];

        if (lastUserMessage) {
          addSearchToHistory(
            lastUserMessage.content,
            data.search_state?.category,
            data.products.length
          );
        }
      }
    }
  }, [
    lastJsonMessage,
    addMessage,
    setLoading,
    setSearchInProgress,
    setCurrentCategory,
    sessionId,
    messages,
    addSearchToHistory,
  ]);

  // Send initial query if provided
  useEffect(() => {
    if (
      initialQuery &&
      !initialQuerySentRef.current &&
      sessionId &&
      isConnected
    ) {
      initialQuerySentRef.current = true;
      sendMessage(initialQuery);
    }
  }, [initialQuery, sessionId, isConnected]);

  const sendMessage = async (message: string) => {
    const textToSend = message.trim();
    if (!textToSend || !isConnected) return;

    const userMessage = {
      id: generateId(),
      role: "user" as const,
      content: textToSend,
      timestamp: Date.now(),
    };

    addMessage(userMessage);
    setLoading(true);

    try {
      sendJsonMessage({
        type: "chat",
        session_id: sessionId,
        message: textToSend,
        country,
        language,
        currency,
        new_search: false,
        current_category: currentCategory,
        ...(accessToken && { access_token: accessToken }),
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
    // Save current search before clearing (if there are messages)
    if (messages.length > 0) {
      saveCurrentSearch();
    }

    processedMessageIds.current.clear();
    initialQuerySentRef.current = false;
    // Don't reset _hasInitialized - locale is already initialized
    // Resetting it would trigger unnecessary re-initialization
    newSearch();
    const newSessionId = generateId();
    setSessionId(newSessionId);
    localStorage.setItem("chat_session_id", newSessionId);
  };

  return {
    sendMessage,
    handleNewSearch,
    connectionStatus,
    isConnected,
    readyState,
  };
}
