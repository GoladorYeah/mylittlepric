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
  addSearchToHistory: (query: string, category?: string, productsCount?: number) => void;
  loadSearchFromHistory: (historyItem: SearchHistoryItem) => void;
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
        if (!state.country || !state.language) {
          const country = await detectCountry();
          set({
            country,
            language: detectLanguage(),
          });
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

      loadSearchFromHistory: (historyItem) => {
        set({
          messages: [],
          searchInProgress: false,
          isLoading: false,
          currentCategory: historyItem.category || "",
          sessionId: historyItem.sessionId,
        });
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