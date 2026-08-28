package apis

import (
	"net/http"
	"time"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

// Response DTOs use the JSON shape the React client's requests/conversations
// module expects (see UI/src/requests/conversations/types.ts) rather than
// the system package's own field tags.

type conversationSummaryResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type messageSourceResponse struct {
	MaterialID string `json:"materialId"`
	Name       string `json:"name"`
}

type conversationMessageResponse struct {
	ID        string                  `json:"id"`
	Role      system.MessageRole      `json:"role"`
	Content   string                  `json:"content"`
	Status    string                  `json:"status"`
	Sources   []messageSourceResponse `json:"sources,omitempty"`
	CreatedAt time.Time               `json:"createdAt"`
}

type conversationDetailResponse struct {
	ID       string                        `json:"id"`
	Title    string                        `json:"title"`
	Messages []conversationMessageResponse `json:"messages"`
}

func toConversationSummary(c *system.Conversation) conversationSummaryResponse {
	return conversationSummaryResponse{ID: c.ID, Title: c.Title, UpdatedAt: c.UpdatedAt}
}

func toMessageSources(sources []system.MessageSource) []messageSourceResponse {
	if len(sources) == 0 {
		return nil
	}
	out := make([]messageSourceResponse, len(sources))
	for i, s := range sources {
		out[i] = messageSourceResponse{MaterialID: s.MaterialID, Name: s.Name}
	}
	return out
}

func toConversationDetail(c *system.Conversation) conversationDetailResponse {
	messages := make([]conversationMessageResponse, len(c.Messages))
	for i, m := range c.Messages {
		messages[i] = conversationMessageResponse{
			ID:        m.ID,
			Role:      m.Role,
			Content:   m.Content,
			Status:    "complete",
			Sources:   toMessageSources(m.Sources),
			CreatedAt: m.CreatedAt,
		}
	}
	return conversationDetailResponse{ID: c.ID, Title: c.Title, Messages: messages}
}

// handleListConversations implements GET /api/v1/conversations.
func (a *App) handleListConversations(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	convs := a.Store.ListConversations(user.ID)

	out := make([]conversationSummaryResponse, len(convs))
	for i, c := range convs {
		out[i] = toConversationSummary(c)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleCreateConversation implements POST /api/v1/conversations.
func (a *App) handleCreateConversation(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	conv := a.Store.CreateConversation(user.ID)

	logs.FromContext(r.Context()).Log("info", "created conversation "+conv.ID)
	writeJSON(w, http.StatusCreated, toConversationSummary(conv))
}

// handleGetConversation implements GET /api/v1/conversations/{id}.
func (a *App) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	conv, err := a.Store.GetConversation(user.ID, r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toConversationDetail(conv))
}

// handleDeleteConversation implements DELETE /api/v1/conversations/{id}.
func (a *App) handleDeleteConversation(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	id := r.PathValue("id")
	if err := a.Store.DeleteConversation(user.ID, id); err != nil {
		writeSystemError(w, err)
		return
	}
	logs.FromContext(r.Context()).Log("info", "deleted conversation "+id)
	writeJSON(w, http.StatusOK, nil)
}
