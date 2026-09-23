import { useEffect, useRef } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { useParams } from "react-router-dom";
import { useConversationStore } from "../../store/conversationStore";
import { MessageBubble } from "./MessageBubble";
import { Composer } from "./Composer";
import { EmptyState } from "../../components/EmptyState";
import "./Conversation.css";

export function Conversation() {
  const { t } = useTranslation("conversation");
  const { conversationId } = useParams();

  const messagesByConversation = useConversationStore((s) => s.messagesByConversation);
  const selectConversation = useConversationStore((s) => s.selectConversation);
  const startNewConversation = useConversationStore((s) => s.startNewConversation);
  const send = useConversationStore((s) => s.send);
  const retryLastMessage = useConversationStore((s) => s.retryLastMessage);
  const isSending = useConversationStore((s) => s.isSending);
  const sendError = useConversationStore((s) => s.sendError);

  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (conversationId) {
      selectConversation(conversationId);
    } else {
      startNewConversation();
    }
  }, [conversationId, selectConversation, startNewConversation]);

  const messages = conversationId ? (messagesByConversation[conversationId] ?? []) : [];

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages.length]);

  return (
    <div className="conversation-view">
      <div className="conversation-thread" ref={scrollRef}>
        {messages.length === 0 ? (
          <EmptyState title={t("emptyState.title")} subtitle={t("emptyState.subtitle")} />
        ) : (
          <div className="conversation-messages">
            {messages.map((message) => (
              <MessageBubble key={message.id} message={message} onRetry={retryLastMessage} />
            ))}
          </div>
        )}
      </div>
      <Composer onSend={send} isSending={isSending} error={sendError ? t(sendError) : null} />
    </div>
  );
}
