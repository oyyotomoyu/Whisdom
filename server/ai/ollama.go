package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/whisdom/server/config"
)

const libraryOnlyRefusal = "The library does not contain enough relevant information to answer this question."

// OllamaModel implements ModelService by calling a local Ollama runtime.
type OllamaModel struct {
	client      *http.Client
	model       string
	runtimeURL  string
	temperature float64
}

// NewOllamaModel returns a ModelService backed by a local Ollama model.
func NewOllamaModel(cfg config.ModelConfig) *OllamaModel {
	runtimeURL := strings.TrimRight(cfg.RuntimeURL, "/")
	if runtimeURL == "" {
		runtimeURL = "http://127.0.0.1:11434"
	}
	model := strings.TrimSpace(cfg.Name)
	if model == "" {
		model = "gemma4"
	}
	return &OllamaModel{
		client:      &http.Client{Timeout: 2 * time.Minute},
		model:       model,
		runtimeURL:  runtimeURL,
		temperature: cfg.Temperature,
	}
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  ollamaOptions   `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float64 `json:"temperature,omitempty"`
}

type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
	Error   string        `json:"error,omitempty"`
}

// Generate implements ModelService.
func (m *OllamaModel) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	if err := ctx.Err(); err != nil {
		return GenerateResponse{}, fmt.Errorf("model request canceled: %w", err)
	}
	if len(req.Context) == 0 {
		return GenerateResponse{Answer: libraryOnlyRefusal}, nil
	}

	body, err := json.Marshal(ollamaChatRequest{
		Model:    m.model,
		Messages: m.messages(req),
		Stream:   false,
		Options:  ollamaOptions{Temperature: m.temperature},
	})
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("encode ollama request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.runtimeURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("build ollama request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(httpReq)
	if err != nil {
		return GenerateResponse{}, fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GenerateResponse{}, fmt.Errorf("decode ollama response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if result.Error != "" {
			return GenerateResponse{}, fmt.Errorf("ollama returned %s: %s", resp.Status, result.Error)
		}
		return GenerateResponse{}, fmt.Errorf("ollama returned %s", resp.Status)
	}
	if result.Error != "" {
		return GenerateResponse{}, fmt.Errorf("ollama error: %s", result.Error)
	}

	answer := strings.TrimSpace(result.Message.Content)
	if answer == "" {
		return GenerateResponse{}, fmt.Errorf("ollama returned an empty answer")
	}
	return GenerateResponse{Answer: answer}, nil
}

// Summarize implements ModelService.
func (m *OllamaModel) Summarize(ctx context.Context, content string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("summarize request canceled: %w", err)
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return "", nil
	}

	body, err := json.Marshal(ollamaChatRequest{
		Model: m.model,
		Messages: []ollamaMessage{
			{Role: "system", Content: "Summarize the user's message as a short title of 6 words or fewer. Reply with only the title, no punctuation at the end, no quotes."},
			{Role: "user", Content: content},
		},
		Stream:  false,
		Options: ollamaOptions{Temperature: m.temperature},
	})
	if err != nil {
		return "", fmt.Errorf("encode ollama request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, m.runtimeURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build ollama request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call ollama: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode ollama response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if result.Error != "" {
			return "", fmt.Errorf("ollama returned %s: %s", resp.Status, result.Error)
		}
		return "", fmt.Errorf("ollama returned %s", resp.Status)
	}
	if result.Error != "" {
		return "", fmt.Errorf("ollama error: %s", result.Error)
	}

	return strings.Trim(strings.TrimSpace(result.Message.Content), "\""), nil
}

func (m *OllamaModel) messages(req GenerateRequest) []ollamaMessage {
	messages := []ollamaMessage{
		{Role: "system", Content: libraryOnlySystemPrompt(req.Context)},
	}
	for _, turn := range req.History {
		role := turn.Role
		if role != "user" && role != "assistant" {
			continue
		}
		content := strings.TrimSpace(turn.Content)
		if content == "" {
			continue
		}
		messages = append(messages, ollamaMessage{Role: role, Content: content})
	}
	messages = append(messages, ollamaMessage{Role: "user", Content: strings.TrimSpace(req.Message)})
	return messages
}

func libraryOnlySystemPrompt(chunks []ContextChunk) string {
	var b strings.Builder
	b.WriteString("You are Whisdom's local language model. Answer only using the library context below. ")
	b.WriteString("Do not use general knowledge. If the library context does not contain enough relevant information, answer exactly: ")
	b.WriteString(libraryOnlyRefusal)
	b.WriteString("\n\nLibrary context:\n")
	for i, chunk := range chunks {
		b.WriteString(fmt.Sprintf("[%d] %s\n", i+1, chunk.Name))
		text := strings.TrimSpace(chunk.Text)
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n")
		}
	}
	return b.String()
}
