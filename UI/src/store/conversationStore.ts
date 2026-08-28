import { create } from "zustand";
import {
  createConversation,
  deleteConversation as deleteConversationRequest,
  getConversation,
  listConversations,
  sendMessage as sendMessageRequest,
} from "../requests/conversations";
import type { ConversationMessage, ConversationSummary } from "../requests/conversations/types";
import { ApiError } from "../requests/core";

interface ConversationState {
  conversations: ConversationSummary[];
  messagesByConversation: Record<string, ConversationMessage[]>;
  activeConversationId: string | null;
  isLoadingList: boolean;
  isLoadingMessages: boolean;
  isSending: boolean;
  listError: string | null;
  sendError: string | null;

  loadConversations: () => Promise<void>;
  selectConversation: (id: string) => Promise<void>;
  startNewConversation: () => void;
  removeConversation: (id: string) => Promise<void>;
  send: (text: string) => Promise<void>;
  retryLastMessage: () => Promise<void>;
}

function tempId(prefix: string) {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`;
}

export const useConversationStore = create<ConversationState>((set, get) => ({
  conversations: [],
  messagesByConversation: {},
  activeConversationId: null,
  isLoadingList: false,
  isLoadingMessages: false,
  isSending: false,
  listError: null,
  sendError: null,

  loadConversations: async () => {
    set({ isLoadingList: true, listError: null });
    try {
      const conversations = await listConversations();
      set({ conversations, isLoadingList: false });
    } catch (error) {
      set({ isLoadingList: false, listError: error instanceof ApiError ? error.message : "errors.generic" });
    }
  },

  selectConversation: async (id) => {
    set({ activeConversationId: id, isLoadingMessages: true });
    try {
      const conversation = await getConversation(id);
      set((state) => ({
        isLoadingMessages: false,
        messagesByConversation: { ...state.messagesByConversation, [id]: conversation.messages },
      }));
    } catch {
      set({ isLoadingMessages: false });
    }
  },

  startNewConversation: () => {
    set({ activeConversationId: null });
  },

  removeConversation: async (id) => {
    const previous = get().conversations;
    set((state) => ({ conversations: state.conversations.filter((c) => c.id !== id) }));
    try {
      await deleteConversationRequest(id);
      if (get().activeConversationId === id) {
        set({ activeConversationId: null });
      }
    } catch {
      set({ conversations: previous });
    }
  },

  send: async (text) => {
    const trimmed = text.trim();
    if (!trimmed) return;

    set({ isSending: true, sendError: null });

    let conversationId = get().activeConversationId;
    try {
      if (!conversationId) {
        const created = await createConversation();
        conversationId = created.id;
        set((state) => ({
          activeConversationId: created.id,
          conversations: [created, ...state.conversations],
        }));
      }

      const userMessage: ConversationMessage = {
        id: tempId("msg"),
        role: "user",
        content: trimmed,
        status: "complete",
        createdAt: new Date().toISOString(),
      };
      const pendingAssistant: ConversationMessage = {
        id: tempId("msg"),
        role: "assistant",
        content: "",
        status: "pending",
        createdAt: new Date().toISOString(),
      };

      const idForState = conversationId;
      set((state) => ({
        messagesByConversation: {
          ...state.messagesByConversation,
          [idForState]: [...(state.messagesByConversation[idForState] ?? []), userMessage, pendingAssistant],
        },
      }));

      const result = await sendMessageRequest(conversationId, trimmed);

      set((state) => ({
        isSending: false,
        messagesByConversation: {
          ...state.messagesByConversation,
          [idForState]: (state.messagesByConversation[idForState] ?? []).map((m) =>
            m.id === pendingAssistant.id
              ? { ...m, id: result.messageId, content: result.answer, sources: result.sources, status: "complete" }
              : m
          ),
        },
      }));
    } catch (error) {
      const message = error instanceof ApiError ? error.message : "conversation.message.error";
      set((state) => {
        const idForState = conversationId;
        if (!idForState) return { isSending: false, sendError: message };
        return {
          isSending: false,
          sendError: message,
          messagesByConversation: {
            ...state.messagesByConversation,
            [idForState]: (state.messagesByConversation[idForState] ?? []).map((m) =>
              m.status === "pending" ? { ...m, status: "error" } : m
            ),
          },
        };
      });
    }
  },

  retryLastMessage: async () => {
    const id = get().activeConversationId;
    if (!id) return;
    const messages = get().messagesByConversation[id] ?? [];
    const lastUser = [...messages].reverse().find((m) => m.role === "user");
    if (!lastUser) return;
    set((state) => ({
      messagesByConversation: {
        ...state.messagesByConversation,
        [id]: (state.messagesByConversation[id] ?? []).filter((m) => m.status !== "error"),
      },
    }));
    await get().send(lastUser.content);
  },
}));
