import { useEffect, useRef } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { useChatStore } from "@/shared/lib";
import { useAuthStore } from "@/shared/lib";
import { SessionAPI } from "@/shared/lib";
import { generateId } from "@/shared/lib";
import { reconnectManager } from "@/shared/lib/reconnect-manager";

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
    checkSavedSearchPrompt,
  } = useChatStore();

  const { accessToken } = useAuthStore();

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    getWebSocketUrl(accessToken),
    {
      shouldReconnect: () => true,
      reconnectAttempts,
      reconnectInterval,
      onOpen: async () => {
        console.log("✅ WebSocket connected");

        // Recover missed messages after reconnect
        const lastTimestamp = reconnectManager.getLastMessageTimestamp();
        if (sessionId && lastTimestamp) {
          console.log("🔄 Recovering missed messages since:", lastTimestamp.toISOString());
          setLoading(true);

          try {
            const missedMessages = await reconnectManager.recoverMissedMessages(sessionId);

            // Add missed messages to store
            missedMessages.forEach((msg) => {
              addMessage({
                id: generateId(),
                role: msg.role as "user" | "assistant",
                content: msg.content,
                timestamp: msg.timestamp ? new Date(msg.timestamp).getTime() : Date.now(),
                quick_replies: msg.quick_replies,
                products: msg.products,
                search_type: msg.search_type,
                isLocal: false, // Recovered messages are not local
              });
            });

            if (missedMessages.length > 0) {
              console.log(`✅ Synced ${missedMessages.length} missed messages`);
            }
          } catch (error) {
            console.error("Failed to sync missed messages:", error);
          } finally {
            setLoading(false);
          }
        }
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
      console.log("🔧 Initializing session:", {
        hasInitialized: store._hasInitialized,
        initialSessionId,
        currentSessionId: store.sessionId,
        messageCount: store.messages.length,
      });

      if (store._hasInitialized && !initialSessionId) {
        console.log("⏭️ Session already initialized, skipping");
        return;
      }

      useChatStore.setState({ _hasInitialized: true });

      // If session_id is provided in URL, use it and load messages
      if (initialSessionId && !sessionLoadedRef.current) {
        sessionLoadedRef.current = true;
        console.log("🔗 Loading session from URL:", initialSessionId);
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

            // If server has a different session, check if we should switch
            if (localSessionId && localSessionId !== serverSessionId) {
              // Don't switch to server session if:
              // 1. We already have messages loaded (user is actively using this session)
              // 2. OR we came from a URL with a session ID (user wants to view this specific session)
              if (store.messages.length > 0 || initialSessionId) {
                console.log("⏭️ Keeping current session (has messages or from URL):", {
                  localSessionId,
                  messageCount: store.messages.length,
                  fromURL: !!initialSessionId,
                  serverSessionAvailable: serverSessionId,
                });
                // Keep using local session - user is actively working with it
                return;
              }

              // No messages in current session, safe to switch to server session
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
        console.log("✅ Session restored from localStorage:", {
          sessionId: store.sessionId,
          messageCount: store.messages.length,
          messages: store.messages.map(m => ({
            id: m.id,
            role: m.role,
            content: m.content.substring(0, 30),
          })),
        });
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

  // Check if we should show savedSearch prompt after initialization
  useEffect(() => {
    if (_hasInitialized && !initialSessionId) {
      // Small delay to let state settle
      const timer = setTimeout(() => {
        checkSavedSearchPrompt();
      }, 100);
      return () => clearTimeout(timer);
    }
  }, [_hasInitialized, initialSessionId, checkSavedSearchPrompt]);

  // Sync sessionId to localStorage when it changes
  useEffect(() => {
    if (sessionId && _hasInitialized) {
      const currentStoredId = localStorage.getItem("chat_session_id");
      if (currentStoredId !== sessionId) {
        localStorage.setItem("chat_session_id", sessionId);
      }
    }
  }, [sessionId, _hasInitialized]);

  // Get signed session ID for authenticated users
  useEffect(() => {
    const signSessionIfAuthenticated = async () => {
      // Only sign sessions for authenticated users
      if (!accessToken || !sessionId) {
        return;
      }

      // Check if we already have a valid signed session
      const store = useChatStore.getState();
      if (store.signedSessionId) {
        return;
      }

      try {
        const signedResponse = await SessionAPI.signSession(sessionId);
        console.log("🔐 Session signed:", signedResponse.signed_session_id);
        store.setSignedSessionId(signedResponse.signed_session_id);
      } catch (error) {
        console.error("Failed to sign session:", error);
        // Continue with unsigned session (backward compatible)
      }
    };

    signSessionIfAuthenticated();
  }, [accessToken, sessionId]);

  // Handle incoming WebSocket messages
  useEffect(() => {
    if (lastJsonMessage !== null) {
      const data: any = lastJsonMessage;

      console.log("📨 WebSocket message received:", {
        type: data.type,
        message_id: data.message_id,
        content: data.output?.substring(0, 50),
        session_id: data.session_id,
      });

      if (data.type === "pong") {
        return;
      }

      // Use message_id for deduplication if available, otherwise fall back to hash
      const messageId = data.message_id || JSON.stringify(data);

      if (processedMessageIds.current.has(messageId)) {
        console.log("🔄 Skipping duplicate message:", data.type, messageId);
        return;
      }

      processedMessageIds.current.add(messageId);

      // Clean up old message IDs to prevent memory leak (keep last 100)
      if (processedMessageIds.current.size > 100) {
        const idsArray = Array.from(processedMessageIds.current);
        processedMessageIds.current = new Set(idsArray.slice(-100));
      }

      setLoading(false);

      // Handle realtime sync messages
      if (data.type === "user_message_sync") {
        // User message from another device
        console.log("📱 Received user message sync from another device", {
          message_id: data.message_id,
          content: data.output?.substring(0, 50),
          session: data.session_id
        });

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
        console.log("📱 Received assistant message sync from another device", {
          message_id: data.message_id,
          content: data.output?.substring(0, 50),
          has_products: !!data.products,
          session: data.session_id
        });

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

        // Check if it's a rate limit error
        if (data.error === "rate_limit_exceeded" || errorMessage.toLowerCase().includes("rate limit exceeded")) {
          console.warn("⚠️ Rate limit exceeded:", data);

          // Parse retry_after from message if available
          const retryMatch = errorMessage.match(/retry after (\d+) seconds?/i);
          const retryAfter = retryMatch ? parseInt(retryMatch[1], 10) : 30;

          // Set rate limit state
          const expiresAt = new Date(Date.now() + retryAfter * 1000);
          const store = useChatStore.getState();
          store.setRateLimitState({
            isLimited: true,
            reason: errorMessage,
            retryAfter,
            expiresAt,
          });

          // Auto-clear after retry_after seconds
          setTimeout(() => {
            const store = useChatStore.getState();
            store.clearRateLimitState();
          }, retryAfter * 1000);

          // Don't add error message to chat (notification will be shown instead)
          return;
        }

        // Check if it's a session ownership error
        if (errorMessage.toLowerCase().includes("session ownership") || errorMessage.toLowerCase().includes("unauthorized")) {
          console.error("❌ Session ownership validation failed");

          // Clear invalid session and start fresh
          const newSessionId = generateId();
          setSessionId(newSessionId);
          const store = useChatStore.getState();
          store.setSignedSessionId(null);
          localStorage.setItem("chat_session_id", newSessionId);
          reconnectManager.reset();
          newSearch();

          // Show user-friendly error
          addMessage({
            id: generateId(),
            role: "assistant",
            content: "Your session has expired. Please start a new conversation.",
            timestamp: Date.now(),
            isLocal: true,
          });
          return;
        }

        // Regular error handling
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
    if (!textToSend || !isConnected) {
      console.warn("⚠️ Cannot send message:", { textToSend, isConnected });
      return;
    }

    const messageId = generateId();
    const userMessage = {
      id: messageId,
      role: "user" as const,
      content: textToSend,
      timestamp: Date.now(),
      isLocal: true, // Message sent from this device
      status: "pending" as const, // Mark as pending
    };

    console.log("📤 Sending user message:", {
      messageId,
      content: textToSend.substring(0, 50),
      sessionId,
    });

    addMessage(userMessage);
    setLoading(true);

    try {
      const store = useChatStore.getState();
      // Prefer signed session ID if available
      const sessionIdToSend = store.signedSessionId || sessionId;

      sendJsonMessage({
        type: "chat",
        session_id: sessionIdToSend,
        message: textToSend,
        country,
        language,
        currency,
        new_search: false,
        current_category: currentCategory,
        ...(accessToken && { access_token: accessToken }),
      });

      // Mark as sent after successful send
      store.updateMessageStatus(messageId, "sent");

      // Update reconnectManager timestamp
      reconnectManager.setLastMessageTimestamp(new Date());
    } catch (error) {
      console.error("Error sending message:", error);
      setLoading(false);

      // Mark as failed
      const store = useChatStore.getState();
      store.updateMessageStatus(messageId, "failed", "Failed to send message");

      // Show error to user
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
    // Save current search before clearing (only if user has sent messages)
    const hasUserMessages = messages.some(msg => msg.role === "user");
    if (hasUserMessages) {
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
