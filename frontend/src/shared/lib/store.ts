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
            }));

            set({ messages: chatMessages });

            if (response.search_state) {
              set({
                currentCategory: response.search_state.category || "",
                searchInProgress: response.search_state.status === "completed",
              });
            }
          }
        } catch (error) {
          console.error("Failed to load session messages:", error);
        }
      },

      toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),

      setSidebarOpen: (open) => set({ isSidebarOpen: open }),

      saveCurrentSearch: () => {
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
        }
      },

      restoreSavedSearch: () => {
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
        }
      },

      clearSavedSearch: () => set({ savedSearch: null }),
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
      }),
    }
  )
);