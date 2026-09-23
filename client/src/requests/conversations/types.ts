export interface ConversationSummary {
  id: string;
  title: string;
  updatedAt: string;
}

export interface MessageSource {
  materialId: string;
  name: string;
}

export type MessageRole = "user" | "assistant";
export type MessageStatus = "pending" | "streaming" | "complete" | "error";

export interface ConversationMessage {
  id: string;
  role: MessageRole;
  content: string;
  status: MessageStatus;
  sources?: MessageSource[];
  createdAt: string;
}

export interface Conversation {
  id: string;
  title: string;
  messages: ConversationMessage[];
}
