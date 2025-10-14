import { create } from "zustand";
import { ChatMessage, ProductCard } from "@/types";

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
}

export const useChatStore = create<ChatStore>((set) => ({
  messages: [],
  sessionId: "",
  isLoading: false,
  country: "CH",
  language: "de",
  searchInProgress: false,

  addMessage: (message) =>
    set((state) => ({ messages: [...state.messages, message] })),

  setMessages: (messages) => set({ messages }),

  setSessionId: (id) => set({ sessionId: id }),

  setLoading: (loading) => set({ isLoading: loading }),

  setCountry: (country) => set({ country }),

  setLanguage: (language) => set({ language }),

  setSearchInProgress: (inProgress) => set({ searchInProgress: inProgress }),

  clearMessages: () => set({ messages: [], sessionId: "" }),

  newSearch: () =>
    set({
      messages: [],
      searchInProgress: false,
    }),
}));