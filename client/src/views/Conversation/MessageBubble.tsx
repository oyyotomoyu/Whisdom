import { useState } from "react";
import { useTranslation } from "../../../node_modules/react-i18next";
import { useNavigate } from "react-router-dom";
import type { ConversationMessage } from "../../requests/conversations/types";
import { useAuth } from "../../hooks/useAuth";
import "./MessageBubble.css";

interface MessageBubbleProps {
  message: ConversationMessage;
  onRetry: () => void;
}

export function MessageBubble({ message, onRetry }: MessageBubbleProps) {
  const { t } = useTranslation("conversation");
  const { isAdmin } = useAuth();
  const navigate = useNavigate();
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(message.content);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const handleCorrect = () => {
    const params = new URLSearchParams({
      messageId: message.id,
      answer: message.content,
    });
    navigate(`/admin/corrections?${params.toString()}`);
  };

  const isUser = message.role === "user";

  return (
    <div className={`message-row ${isUser ? "user" : "assistant"}`}>
      <div className={`message-bubble ${message.status === "error" ? "error" : ""}`}>
        {message.status === "pending" ? (
          <span className="message-thinking">{t("message.thinking")}</span>
        ) : message.status === "error" ? (
          <span>{t("message.error")}</span>
        ) : (
          <p className="message-content">{message.content}</p>
        )}

        {message.sources && message.sources.length > 0 && (
          <div className="message-sources">
            <span className="message-sources-label">{t("message.sources")}</span>
            <ul>
              {message.sources.map((source) => (
                <li key={source.materialId}>{source.name}</li>
              ))}
            </ul>
          </div>
        )}

        {!isUser && message.status === "complete" && (
          <div className="message-actions">
            <button type="button" onClick={handleCopy}>
              {copied ? t("message.copy") + " ✓" : t("message.copy")}
            </button>
            {isAdmin && (
              <button type="button" onClick={handleCorrect}>
                {t("message.correct")}
              </button>
            )}
          </div>
        )}

        {message.status === "error" && (
          <div className="message-actions">
            <button type="button" onClick={onRetry}>
              {t("message.retry")}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
