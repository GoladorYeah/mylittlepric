import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { ChatMessage } from "@/types";
import { detectCountry, detectLanguage } from "./locale";

interface ChatStore {
  messages: ChatMessage[];
  sessionId: string;
  isLoading: boolean;
  country: string;
  language: string;
  searchInProgress: boolean;
  currentCategory: string;
  
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
    }),
    {
      name: "chat-storage",
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        country: state.country,
        language: state.language,
      }),
    }
  )
);