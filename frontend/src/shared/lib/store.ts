import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { ChatMessage } from "@/shared/types";
import { detectCountry, detectLanguage, getCurrencyForCountry } from "./locale";

export interface SavedSearch {
  messages: ChatMessage[];
  sessionId: string;
  category: string;
  timestamp: number;
}

type WebSocketSender = (message: any) => void;

interface ChatStore {
  messages: ChatMessage[];
  sessionId: string;
  isLoading: boolean;
  country: string;
  language: string;
  currency: string;
  searchInProgress: boolean;
  currentCategory: string;
  isSidebarOpen: boolean;
  _hasInitialized: boolean; // Internal flag to track initialization
  savedSearch: SavedSearch | null; // Last search before "New Search" was clicked
  _wsSender: WebSocketSender | null; // Internal WebSocket sender for realtime sync

  addMessage: (message: ChatMessage) => void;
  setMessages: (messages: ChatMessage[]) => void;
  setSessionId: (id: string) => void;
  setLoading: (loading: boolean) => void;
  setCountry: (country: string) => void;
  setLanguage: (language: string) => void;
  setCurrency: (currency: string) => void;
  setSearchInProgress: (inProgress: boolean) => void;
  setCurrentCategory: (category: string) => void;
  clearMessages: () => void;
  clearAll: () => void;
  newSearch: () => void;
  initializeLocale: () => Promise<void>;
  loadSessionMessages: (sessionId: string) => Promise<void>;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  saveCurrentSearch: () => void;
  restoreSavedSearch: () => void;
  clearSavedSearch: () => void;
  syncPreferencesFromServer: () => Promise<void>;
  syncPreferencesToServer: () => Promise<void>;
  registerWebSocketSender: (sender: WebSocketSender | null) => void;
}

