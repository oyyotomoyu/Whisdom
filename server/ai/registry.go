package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// SwitchableModel wraps a ModelService so the target it dispatches to can be
// swapped at runtime (see apis' POST /api/v1/models/{id}/activate), without
// callers holding a stale reference to the old one.
type SwitchableModel struct {
	mu      sync.RWMutex
	current ModelService
}

// NewSwitchableModel returns a ModelService that starts out dispatching to
// initial and can be redirected later via Set.
func NewSwitchableModel(initial ModelService) *SwitchableModel {
	return &SwitchableModel{current: initial}
}

// Generate implements ModelService.
func (m *SwitchableModel) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	return m.snapshot().Generate(ctx, req)
}

// Summarize implements ModelService.
func (m *SwitchableModel) Summarize(ctx context.Context, content string) (string, error) {
	return m.snapshot().Summarize(ctx, content)
}

// Set redirects future calls to ms.
func (m *SwitchableModel) Set(ms ModelService) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = ms
}

func (m *SwitchableModel) snapshot() ModelService {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

// ListOllamaModels queries a local Ollama runtime for the models it
// currently has pulled, via its /api/tags endpoint.
func ListOllamaModels(ctx context.Context, runtimeURL string) ([]string, error) {
	runtimeURL = strings.TrimRight(runtimeURL, "/")
	if runtimeURL == "" {
		runtimeURL = "http://127.0.0.1:11434"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, runtimeURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("build ollama tags request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ollama tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ollama tags returned %s", resp.Status)
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode ollama tags response: %w", err)
	}

	names := make([]string, len(result.Models))
	for i, m := range result.Models {
		names[i] = m.Name
	}
	return names, nil
}
