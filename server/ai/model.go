// Package ai defines the model and retrieval interfaces the chat API is
// built against, plus lightweight placeholder implementations. Swapping in a
// real local model (the default target is Gemma 4, per docs/architecture.md)
// or a real vector-search RAG backend means implementing these same
// interfaces — callers never need to change.
package ai

import (
	"context"
	"fmt"
	"strings"
)

// ContextChunk is one piece of retrieved library knowledge handed to the
// model as grounding, and cited back to the user as a source.
type ContextChunk struct {
	MaterialID string
	Name       string
	Text       string
}

// HistoryTurn is one prior message in the conversation, oldest first.
type HistoryTurn struct {
	Role    string
	Content string
}

// GenerateRequest is one chat turn to answer.
type GenerateRequest struct {
	Message string
	Context []ContextChunk
	History []HistoryTurn
}

// GenerateResponse is the model's answer.
type GenerateResponse struct {
	Answer string
}

// ModelService generates chat answers, optionally grounded in retrieved
// context. Implementations must respect ctx cancellation/timeout and return
// a clear error rather than a partial or malformed answer.
type ModelService interface {
	Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
	// Summarize condenses content into a short, single-line summary (used,
	// e.g., to title a conversation from its first message).
	Summarize(ctx context.Context, content string) (string, error)
}

// StubModel is a placeholder ModelService used until a real local model
// runtime is wired in. It never calls out to a network or GPU; it composes a
// deterministic answer so the rest of the chat flow (storage, sources,
// corrections) can be built and tested end to end today.
type StubModel struct{}

// NewStubModel returns a ModelService placeholder.
func NewStubModel() *StubModel { return &StubModel{} }

// Generate implements ModelService.
func (m *StubModel) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	if err := ctx.Err(); err != nil {
		return GenerateResponse{}, fmt.Errorf("model request canceled: %w", err)
	}

	if len(req.Context) == 0 {
		return GenerateResponse{
			Answer: libraryOnlyRefusal,
		}, nil
	}

	var b strings.Builder
	b.WriteString("Based on the retrieved library knowledge:\n\n")
	for _, c := range req.Context {
		b.WriteString("- " + c.Name + "\n")
	}
	b.WriteString("\n(No local model is configured yet; this placeholder answer is assembled only from retrieved library sources.)")
	return GenerateResponse{Answer: b.String()}, nil
}

// summarizeMaxLen bounds a StubModel/OllamaModel summary; conversation
// titles have no business being long.
const summarizeMaxLen = 60

// Summarize implements ModelService. Without a local model configured, this
// is a first-sentence-or-truncation heuristic rather than a real summary.
func (m *StubModel) Summarize(ctx context.Context, content string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("summarize request canceled: %w", err)
	}
	return truncateSummary(content, summarizeMaxLen), nil
}

func truncateSummary(content string, maxLen int) string {
	content = strings.TrimSpace(content)
	if end := strings.IndexAny(content, ".!?\n"); end > 0 && end < maxLen {
		return content[:end]
	}
	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}
	return string(runes[:maxLen]) + "…"
}
