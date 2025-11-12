# Frontend Improvements Required for Backend Changes

**Дата:** 12 ноября 2024
**Проект:** MyLittlePrice - AI Shopping Assistant
**Основано на:** BACKEND_IMPROVEMENTS.md (v2.2.0)

---

## 📋 Оглавление

- [Текущее состояние Frontend](#текущее-состояние-frontend)
- [Необходимые изменения](#необходимые-изменения)
  - [Приоритет 1: Критично для Production](#приоритет-1-критично-для-production)
  - [Приоритет 2: Улучшение UX](#приоритет-2-улучшение-ux)
  - [Приоритет 3: Nice to have](#приоритет-3-nice-to-have)
- [Детальные спецификации](#детальные-спецификации)
- [План реализации](#план-реализации)

---

## 📊 Текущее состояние Frontend

### ✅ Что уже реализовано:

1. **Базовый WebSocket connection** (`use-chat.ts`)
   - Подключение через `react-use-websocket`
   - Автоматический reconnect (10 попыток, 3 сек интервал)
   - Message deduplication
   - Ping/Pong handling

2. **Session Management** (`session-api.ts`)
   - `getActiveSession()` - получение активной сессии
   - `linkSessionToUser()` - связывание сессии с пользователем

3. **Multi-device Sync**
   - Обработка `user_message_sync` / `assistant_message_sync`
   - Обработка `preferences_updated` / `saved_search_updated`
   - Обработка `session_changed`

4. **Error Handling**
   - Базовая обработка WebSocket ошибок
   - Обработка `type: "error"` messages

### ❌ Что отсутствует:

1. **Reconnect с восстановлением пропущенных сообщений**
   - Нет вызова endpoint `/api/chat/messages/since`
   - Не хранится timestamp последнего сообщения

2. **Rate Limiting Error Handling**
   - Нет обработки `rate_limit_exceeded` errors
   - Нет UI feedback для пользователя
   - Нет retry logic с учетом `retry_after`

3. **Session Ownership Validation**
   - Не используются signed session IDs
   - Нет вызова endpoint `POST /api/sessions/sign`

4. **Rate Limit Headers Tracking**
   - Не парсятся `X-RateLimit-*` headers
   - Нет превентивного блокирования

---

## 🚀 Необходимые изменения

### Приоритет 1: Критично для Production

#### 1. ✨ Reconnect Mechanism с восстановлением сообщений

**Проблема:**
При разрыве WebSocket соединения клиент пропускает сообщения, отправленные во время disconnect. После reconnect истории сообщений не синхронизируется.

**Решение:**

**Новые файлы:**
- `frontend/src/shared/lib/reconnect-manager.ts` - менеджер для reconnect логики

**Обновить файлы:**
- `frontend/src/shared/lib/api.ts` - добавить метод `getMessagesSince()`
- `frontend/src/features/chat/hooks/use-chat.ts` - интегрировать reconnect logic
- `frontend/src/shared/lib/store.ts` - хранить timestamp последнего сообщения

**Функционал:**

```typescript
// api.ts
export interface MessagesSinceResponse {
  messages: Array<{
    role: string;
    content: string;
    timestamp: string;
    quick_replies?: string[];
    products?: any[];
    search_type?: string;
  }>;
  session_id: string;
  message_count: number;
  since: string;
}

export async function getMessagesSince(
  sessionId: string,
  since: Date
): Promise<MessagesSinceResponse> {
  const accessToken = useAuthStore.getState().accessToken;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (accessToken) {
    headers["Authorization"] = `Bearer ${accessToken}`;
  }

  const sinceISO = since.toISOString();
  const url = `${API_URL}/api/chat/messages/since?session_id=${encodeURIComponent(sessionId)}&since=${encodeURIComponent(sinceISO)}`;

  const response = await fetch(url, {
    method: "GET",
    headers,
  });

  if (!response.ok) {
    throw new Error("Failed to fetch messages since timestamp");
  }

  return response.json();
}
```

```typescript
// reconnect-manager.ts
export class ReconnectManager {
  private lastMessageTimestamp: Date | null = null;
  private isRecovering = false;

  setLastMessageTimestamp(timestamp: Date) {
    this.lastMessageTimestamp = timestamp;
  }

  async recoverMissedMessages(sessionId: string): Promise<any[]> {
    if (!this.lastMessageTimestamp || this.isRecovering) {
      return [];
    }

    this.isRecovering = true;

    try {
      const response = await getMessagesSince(sessionId, this.lastMessageTimestamp);
      console.log(`🔄 Recovered ${response.message_count} missed messages`);
      return response.messages;
    } catch (error) {
      console.error("Failed to recover missed messages:", error);
      return [];
    } finally {
      this.isRecovering = false;
    }
  }
}
```

```typescript
// use-chat.ts - добавить в useWebSocket onOpen
onOpen: async () => {
  console.log("✅ WebSocket connected");

  // Recover missed messages after reconnect
  if (sessionId && reconnectManager.lastMessageTimestamp) {
    setLoading(true);
    try {
      const missedMessages = await reconnectManager.recoverMissedMessages(sessionId);

      // Add missed messages to store
      missedMessages.forEach(msg => {
        addMessage({
          id: generateId(),
          role: msg.role as "user" | "assistant",
          content: msg.content,
          timestamp: new Date(msg.timestamp).getTime(),
          quick_replies: msg.quick_replies,
          products: msg.products,
          search_type: msg.search_type,
          isLocal: false, // Recovered messages are not local
        });
      });

      console.log(`✅ Synced ${missedMessages.length} missed messages`);
    } catch (error) {
      console.error("Failed to sync missed messages:", error);
    } finally {
      setLoading(false);
    }
  }
}
```

**Обновить store:**

```typescript
// store.ts - добавить в ChatStore
interface ChatStore {
  // ... existing fields
  lastMessageTimestamp: Date | null;
  setLastMessageTimestamp: (timestamp: Date) => void;
}

// При добавлении сообщения обновлять timestamp
addMessage: (message) => {
  set(state => ({
    messages: [...state.messages, message],
    lastMessageTimestamp: new Date(),
  }));
}
```

**Результат:**
- ✅ Клиент восстанавливает пропущенные сообщения после reconnect
- ✅ Работает при нестабильной сети
- ✅ Синхронизация с точностью до миллисекунды

---

#### 2. ✨ Rate Limiting Error Handling

**Проблема:**
Backend отправляет `rate_limit_exceeded` ошибки, но frontend не обрабатывает их специальным образом. Пользователь не знает, что он превысил лимит и когда сможет снова отправлять сообщения.

**Решение:**

**Новые файлы:**
- `frontend/src/features/chat/components/RateLimitNotification.tsx` - UI компонент для уведомлений

**Обновить файлы:**
- `frontend/src/features/chat/hooks/use-chat.ts` - обработка rate limit errors
- `frontend/src/features/chat/components/ChatInterface.tsx` - показ уведомлений
- `frontend/src/shared/lib/store.ts` - хранение rate limit state

**Функционал:**

```typescript
// store.ts - добавить в ChatStore
interface ChatStore {
  // ... existing fields
  rateLimitState: {
    isLimited: boolean;
    reason: string | null;
    retryAfter: number | null; // seconds
    expiresAt: Date | null;
  };
  setRateLimitState: (state: Partial<RateLimitState>) => void;
  clearRateLimitState: () => void;
}
```

```typescript
// use-chat.ts - добавить обработку в useEffect для lastJsonMessage
if (data.type === "error") {
  const errorMessage = data.message || data.error || "An error occurred";

  // Check if it's a rate limit error
  if (data.error === "rate_limit_exceeded" || errorMessage.includes("rate limit exceeded")) {
    console.warn("⚠️ Rate limit exceeded:", data);

    // Parse retry_after from message if available
    const retryMatch = errorMessage.match(/retry after (\d+) seconds?/i);
    const retryAfter = retryMatch ? parseInt(retryMatch[1], 10) : 30;

    // Set rate limit state
    const expiresAt = new Date(Date.now() + retryAfter * 1000);
    useChatStore.getState().setRateLimitState({
      isLimited: true,
      reason: errorMessage,
      retryAfter,
      expiresAt,
    });

    // Auto-clear after retry_after seconds
    setTimeout(() => {
      useChatStore.getState().clearRateLimitState();
    }, retryAfter * 1000);

    // Don't add error message to chat (show notification instead)
    return;
  }

  // Regular error handling
  addMessage({
    id: generateId(),
    role: "assistant",
    content: errorMessage,
    timestamp: Date.now(),
    isLocal: true,
  });
  return;
}
```

```typescript
// RateLimitNotification.tsx
import { useEffect, useState } from "react";
import { useChatStore } from "@/shared/lib";

export function RateLimitNotification() {
  const { rateLimitState } = useChatStore();
  const [timeRemaining, setTimeRemaining] = useState<number>(0);

  useEffect(() => {
    if (!rateLimitState.isLimited || !rateLimitState.expiresAt) {
      return;
    }

    const interval = setInterval(() => {
      const now = Date.now();
      const remaining = Math.max(0, Math.floor((rateLimitState.expiresAt.getTime() - now) / 1000));
      setTimeRemaining(remaining);

      if (remaining === 0) {
        clearInterval(interval);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [rateLimitState.isLimited, rateLimitState.expiresAt]);

  if (!rateLimitState.isLimited) {
    return null;
  }

  return (
    <div className="fixed top-4 right-4 z-50 max-w-md p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg shadow-lg">
      <div className="flex items-start gap-3">
        <div className="flex-shrink-0">
          <svg className="w-5 h-5 text-yellow-600 dark:text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
        </div>
        <div className="flex-1">
          <h3 className="text-sm font-medium text-yellow-800 dark:text-yellow-200">
            Rate Limit Exceeded
          </h3>
          <p className="mt-1 text-sm text-yellow-700 dark:text-yellow-300">
            {rateLimitState.reason || "You've sent too many messages. Please wait before sending more."}
          </p>
          {timeRemaining > 0 && (
            <p className="mt-2 text-sm font-semibold text-yellow-800 dark:text-yellow-200">
              Retry in {timeRemaining} seconds
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
```

```typescript
// ChatInterface.tsx - добавить компонент
import { RateLimitNotification } from "./RateLimitNotification";

export function ChatInterface() {
  return (
    <>
      <RateLimitNotification />
      {/* ... existing chat interface */}
    </>
  );
}
```

```typescript
// chat-input.tsx - блокировать ввод при rate limit
const { rateLimitState } = useChatStore();
const isDisabled = rateLimitState.isLimited || !isConnected || loading;

// Update textarea
<textarea
  disabled={isDisabled}
  placeholder={
    rateLimitState.isLimited
      ? `Rate limit exceeded. Retry in ${Math.ceil((rateLimitState.expiresAt.getTime() - Date.now()) / 1000)}s`
      : "Ask me anything..."
  }
  // ... rest of props
/>
```

**Результат:**
- ✅ Пользователь видит notification при превышении лимита
- ✅ Countdown timer показывает время до разблокировки
- ✅ Поле ввода автоматически блокируется
- ✅ Автоматическая разблокировка после истечения времени

---

#### 3. ✨ Session Ownership Validation (Signed Sessions)

**Проблема:**
Backend реализовал HMAC-подпись session IDs для защиты от session hijacking, но frontend не использует signed session IDs.

**Решение:**

**Обновить файлы:**
- `frontend/src/shared/lib/session-api.ts` - добавить метод `signSession()`
- `frontend/src/features/chat/hooks/use-chat.ts` - использовать signed sessions
- `frontend/src/shared/lib/store.ts` - хранить signed session ID

**Функционал:**

```typescript
// session-api.ts - добавить метод
export interface SignSessionResponse {
  signed_session_id: string;
  session_id: string;
  expires_at: string;
}

export class SessionAPI {
  // ... existing methods

  /**
   * Get signed session ID for additional security
   * Signed sessions are protected with HMAC and include ownership validation
   */
  static async signSession(sessionId: string): Promise<SignSessionResponse> {
    const url = `${API_URL}/api/sessions/sign`;
    const response = await fetchWithAuth(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ session_id: sessionId }),
    });

    if (!response.ok) {
      let errorMessage = 'Failed to sign session';
      try {
        const error = await response.json();
        errorMessage = error.message || error.error || errorMessage;
      } catch (e) {
        // If response is not JSON, use default error message
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }
}
```

```typescript
// store.ts - добавить в ChatStore
interface ChatStore {
  // ... existing fields
  signedSessionId: string | null;
  setSignedSessionId: (signedSessionId: string | null) => void;
}
```

```typescript
// use-chat.ts - получить signed session ID после авторизации
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
      setSignedSessionId(signedResponse.signed_session_id);
    } catch (error) {
      console.error("Failed to sign session:", error);
      // Continue with unsigned session (backward compatible)
    }
  };

  signSessionIfAuthenticated();
}, [accessToken, sessionId]);
```

```typescript
// use-chat.ts - использовать signed session ID в WebSocket messages
const sendMessage = async (message: string) => {
  const textToSend = message.trim();
  if (!textToSend || !isConnected) return;

  // ... existing code for userMessage and addMessage

  try {
    const store = useChatStore.getState();
    const sessionIdToSend = store.signedSessionId || sessionId; // Prefer signed

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
  } catch (error) {
    // ... existing error handling
  }
};
```

```typescript
// api.ts - использовать signed session ID в HTTP запросах
export async function getSessionMessages(
  sessionId: string
): Promise<SessionMessagesResponse> {
  const accessToken = useAuthStore.getState().accessToken;
  const signedSessionId = useChatStore.getState().signedSessionId;

  // Prefer signed session ID if available
  const sessionIdToUse = signedSessionId || sessionId;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (accessToken) {
    headers["Authorization"] = `Bearer ${accessToken}`;
  }

  const response = await fetch(
    `${API_URL}/api/chat/messages?session_id=${encodeURIComponent(sessionIdToUse)}`,
    {
      method: "GET",
      headers,
    }
  );

  // ... existing error handling
}
```

**Обработка ошибок валидации:**

```typescript
// use-chat.ts - обработка ownership errors
if (data.type === "error") {
  const errorMessage = data.message || data.error || "An error occurred";

  // Check if it's a session ownership error
  if (errorMessage.includes("session ownership") || errorMessage.includes("unauthorized")) {
    console.error("❌ Session ownership validation failed");

    // Clear invalid session and start fresh
    const newSessionId = generateId();
    setSessionId(newSessionId);
    setSignedSessionId(null);
    localStorage.setItem("chat_session_id", newSessionId);
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

  // ... existing error handling
}
```

**Результат:**
- ✅ Защита от session hijacking для authenticated users
- ✅ HMAC-подпись с таймаутом (24 часа)
- ✅ Backward compatible (работает с обычными session IDs)
- ✅ Автоматическое обновление signed session ID

---

### Приоритет 2: Улучшение UX

#### 4. 📊 Rate Limit Headers Tracking

**Проблема:**
Backend отправляет `X-RateLimit-*` headers в HTTP ответах, но frontend не использует эту информацию для превентивного блокирования.

**Решение:**

**Новые файлы:**
- `frontend/src/shared/lib/rate-limit-tracker.ts` - трекинг rate limits

**Функционал:**

```typescript
// rate-limit-tracker.ts
export interface RateLimitInfo {
  limit: number;
  remaining: number;
  reset: Date;
  percentage: number; // 0-100
}

export class RateLimitTracker {
  private info: RateLimitInfo | null = null;
  private listeners: Set<(info: RateLimitInfo | null) => void> = new Set();

  updateFromHeaders(headers: Headers) {
    const limit = headers.get("X-RateLimit-Limit");
    const remaining = headers.get("X-RateLimit-Remaining");
    const reset = headers.get("X-RateLimit-Reset");

    if (limit && remaining && reset) {
      this.info = {
        limit: parseInt(limit, 10),
        remaining: parseInt(remaining, 10),
        reset: new Date(parseInt(reset, 10) * 1000),
        percentage: (parseInt(remaining, 10) / parseInt(limit, 10)) * 100,
      };

      this.notifyListeners();
    }
  }

  subscribe(listener: (info: RateLimitInfo | null) => void) {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  private notifyListeners() {
    this.listeners.forEach(listener => listener(this.info));
  }

  getInfo(): RateLimitInfo | null {
    return this.info;
  }

  isNearLimit(threshold = 10): boolean {
    return this.info ? this.info.percentage < threshold : false;
  }
}

export const rateLimitTracker = new RateLimitTracker();
```

```typescript
// api.ts - обновить все fetch calls
export async function getSessionMessages(
  sessionId: string
): Promise<SessionMessagesResponse> {
  // ... existing code

  const response = await fetch(url, { method: "GET", headers });

  // Track rate limit headers
  rateLimitTracker.updateFromHeaders(response.headers);

  if (!response.ok) {
    throw new Error("Failed to fetch session messages");
  }

  return response.json();
}
```

```typescript
// RateLimitIndicator.tsx - новый компонент
import { useEffect, useState } from "react";
import { rateLimitTracker, RateLimitInfo } from "@/shared/lib/rate-limit-tracker";

export function RateLimitIndicator() {
  const [info, setInfo] = useState<RateLimitInfo | null>(null);

  useEffect(() => {
    return rateLimitTracker.subscribe(setInfo);
  }, []);

  if (!info || info.percentage > 20) {
    return null; // Only show when less than 20% remaining
  }

  const color = info.percentage < 10 ? "text-red-500" : "text-yellow-500";

  return (
    <div className={`text-xs ${color} px-2 py-1`}>
      {info.remaining}/{info.limit} requests remaining
    </div>
  );
}
```

**Результат:**
- ✅ Tracking rate limit headers из API ответов
- ✅ UI indicator при приближении к лимиту
- ✅ Превентивное предупреждение пользователя

---

#### 5. 🔌 Improved Connection Status Indicators

**Проблема:**
Текущий индикатор статуса подключения слишком простой. Не показывает состояние "Syncing missed messages..." или "Reconnecting...".

**Решение:**

**Обновить файлы:**
- `frontend/src/features/chat/components/chat-header.tsx` - улучшенный status indicator

**Функционал:**

```typescript
// chat-header.tsx
export function ChatHeader() {
  const { connectionStatus, isConnected } = useChat();
  const { loading, rateLimitState } = useChatStore();
  const [isSyncing, setIsSyncing] = useState(false);

  // Detect syncing state
  useEffect(() => {
    if (isConnected && loading) {
      // Check if we're syncing (vs regular loading)
      const store = useChatStore.getState();
      setIsSyncing(store.lastMessageTimestamp !== null);
    } else {
      setIsSyncing(false);
    }
  }, [isConnected, loading]);

  const getStatusColor = () => {
    if (rateLimitState.isLimited) return "bg-yellow-500";
    if (isConnected) return "bg-green-500";
    if (connectionStatus === "Connecting") return "bg-yellow-500";
    return "bg-red-500";
  };

  const getStatusText = () => {
    if (rateLimitState.isLimited) {
      return `Rate limited (${Math.ceil((rateLimitState.expiresAt.getTime() - Date.now()) / 1000)}s)`;
    }
    if (isSyncing) {
      return "Syncing missed messages...";
    }
    return connectionStatus;
  };

  return (
    <header className="flex items-center justify-between p-4 border-b">
      {/* ... existing header content */}

      <div className="flex items-center gap-2">
        <div className={`w-2 h-2 rounded-full ${getStatusColor()} animate-pulse`} />
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {getStatusText()}
        </span>
      </div>
    </header>
  );
}
```

**Результат:**
- ✅ Четкие индикаторы статуса (Connected, Reconnecting, Syncing, Rate limited)
- ✅ Анимация для визуального feedback
- ✅ Countdown для rate limit и syncing

---

#### 6. ⚡ Optimistic Updates с Rollback

**Проблема:**
Сообщения пользователя уже показываются сразу (optimistic updates), но нет rollback при ошибках отправки.

**Решение:**

**Обновить файлы:**
- `frontend/src/shared/lib/store.ts` - добавить pending status
- `frontend/src/features/chat/hooks/use-chat.ts` - rollback logic
- `frontend/src/features/chat/components/ChatMessage.tsx` - показ pending state

**Функционал:**

```typescript
// store.ts - добавить pending status
export interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: number;
  isLocal?: boolean;
  quick_replies?: string[];
  products?: any[];
  search_type?: string;

  // New fields for optimistic updates
  status?: "pending" | "sent" | "failed";
  error?: string;
}

interface ChatStore {
  // ... existing fields
  updateMessageStatus: (messageId: string, status: "pending" | "sent" | "failed", error?: string) => void;
  removeMessage: (messageId: string) => void;
}

// Implementation
updateMessageStatus: (messageId, status, error) => {
  set(state => ({
    messages: state.messages.map(msg =>
      msg.id === messageId ? { ...msg, status, error } : msg
    ),
  }));
},

removeMessage: (messageId) => {
  set(state => ({
    messages: state.messages.filter(msg => msg.id !== messageId),
  }));
},
```

```typescript
// use-chat.ts - обновить sendMessage
const sendMessage = async (message: string) => {
  const textToSend = message.trim();
  if (!textToSend || !isConnected) return;

  const messageId = generateId();
  const userMessage = {
    id: messageId,
    role: "user" as const,
    content: textToSend,
    timestamp: Date.now(),
    isLocal: true,
    status: "pending" as const, // Mark as pending
  };

  addMessage(userMessage);
  setLoading(true);

  try {
    const store = useChatStore.getState();
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
      isLocal: true,
    });
  }
};
```

```typescript
// ChatMessage.tsx - показать pending/failed state
export function ChatMessage({ message }: { message: Message }) {
  const isUser = message.role === "user";
  const isPending = message.status === "pending";
  const isFailed = message.status === "failed";

  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"} mb-4`}>
      <div className={`max-w-[80%] ${isUser ? "bg-blue-500 text-white" : "bg-gray-100 dark:bg-gray-800"} rounded-lg p-3`}>
        <p>{message.content}</p>

        {/* Show status indicators */}
        {isPending && (
          <div className="flex items-center gap-1 mt-1 text-xs opacity-70">
            <span className="animate-pulse">Sending...</span>
          </div>
        )}

        {isFailed && (
          <div className="flex items-center gap-1 mt-1 text-xs text-red-300">
            <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <span>Failed to send</span>
            <button
              onClick={() => handleRetry(message)}
              className="underline ml-1"
            >
              Retry
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
```

**Результат:**
- ✅ Визуальный feedback для pending messages
- ✅ Показ ошибки отправки
- ✅ Кнопка "Retry" для failed messages
- ✅ Улучшенный UX при проблемах с сетью

---

### Приоритет 3: Nice to have

#### 7. 📈 WebSocket Metrics на Frontend

**Проблема:**
Нет visibility в производительность WebSocket соединения на стороне клиента.

**Решение:**

**Новые файлы:**
- `frontend/src/shared/lib/ws-metrics.ts` - метрики WebSocket

**Функционал:**

```typescript
// ws-metrics.ts
export interface WSMetrics {
  connectionCount: number;
  messagesSent: number;
  messagesReceived: number;
  errors: number;
  averageLatency: number; // ms
  lastConnectedAt: Date | null;
  totalUptime: number; // ms
}

export class WSMetricsTracker {
  private metrics: WSMetrics = {
    connectionCount: 0,
    messagesSent: 0,
    messagesReceived: 0,
    errors: 0,
    averageLatency: 0,
    lastConnectedAt: null,
    totalUptime: 0,
  };

  private latencies: number[] = [];
  private connectedAt: Date | null = null;

  onConnect() {
    this.metrics.connectionCount++;
    this.connectedAt = new Date();
    this.metrics.lastConnectedAt = this.connectedAt;
  }

  onDisconnect() {
    if (this.connectedAt) {
      const uptime = Date.now() - this.connectedAt.getTime();
      this.metrics.totalUptime += uptime;
      this.connectedAt = null;
    }
  }

  onMessageSent() {
    this.metrics.messagesSent++;
  }

  onMessageReceived(latency?: number) {
    this.metrics.messagesReceived++;

    if (latency !== undefined) {
      this.latencies.push(latency);
      if (this.latencies.length > 100) {
        this.latencies.shift();
      }

      this.metrics.averageLatency =
        this.latencies.reduce((a, b) => a + b, 0) / this.latencies.length;
    }
  }

  onError() {
    this.metrics.errors++;
  }

  getMetrics(): WSMetrics {
    return { ...this.metrics };
  }

  reset() {
    this.metrics = {
      connectionCount: 0,
      messagesSent: 0,
      messagesReceived: 0,
      errors: 0,
      averageLatency: 0,
      lastConnectedAt: null,
      totalUptime: 0,
    };
    this.latencies = [];
  }
}

export const wsMetrics = new WSMetricsTracker();
```

```typescript
// use-chat.ts - интегрировать metrics
const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
  getWebSocketUrl(accessToken),
  {
    shouldReconnect: () => true,
    reconnectAttempts,
    reconnectInterval,
    onOpen: () => {
      console.log("✅ WebSocket connected");
      wsMetrics.onConnect();
      // ... existing onOpen logic
    },
    onError: (event) => {
      console.error("❌ WebSocket error:", event);
      wsMetrics.onError();
    },
    onClose: (event) => {
      console.log("🔌 WebSocket closed:", event.code, event.reason);
      wsMetrics.onDisconnect();
    },
    onMessage: () => {
      wsMetrics.onMessageReceived();
    },
  }
);

// Track sent messages
const sendMessage = async (message: string) => {
  // ... existing code

  sendJsonMessage({
    // ... message data
  });

  wsMetrics.onMessageSent();
};
```

**Результат:**
- ✅ Tracking connection metrics
- ✅ Latency monitoring
- ✅ Error rate tracking
- ✅ Данные для debugging и optimization

---

#### 8. 🔄 Advanced Retry Logic

**Проблема:**
Текущий reconnect использует фиксированный интервал (3 секунды). При частых разрывах может быть неэффективным.

**Решение:**

**Новые файлы:**
- `frontend/src/shared/lib/reconnect-strategy.ts` - продвинутая логика reconnect

**Функционал:**

```typescript
// reconnect-strategy.ts
export class ExponentialBackoffStrategy {
  private attempt = 0;
  private maxAttempts = 10;
  private baseDelay = 1000; // 1 second
  private maxDelay = 30000; // 30 seconds

  getNextDelay(): number {
    const delay = Math.min(
      this.baseDelay * Math.pow(2, this.attempt),
      this.maxDelay
    );

    // Add jitter (±20%)
    const jitter = delay * 0.2 * (Math.random() * 2 - 1);

    this.attempt++;
    return Math.floor(delay + jitter);
  }

  reset() {
    this.attempt = 0;
  }

  shouldRetry(): boolean {
    return this.attempt < this.maxAttempts;
  }
}

export class CircuitBreaker {
  private failureCount = 0;
  private failureThreshold = 5;
  private resetTimeout = 60000; // 1 minute
  private state: "closed" | "open" | "half-open" = "closed";
  private resetTimer: NodeJS.Timeout | null = null;

  recordSuccess() {
    this.failureCount = 0;
    this.state = "closed";

    if (this.resetTimer) {
      clearTimeout(this.resetTimer);
      this.resetTimer = null;
    }
  }

  recordFailure() {
    this.failureCount++;

    if (this.failureCount >= this.failureThreshold) {
      this.state = "open";
      console.warn("🔴 Circuit breaker opened - too many failures");

      // Try to reset after timeout
      this.resetTimer = setTimeout(() => {
        this.state = "half-open";
        this.failureCount = 0;
        console.log("🟡 Circuit breaker half-open - attempting reconnect");
      }, this.resetTimeout);
    }
  }

  canAttempt(): boolean {
    return this.state !== "open";
  }

  getState(): string {
    return this.state;
  }
}
```

**Результат:**
- ✅ Exponential backoff для reconnect
- ✅ Circuit breaker pattern для защиты от бесконечных попыток
- ✅ Jitter для предотвращения thundering herd
- ✅ Улучшенная стабильность при плохой сети

---

## 📝 План реализации

### Фаза 1: Критичные функции (1-2 недели)

**Неделя 1:**
- [ ] Реализовать Reconnect Mechanism
  - [ ] Создать `reconnect-manager.ts`
  - [ ] Добавить `getMessagesSince()` в `api.ts`
  - [ ] Интегрировать в `use-chat.ts`
  - [ ] Обновить store для хранения timestamp
  - [ ] Тестирование reconnect логики

**Неделя 2:**
- [ ] Реализовать Rate Limiting Error Handling
  - [ ] Создать `RateLimitNotification.tsx`
  - [ ] Обновить `use-chat.ts` для обработки errors
  - [ ] Обновить store для rate limit state
  - [ ] Блокировать chat-input при rate limit
  - [ ] Тестирование rate limit UI

- [ ] Реализовать Session Ownership Validation
  - [ ] Добавить `signSession()` в `session-api.ts`
  - [ ] Интегрировать signed sessions в `use-chat.ts`
  - [ ] Обновить все API calls для использования signed sessions
  - [ ] Обработка ownership errors
  - [ ] Тестирование signed sessions

### Фаза 2: UX Improvements (1 неделя)

- [ ] Rate Limit Headers Tracking
  - [ ] Создать `rate-limit-tracker.ts`
  - [ ] Обновить все API calls
  - [ ] Создать `RateLimitIndicator.tsx`
  - [ ] Интегрировать в UI

- [ ] Improved Connection Status Indicators
  - [ ] Обновить `chat-header.tsx`
  - [ ] Добавить анимации
  - [ ] Тестирование UI states

- [ ] Optimistic Updates с Rollback
  - [ ] Обновить Message type
  - [ ] Добавить rollback logic
  - [ ] Обновить `ChatMessage.tsx`
  - [ ] Добавить Retry button

### Фаза 3: Nice to have (опционально)

- [ ] WebSocket Metrics
  - [ ] Создать `ws-metrics.ts`
  - [ ] Интегрировать в `use-chat.ts`
  - [ ] Создать metrics dashboard (опционально)

- [ ] Advanced Retry Logic
  - [ ] Создать `reconnect-strategy.ts`
  - [ ] Интегрировать exponential backoff
  - [ ] Добавить circuit breaker
  - [ ] Тестирование

---

## 🧪 Тестирование

### Unit Tests

```bash
# Test reconnect logic
npm test reconnect-manager.test.ts

# Test rate limit tracking
npm test rate-limit-tracker.test.ts

# Test session signing
npm test session-api.test.ts
```

### Integration Tests

```typescript
// Test reconnect with missed messages
describe("Reconnect Mechanism", () => {
  it("should recover missed messages after reconnect", async () => {
    // 1. Connect WebSocket
    // 2. Send message
    // 3. Disconnect
    // 4. Send message from another device (via API)
    // 5. Reconnect
    // 6. Verify missed message is recovered
  });
});

// Test rate limiting
describe("Rate Limiting", () => {
  it("should block input when rate limited", async () => {
    // 1. Send many messages quickly
    // 2. Receive rate_limit_exceeded error
    // 3. Verify input is blocked
    // 4. Wait for retry_after
    // 5. Verify input is unblocked
  });
});

// Test signed sessions
describe("Session Ownership", () => {
  it("should use signed session ID for authenticated users", async () => {
    // 1. Login
    // 2. Create session
    // 3. Verify signed session ID is obtained
    // 4. Send message with signed session
    // 5. Verify backend accepts signed session
  });
});
```

### Manual Testing Checklist

- [ ] Test reconnect after WiFi disconnect
- [ ] Test reconnect with missed messages
- [ ] Test rate limiting by sending many messages
- [ ] Test rate limit notification UI
- [ ] Test signed sessions after login
- [ ] Test session ownership errors
- [ ] Test optimistic updates with network failure
- [ ] Test retry button for failed messages
- [ ] Test connection status indicators
- [ ] Test rate limit indicator

---

## 📚 Документация для разработчиков

### Новые API endpoints (используемые frontend)

```
GET /api/chat/messages/since?session_id=xxx&since=2024-01-01T00:00:00Z
Response:
{
  "messages": [...],
  "session_id": "abc123",
  "message_count": 5,
  "since": "2024-01-01T00:00:00Z"
}

POST /api/sessions/sign
Request: { "session_id": "abc123" }
Response:
{
  "signed_session_id": "abc123.1699999999.uuid.signature",
  "session_id": "abc123",
  "expires_at": "2024-01-02T00:00:00Z"
}
```

### WebSocket Protocol Updates

```javascript
// Rate limit error
{
  "type": "error",
  "error": "rate_limit_exceeded",
  "message": "Rate limit exceeded (connection): 25 messages in 1m0s. Blocked for 30s. Retry after 30 seconds"
}

// Session ownership error
{
  "type": "error",
  "error": "unauthorized",
  "message": "Session ownership validation failed"
}
```

---

## 🎯 Итоги

### Что нужно реализовать:

**Приоритет 1 (Критично):**
1. ✅ Reconnect Mechanism с восстановлением сообщений
2. ✅ Rate Limiting Error Handling
3. ✅ Session Ownership Validation

**Приоритет 2 (UX):**
4. ✅ Rate Limit Headers Tracking
5. ✅ Improved Connection Status
6. ✅ Optimistic Updates с Rollback

**Приоритет 3 (Nice to have):**
7. ✅ WebSocket Metrics
8. ✅ Advanced Retry Logic

### Оценка трудозатрат:

- **Приоритет 1:** ~2 недели (1 developer)
- **Приоритет 2:** ~1 неделя (1 developer)
- **Приоритет 3:** ~3-5 дней (опционально)

**Итого:** ~3-4 недели для полной реализации

---

## 🔗 Связанные документы

- [BACKEND_IMPROVEMENTS.md](./BACKEND_IMPROVEMENTS.md) - Backend изменения, требующие frontend обновлений
- [CLAUDE.md](./CLAUDE.md) - Общая архитектура проекта
- [MONITORING.md](./MONITORING.md) - Мониторинг и метрики

---

**Автор:** Claude (Anthropic AI)
**Дата:** 12 ноября 2024
**Версия:** 1.0.0
