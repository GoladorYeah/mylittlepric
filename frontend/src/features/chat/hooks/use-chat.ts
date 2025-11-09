import { useEffect, useRef } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useChatStore } from "@/shared/lib";
import { useAuthStore } from "@/shared/lib";
import { SessionAPI } from "@/shared/lib";
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
  syncSession: (newSessionId: string) => void;
  syncPreferences: () => void;
  syncSavedSearch: () => void;
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
    saveCurrentSearch,
    registerWebSocketSender,
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

  // Register WebSocket sender in store for realtime sync
  useEffect(() => {
    if (isConnected) {
      registerWebSocketSender(sendJsonMessage);
    } else {
      registerWebSocketSender(null);
    }
  }, [isConnected, sendJsonMessage, registerWebSocketSender]);

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

      // For authenticated users, check for active session on server
      if (accessToken) {
        // Load user preferences from server
        try {
          await store.syncPreferencesFromServer();
        } catch (error) {
          console.error("Failed to sync preferences from server:", error);
        }

        try {
          const activeSessionResponse = await SessionAPI.getActiveSession();

          if (activeSessionResponse.has_active_session && activeSessionResponse.session) {
            const serverSessionId = activeSessionResponse.session.session_id;
            const localSessionId = store.sessionId || localStorage.getItem("chat_session_id");

            // If server has a different session, ask user which one to use
            if (localSessionId && localSessionId !== serverSessionId) {
              // We have both local and server session - prefer server (most recent)
              console.log("📱 Multi-device sync: Using server session", serverSessionId);
              setSessionId(serverSessionId);
              localStorage.setItem("chat_session_id", serverSessionId);

              try {
                await loadSessionMessages(serverSessionId);
              } catch (error) {
                console.error("Failed to load server session:", error);
              }
              return;
            } else if (!localSessionId) {
              // No local session, use server session
              console.log("☁️ Restoring session from server:", serverSessionId);
              setSessionId(serverSessionId);
              localStorage.setItem("chat_session_id", serverSessionId);

              try {
                await loadSessionMessages(serverSessionId);
              } catch (error) {
                console.error("Failed to load server session:", error);
              }
              return;
            }
          } else if (store.sessionId) {
            // No active session on server, but we have local session
            // Link it to user account
            console.log("🔗 Linking local session to user account");
            try {
              await SessionAPI.linkSessionToUser(store.sessionId);
            } catch (error) {
              console.error("Failed to link session to user:", error);
            }
          }
        } catch (error) {
          console.error("Failed to check active session:", error);
          // Continue with local session logic
        }
      }

      // If store already has sessionId (restored from persist), don't reload
      if (store.sessionId && store.messages.length > 0) {
        console.log("✅ Session restored from localStorage:", store.sessionId);
        localStorage.setItem("chat_session_id", store.sessionId);
        return;
      }

      // Otherwise, check for session in localStorage or create new one
      const savedSessionId = localStorage.getItem("chat_session_id");

      if (savedSessionId && savedSessionId === store.sessionId) {
        // Session ID exists but no messages - this is a fresh session
        console.log("🆕 Fresh session:", savedSessionId);
      } else if (savedSessionId) {
        // Session ID mismatch - load from server
        setSessionId(savedSessionId);
        try {
          await loadSessionMessages(savedSessionId);
        } catch (error) {
          console.error("Failed to load messages:", error);
        }
      } else {
        // No session at all - create new one
        const newSessionId = generateId();
        setSessionId(newSessionId);
        localStorage.setItem("chat_session_id", newSessionId);
      }
    };

    initializeSession();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialSessionId, accessToken]);

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

      // Handle realtime sync messages
      if (data.type === "user_message_sync") {
        // User message from another device
        console.log("📱 Received user message sync from another device");

        // Ignore if session doesn't match
        if (data.session_id && data.session_id !== sessionId) {
          console.warn("⚠️ Ignoring sync from different session");
          return;
        }

        const userMessage = {
          id: generateId(),
          role: "user" as const,
          content: data.output || "",
          timestamp: Date.now(),
          isLocal: false, // This message is from another device
        };

        addMessage(userMessage);
        setLoading(true);
        return;
      }

      if (data.type === "assistant_message_sync") {
        // Assistant message from another device
        console.log("📱 Received assistant message sync from another device");

        // Ignore if session doesn't match
        if (data.session_id && data.session_id !== sessionId) {
          console.warn("⚠️ Ignoring sync from different session");
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
          isLocal: false, // Synced messages are not local
        };

        addMessage(assistantMessage);
        setLoading(false);

        if (data.search_state && data.search_state.category) {
          setCurrentCategory(data.search_state.category);
        }

        if (data.search_state) {
          setSearchInProgress(data.search_state.status === "completed");
        }
        return;
      }

      // Support old message_sync type for backwards compatibility
      if (data.type === "message_sync") {
        // Message from another device (legacy format - assistant only)
        console.log("📱 Received message sync from another device (legacy)");

        // Ignore if session doesn't match
        if (data.session_id && data.session_id !== sessionId) {
          console.warn("⚠️ Ignoring sync from different session");
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
          isLocal: false, // Legacy synced messages are not local
        };

        addMessage(assistantMessage);

        if (data.search_state && data.search_state.category) {
          setCurrentCategory(data.search_state.category);
        }

        if (data.search_state) {
          setSearchInProgress(data.search_state.status === "completed");
        }
        return;
      }

      if (data.type === "preferences_updated") {
        // Preferences changed on another device
        console.log("📱 Preferences updated on another device");
        const store = useChatStore.getState();
        store.syncPreferencesFromServer();
        return;
      }

      if (data.type === "saved_search_updated") {
        // Saved search changed on another device
        console.log("📱 Saved search updated on another device");
        const store = useChatStore.getState();
        store.syncPreferencesFromServer();
        return;
      }

      if (data.type === "session_changed") {
        // Session changed on another device (e.g., New Search)
        console.log("📱 Session changed on another device");
        if (data.session_id && data.session_id !== sessionId) {
          setSessionId(data.session_id);
          localStorage.setItem("chat_session_id", data.session_id);
          loadSessionMessages(data.session_id);
        }
        return;
      }

      if (data.type === "error") {
        const errorMessage = data.message || data.error || "An error occurred";
        addMessage({
          id: generateId(),
          role: "assistant",
          content: errorMessage,
          timestamp: Date.now(),
          isLocal: true, // Error responses are local to this device
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
        isLocal: true, // Direct response to local message
      };

      addMessage(assistantMessage);

      if (data.search_state && data.search_state.category) {
        setCurrentCategory(data.search_state.category);
      }

      if (data.search_state) {
        setSearchInProgress(data.search_state.status === "completed");
      }

      // Note: Search history is now managed entirely by the backend
      // History is accessible via SearchHistoryAPI.getSearchHistory()
    }
  }, [
    lastJsonMessage,
    addMessage,
    setLoading,
    setSearchInProgress,
    setCurrentCategory,
    sessionId,
    setSessionId,
    loadSessionMessages,
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
      isLocal: true, // Message sent from this device
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
        isLocal: true, // Error messages are local
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

    // Sync new session to other devices
    if (accessToken && isConnected) {
      sendJsonMessage({
        type: "sync_session",
        session_id: newSessionId,
        access_token: accessToken,
      });
    }
  };

  const syncSession = (newSessionId: string) => {
    if (accessToken && isConnected) {
      sendJsonMessage({
        type: "sync_session",
        session_id: newSessionId,
        access_token: accessToken,
      });
    }
  };

  const syncPreferences = () => {
    if (accessToken && isConnected) {
      sendJsonMessage({
        type: "sync_preferences",
        session_id: sessionId,
        access_token: accessToken,
      });
    }
  };

  const syncSavedSearch = () => {
    if (accessToken && isConnected) {
      sendJsonMessage({
        type: "sync_saved_search",
        session_id: sessionId,
        access_token: accessToken,
      });
    }
  };

  return {
    sendMessage,
    handleNewSearch,
    syncSession,
    syncPreferences,
    syncSavedSearch,
    connectionStatus,
    isConnected,
    readyState,
  };
}