export const useChatStore = create<ChatStore>()(
  persist(
    (set, get) => ({
      messages: [],
      sessionId: "",
      isLoading: false,
      country: "",
      language: "",
      currency: "",
      searchInProgress: false,
      currentCategory: "",
      isSidebarOpen: true, // По умолчанию развернута
      _hasInitialized: false,
      savedSearch: null,
      _wsSender: null,

      addMessage: (message) =>
        set((state) => ({ messages: [...state.messages, message] })),

      setMessages: (messages) => set({ messages }),

      setSessionId: (id) => set({ sessionId: id }),

      setLoading: (loading) => set({ isLoading: loading }),

      setCountry: (country) => set({ country }),

      setLanguage: (language) => set({ language }),

      setCurrency: (currency) => set({ currency }),

      setSearchInProgress: (inProgress) => set({ searchInProgress: inProgress }),

      setCurrentCategory: (category) => set({ currentCategory: category }),

      clearMessages: () => set({ messages: [], isLoading: false, currentCategory: "" }),

      clearAll: () => {
        // Clear all chat data (for logout)
        localStorage.removeItem("chat_session_id");
        set({
          messages: [],
          sessionId: "",
          isLoading: false,
          searchInProgress: false,
          currentCategory: "",
          savedSearch: null,
          _hasInitialized: false,
        });
      },

      newSearch: () =>
        set({
          messages: [],
          searchInProgress: false,
          isLoading: false,
          currentCategory: "",
        }),

      initializeLocale: async () => {
        const state = get();
        // Only initialize if country is not already set (either from localStorage or detection)
        if (!state.country) {
          const country = await detectCountry();
          const currency = getCurrencyForCountry(country);
          set({ country, currency });
        } else if (!state.currency) {
          // If country exists but currency doesn't (migration case)
          const currency = getCurrencyForCountry(state.country);
          set({ currency });
        }
        if (!state.language) {
          set({ language: detectLanguage() });
        }
      },

      loadSessionMessages: async (sessionId: string) => {
        if (!sessionId) {
          console.warn("loadSessionMessages called with empty sessionId");
          return;
        }

        try {
          const { getSessionMessages } = await import("./api");
          const response = await getSessionMessages(sessionId);

          if (response.messages && response.messages.length > 0) {
            const chatMessages: ChatMessage[] = response.messages.map((msg, index) => ({
              id: `${sessionId}-${index}`,
              role: msg.role as "user" | "assistant",
              content: msg.content,
              timestamp: msg.timestamp ? new Date(msg.timestamp).getTime() : Date.now(),
              quick_replies: msg.quick_replies,
              products: msg.products,
              search_type: msg.search_type,
              isLocal: true, // Messages loaded from session are considered local (already sent)
            }));

            set({ messages: chatMessages });

            // Restore search state from server response
            let hasActiveSearch = false;
            let category = "";

            if (response.search_state) {
              category = response.search_state.category || "";
              hasActiveSearch = response.search_state.status === "completed";
            }

            // Also check if the last message has products - if so, consider it an active search
            // This ensures products are displayed when reopening a chat with search results
            if (!hasActiveSearch) {
              for (let i = chatMessages.length - 1; i >= 0; i--) {
                const msg = chatMessages[i];
                if (msg.products && msg.products.length > 0) {
                  hasActiveSearch = true;
                  // If we don't have a category from search_state, try to get it from message
                  if (!category && msg.search_type) {
                    category = msg.search_type;
                  }
                  break;
                }
              }
            }

            set({
              currentCategory: category,
              searchInProgress: hasActiveSearch,
            });

            console.log("✅ Session restored with", chatMessages.length, "messages",
                       hasActiveSearch ? "(with active search)" : "(no active search)");
          }
        } catch (error) {
          console.error("Failed to load session messages:", error);
        }
      },

      toggleSidebar: async () => {
        set((state) => ({ isSidebarOpen: !state.isSidebarOpen }));

        // Sync to server if user is authenticated
        const { useAuthStore } = await import("./auth-store");
        const isAuthenticated = useAuthStore.getState().isAuthenticated;

        if (isAuthenticated) {
          try {
            const { PreferencesAPI } = await import("./preferences-api");
            const newState = get().isSidebarOpen;
            await PreferencesAPI.updateUserPreferences({ sidebar_open: newState });
            console.log("✅ Synced sidebar state to server:", newState);

            // Realtime sync to other devices
            const state = get();
            if (state._wsSender) {
              const accessToken = useAuthStore.getState().accessToken;
              state._wsSender({
                type: "sync_preferences",
                session_id: state.sessionId,
                access_token: accessToken,
              });
            }
          } catch (error) {
            console.error("Failed to sync sidebar state:", error);
          }
        }
      },

      setSidebarOpen: async (open) => {
        set({ isSidebarOpen: open });

        // Sync to server if user is authenticated
        const { useAuthStore } = await import("./auth-store");
        const isAuthenticated = useAuthStore.getState().isAuthenticated;

        if (isAuthenticated) {
          try {
            const { PreferencesAPI } = await import("./preferences-api");
            await PreferencesAPI.updateUserPreferences({ sidebar_open: open });
            console.log("✅ Synced sidebar state to server:", open);

            // Realtime sync to other devices
            const state = get();
            if (state._wsSender) {
              const accessToken = useAuthStore.getState().accessToken;
              state._wsSender({
                type: "sync_preferences",
                session_id: state.sessionId,
                access_token: accessToken,
              });
            }
          } catch (error) {
            console.error("Failed to sync sidebar state:", error);
          }
        }
      },

      saveCurrentSearch: async () => {
        const state = get();
        // Only save if there are messages (otherwise nothing to save)
        if (state.messages.length > 0) {
          set({
            savedSearch: {
              messages: [...state.messages],
              sessionId: state.sessionId,
              category: state.currentCategory,
              timestamp: Date.now(),
            },
          });

          // Realtime sync to other devices
          if (state._wsSender) {
            const { useAuthStore } = await import("./auth-store");
            const accessToken = useAuthStore.getState().accessToken;
            if (accessToken) {
              state._wsSender({
                type: "sync_saved_search",
                session_id: state.sessionId,
                access_token: accessToken,
              });
            }
          }
        }
      },

      restoreSavedSearch: async () => {
        const state = get();
        if (state.savedSearch) {
          set({
            messages: [...state.savedSearch.messages],
            sessionId: state.savedSearch.sessionId,
            currentCategory: state.savedSearch.category,
            searchInProgress: false,
            isLoading: false,
          });
          // Save session ID to localStorage
          localStorage.setItem("chat_session_id", state.savedSearch.sessionId);

          // Realtime sync to other devices
          if (state._wsSender) {
            const { useAuthStore } = await import("./auth-store");
            const accessToken = useAuthStore.getState().accessToken;
            if (accessToken) {
              state._wsSender({
                type: "sync_session",
                session_id: state.savedSearch.sessionId,
                access_token: accessToken,
              });
            }
          }
        }
      },

      clearSavedSearch: async () => {
        const state = get();
        set({ savedSearch: null });

        // Realtime sync to other devices
        if (state._wsSender) {
          const { useAuthStore } = await import("./auth-store");
          const accessToken = useAuthStore.getState().accessToken;
          if (accessToken) {
            state._wsSender({
              type: "sync_saved_search",
              session_id: state.sessionId,
              access_token: accessToken,
            });
          }
        }
      },

      syncPreferencesFromServer: async () => {
        try {
          const { PreferencesAPI } = await import("./preferences-api");
          const response = await PreferencesAPI.getUserPreferences();

          if (response.has_preferences && response.preferences) {
            const prefs = response.preferences;
            const updates: Partial<ChatStore> = {};

            // Only update if server has value (don't override local with null)
            if (prefs.country) updates.country = prefs.country;
            if (prefs.currency) updates.currency = prefs.currency;
            if (prefs.language) updates.language = prefs.language;
            if (prefs.sidebar_open !== undefined) updates.isSidebarOpen = prefs.sidebar_open;

            if (Object.keys(updates).length > 0) {
              set(updates);
              console.log("✅ Synced preferences from server:", updates);
            }
          }
        } catch (error) {
          console.error("Failed to sync preferences from server:", error);
        }
      },

      syncPreferencesToServer: async () => {
        try {
          const state = get();
          const { PreferencesAPI } = await import("./preferences-api");

          const update = {
            country: state.country || undefined,
            currency: state.currency || undefined,
            language: state.language || undefined,
            sidebar_open: state.isSidebarOpen,
          };

          await PreferencesAPI.updateUserPreferences(update);
          console.log("✅ Synced preferences to server:", update);
        } catch (error) {
          console.error("Failed to sync preferences to server:", error);
        }
      },

      registerWebSocketSender: (sender) => {
        set({ _wsSender: sender });
      },
    }),
    {
      name: "chat-storage",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        country: state.country,
        language: state.language,
        currency: state.currency,
        isSidebarOpen: state.isSidebarOpen,
        messages: state.messages,
        sessionId: state.sessionId,
        currentCategory: state.currentCategory,
        searchInProgress: state.searchInProgress,
        savedSearch: state.savedSearch,
        // Exclude _wsSender from persistence
      }),
    }
  )
);