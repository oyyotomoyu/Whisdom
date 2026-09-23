package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/whisdom/server/config"
)

// Embedder turns text into fixed-length vectors for similarity search.
// Implementations must respect ctx cancellation and return one vector per
// input text, in order.
type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float64, error)
}

// hashEmbeddingDims is the fixed vector size HashEmbedder produces.
const hashEmbeddingDims = 256

// HashEmbedder is a dependency-free Embedder: it hashes each word into one
// of a fixed number of buckets (the "hashing trick") and L2-normalizes the
// result. It has no notion of meaning or synonymy, but unlike raw filename
// matching it does score chunks by real word overlap with the query, and it
// works fully offline. Use OllamaEmbedder for real semantic embeddings.
type HashEmbedder struct{}

// NewHashEmbedder returns an offline, dependency-free Embedder.
func NewHashEmbedder() *HashEmbedder { return &HashEmbedder{} }

// Embed implements Embedder.
func (e *HashEmbedder) Embed(ctx context.Context, texts []string) ([][]float64, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out := make([][]float64, len(texts))
	for i, t := range texts {
		out[i] = hashEmbed(t)
	}
	return out, nil
}

func hashEmbed(text string) []float64 {
	vec := make([]float64, hashEmbeddingDims)
	for _, word := range strings.Fields(strings.ToLower(text)) {
		word = strings.Trim(word, ".,!?;:()[]{}\"'")
		if len(word) < 2 {
			continue
		}
		h := fnv.New32a()
		_, _ = h.Write([]byte(word))
		vec[int(h.Sum32()%hashEmbeddingDims)]++
	}
	normalizeVector(vec)
	return vec
}

func normalizeVector(v []float64) {
	var sumSquares float64
	for _, x := range v {
		sumSquares += x * x
	}
	if sumSquares == 0 {
		return
	}
	norm := math.Sqrt(sumSquares)
	for i := range v {
		v[i] /= norm
	}
}

// cosineSimilarity returns a value in [-1, 1] (typically [0, 1] for
// non-negative embeddings); 0 for mismatched or zero vectors.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// OllamaEmbedder implements Embedder by calling a local Ollama runtime's
// embedding endpoint.
type OllamaEmbedder struct {
	client     *http.Client
	model      string
	runtimeURL string
}

// NewOllamaEmbedder returns an Embedder backed by a local Ollama model.
func NewOllamaEmbedder(cfg config.ModelConfig) *OllamaEmbedder {
	runtimeURL := strings.TrimRight(cfg.RuntimeURL, "/")
	if runtimeURL == "" {
		runtimeURL = "http://127.0.0.1:11434"
	}
	model := strings.TrimSpace(cfg.EmbeddingModel)
	if model == "" {
		model = "nomic-embed-text"
	}
	return &OllamaEmbedder{
		client:     &http.Client{Timeout: time.Minute},
		model:      model,
		runtimeURL: runtimeURL,
	}
}

type ollamaEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type ollamaEmbedResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
	Error      string      `json:"error,omitempty"`
}

// Embed implements Embedder.
func (e *OllamaEmbedder) Embed(ctx context.Context, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(ollamaEmbedRequest{Model: e.model, Input: texts})
	if err != nil {
		return nil, fmt.Errorf("encode ollama embed request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.runtimeURL+"/api/embed", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build ollama embed request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call ollama embed: %w", err)
	}
	defer resp.Body.Close()

	var result ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode ollama embed response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if result.Error != "" {
			return nil, fmt.Errorf("ollama embed returned %s: %s", resp.Status, result.Error)
		}
		return nil, fmt.Errorf("ollama embed returned %s", resp.Status)
	}
	if len(result.Embeddings) != len(texts) {
		return nil, fmt.Errorf("ollama embed returned %d vectors for %d inputs", len(result.Embeddings), len(texts))
	}
	return result.Embeddings, nil
}
