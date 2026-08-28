import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { NavLink, useNavigate, useParams } from "react-router-dom";
import { useConversationStore } from "../store/conversationStore";
import { useDebounce } from "../hooks/useDebounce";
import "./ConversationSidebar.css";

export function ConversationSidebar({ onNavigate }: { onNavigate?: () => void }) {
  const { t } = useTranslation("conversation");
  const navigate = useNavigate();
  const { conversationId: activeId } = useParams();

  const conversations = useConversationStore((s) => s.conversations);
  const isLoadingList = useConversationStore((s) => s.isLoadingList);
  const loadConversations = useConversationStore((s) => s.loadConversations);
  const removeConversation = useConversationStore((s) => s.removeConversation);

  const [query, setQuery] = useState("");
  const debouncedQuery = useDebounce(query, 200);

  useEffect(() => {
    loadConversations();
  }, [loadConversations]);

  const filtered = useMemo(() => {
    if (!debouncedQuery.trim()) return conversations;
    const q = debouncedQuery.trim().toLowerCase();
    return conversations.filter((c) => c.title.toLowerCase().includes(q));
  }, [conversations, debouncedQuery]);

  const handleNewChat = () => {
    navigate("/chat");
    onNavigate?.();
  };

  const handleDelete = async (event: React.MouseEvent, id: string) => {
    event.preventDefault();
    event.stopPropagation();
    if (window.confirm(t("deleteConfirm"))) {
      await removeConversation(id);
      if (activeId === id) navigate("/chat");
    }
  };

  return (
    <div className="conversation-sidebar">
      <button type="button" className="new-chat-button" onClick={handleNewChat}>
        + {t("newChat")}
      </button>

      <input
        type="search"
        className="sidebar-search"
        placeholder={t("searchPlaceholder")}
        value={query}
        onChange={(event) => setQuery(event.target.value)}
      />

      <div className="sidebar-section-label">{t("recentTitle")}</div>

      <nav className="conversation-list" aria-busy={isLoadingList}>
        {filtered.length === 0 && !isLoadingList && <p className="conversation-list-empty">{t("emptyRecent")}</p>}
        {filtered.map((conversation) => (
          <NavLink
            key={conversation.id}
            to={`/chat/${conversation.id}`}
            className={({ isActive }) => `conversation-item${isActive ? " active" : ""}`}
            onClick={onNavigate}
          >
            <span className="conversation-item-title">{conversation.title || t("untitled")}</span>
            <button
              type="button"
              className="conversation-item-delete"
              aria-label="Delete conversation"
              onClick={(event) => handleDelete(event, conversation.id)}
            >
              ×
            </button>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}
