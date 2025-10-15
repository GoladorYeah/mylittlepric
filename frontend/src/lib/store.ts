// frontend/src/lib/store.ts
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
  
  addMessage: (message: ChatMessage) => void;
  setMessages: (messages: ChatMessage[]) => void;
  setSessionId: (id: string) => void;
  setLoading: (loading: boolean) => void;
  setCountry: (country: string) => void;
  setLanguage: (language: string) => void;
  setSearchInProgress: (inProgress: boolean) => void;
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

      addMessage: (message) =>
        set((state) => ({ messages: [...state.messages, message] })),

      setMessages: (messages) => set({ messages }),

      setSessionId: (id) => set({ sessionId: id }),

      setLoading: (loading) => set({ isLoading: loading }),

      setCountry: (country) => set({ country }),

      setLanguage: (language) => set({ language }),

      setSearchInProgress: (inProgress) => set({ searchInProgress: inProgress }),

      clearMessages: () => set({ messages: [], isLoading: false }),

      newSearch: () =>
        set({
          messages: [],
          searchInProgress: false,
          isLoading: false,
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