package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/whisdom/server/ai"
	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

const modelRequestTimeout = 30 * time.Second

type sendMessageRequest struct {
	Message string `json:"message"`
}

// chatResponse matches the wire shape UI/src/requests/conversations/index.ts
// expects: a snake_case top level (mirroring docs/server.md's example)
// whose sources reuse the camelCase MessageSource shape used elsewhere.
type chatResponse struct {
	MessageID string                  `json:"message_id"`
	Answer    string                  `json:"answer"`
	Sources   []messageSourceResponse `json:"sources"`
}

// handleSendMessage implements POST /api/v1/conversations/{id}/messages: the
// full chat flow described in docs/server.md — store the user message,
// retrieve relevant knowledge, call the model, store and return the answer.
func (a *App) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r.Context())
	logger := logs.FromContext(r.Context())
	conversationID := r.PathValue("id")

	var req sendMessageRequest
	if err := decodeJSON(r, &req); err != nil || req.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}

	conv, err := a.Store.GetConversation(user.ID, conversationID)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	history := make([]ai.HistoryTurn, len(conv.Messages))
	for i, m := range conv.Messages {
		history[i] = ai.HistoryTurn{Role: string(m.Role), Content: m.Content}
	}

	if _, err := a.Store.AppendMessage(user.ID, conversationID, system.MessageRoleUser, req.Message, nil); err != nil {
		writeSystemError(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), modelRequestTimeout)
	defer cancel()

	chunks, err := a.RAG.Retrieve(ctx, req.Message)
	if err != nil {
		logger.Log("warning", "RAG retrieval failed: "+err.Error())
		chunks = nil
	}

	result, err := a.Model.Generate(ctx, ai.GenerateRequest{Message: req.Message, Context: chunks, History: history})
	if err != nil {
		logger.Log("error", "model request failed: "+err.Error())
		writeError(w, http.StatusBadGateway, "the model could not generate a response")
		return
	}

	sources := make([]system.MessageSource, len(chunks))
	for i, c := range chunks {
		sources[i] = system.MessageSource{MaterialID: c.MaterialID, Name: c.Name}
	}

	assistantMsg, err := a.Store.AppendMessage(user.ID, conversationID, system.MessageRoleAssistant, result.Answer, sources)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logger.Log("info", "answered message in conversation "+conversationID)
	sourceResponses := toMessageSources(sources)
	if sourceResponses == nil {
		sourceResponses = []messageSourceResponse{}
	}
	writeJSON(w, http.StatusOK, chatResponse{
		MessageID: assistantMsg.ID,
		Answer:    assistantMsg.Content,
		Sources:   sourceResponses,
	})
}
