import { httpClient, API_BASE_URL } from "../core";
import { useAuthStore } from "../../store/authStore";
import type { Conversation, ConversationSummary, MessageSource } from "./types";

export async function listConversations(): Promise<ConversationSummary[]> {
  const { data } = await httpClient.get<ConversationSummary[]>("/conversations");
  return data;
}

export async function createConversation(): Promise<ConversationSummary> {
  const { data } = await httpClient.post<ConversationSummary>("/conversations");
  return data;
}

export async function getConversation(id: string): Promise<Conversation> {
  const { data } = await httpClient.get<Conversation>(`/conversations/${id}`);
  return data;
}

export async function deleteConversation(id: string): Promise<void> {
  await httpClient.delete(`/conversations/${id}`);
}

export interface SendMessageResult {
  messageId: string;
  answer: string;
  sources: MessageSource[];
}

export async function sendMessage(conversationId: string, message: string): Promise<SendMessageResult> {
  const { data } = await httpClient.post<{ message_id: string; answer: string; sources: MessageSource[] }>(
    `/conversations/${conversationId}/messages`,
    { message }
  );
  return { messageId: data.message_id, answer: data.answer, sources: data.sources };
}

export interface StreamMessageHandlers {
  onToken: (token: string) => void;
  onSources?: (sources: MessageSource[]) => void;
  onDone: () => void;
  onError: (error: Error) => void;
}

/**
 * Streams an assistant reply over Server-Sent Events. Falls back to a single
 * `sendMessage` call if the backend does not support streaming for this deployment.
 */
export async function streamMessage(
  conversationId: string,
  message: string,
  handlers: StreamMessageHandlers,
  signal?: AbortSignal
): Promise<void> {
  const token = useAuthStore.getState().accessToken;

  try {
    const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/messages`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "text/event-stream",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: JSON.stringify({ message, stream: true }),
      signal,
    });

    if (!response.ok || !response.body) {
      throw new Error(`Request failed with status ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";

    for (;;) {
      const { value, done } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const events = buffer.split("\n\n");
      buffer = events.pop() ?? "";

      for (const event of events) {
        const line = event.trim();
        if (!line.startsWith("data:")) continue;
        const payload = line.slice(5).trim();
        if (payload === "[DONE]") {
          handlers.onDone();
          return;
        }
        const parsed = JSON.parse(payload) as { token?: string; sources?: MessageSource[] };
        if (parsed.token) handlers.onToken(parsed.token);
        if (parsed.sources) handlers.onSources?.(parsed.sources);
      }
    }

    handlers.onDone();
  } catch (error) {
    handlers.onError(error instanceof Error ? error : new Error("Streaming failed"));
  }
}

export type * from "./types";
