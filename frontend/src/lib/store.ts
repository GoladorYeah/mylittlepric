import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { ChatMessage } from "@/types";
import { detectCountry, detectLanguage } from "./locale";

export interface SearchHistoryItem {
  id: string;
  query: string;
  timestamp: number;
  category?: string;
  productsCount?: number;
  sessionId: string;
}

interface ChatStore {
  messages: ChatMessage[];
  sessionId: string;
  isLoading: boolean;
  country: string;
  language: string;
  searchInProgress: boolean;
  currentCategory: string;
  searchHistory: SearchHistoryItem[];
  isSidebarOpen: boolean;
  _hasInitialized: boolean; // Internal flag to track initialization

  addMessage: (message: ChatMessage) => void;
  setMessages: (messages: ChatMessage[]) => void;
  setSessionId: (id: string) => void;
  setLoading: (loading: boolean) => void;
  setCountry: (country: string) => void;
  setLanguage: (language: string) => void;
  setSearchInProgress: (inProgress: boolean) => void;
  setCurrentCategory: (category: string) => void;
  clearMessages: () => void;
  newSearch: () => void;
  initializeLocale: () => Promise<void>;
  loadSessionMessages: (sessionId: string) => Promise<void>;
  addSearchToHistory: (query: string, category?: string, productsCount?: number) => void;
  loadSearchFromHistory: (historyItem: SearchHistoryItem) => Promise<void>;
  clearSearchHistory: () => void;
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
}

export const useChatStore = create<ChatStore>()(
  persist(
    (set, get) => ({
      messages: [],
      sessionId: "",
      isLoading: false,
      country: "",
      language: "",
      searchInProgress: false,
      currentCategory: "",
      searchHistory: [],
      isSidebarOpen: false,
      _hasInitialized: false,

      addMessage: (message) =>
        set((state) => ({ messages: [...state.messages, message] })),

      setMessages: (messages) => set({ messages }),

      setSessionId: (id) => set({ sessionId: id }),

      setLoading: (loading) => set({ isLoading: loading }),

      setCountry: (country) => set({ country }),

      setLanguage: (language) => set({ language }),

      setSearchInProgress: (inProgress) => set({ searchInProgress: inProgress }),

      setCurrentCategory: (category) => set({ currentCategory: category }),

      clearMessages: () => set({ messages: [], isLoading: false, currentCategory: "" }),

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
          set({ country });
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

      addSearchToHistory: (query, category, productsCount) => {
        const state = get();
        const newHistoryItem: SearchHistoryItem = {
          id: `${Date.now()}-${Math.random()}`,
          query,
          timestamp: Date.now(),
          category,
          productsCount,
          sessionId: state.sessionId,
        };

        set((state) => ({
          searchHistory: [newHistoryItem, ...state.searchHistory].slice(0, 50), // Keep last 50 searches
        }));
      },

      loadSearchFromHistory: async (historyItem) => {
        // Clear current messages first
        set({
          messages: [],
          searchInProgress: false,
          isLoading: false,
          currentCategory: historyItem.category || "",
          sessionId: historyItem.sessionId,
        });

        // Save to localStorage
        localStorage.setItem("chat_session_id", historyItem.sessionId);

        // Load messages for this session
        const { loadSessionMessages } = get();
        await loadSessionMessages(historyItem.sessionId);
      },

      clearSearchHistory: () => set({ searchHistory: [] }),

      toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),

      setSidebarOpen: (open) => set({ isSidebarOpen: open }),
    }),
    {
      name: "chat-storage",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        country: state.country,
        language: state.language,
        searchHistory: state.searchHistory,
        isSidebarOpen: state.isSidebarOpen,
      }),
    }
  )
);