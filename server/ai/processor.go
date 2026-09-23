package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/whisdom/server/system"
)

// chunkRuneSize and chunkRuneOverlap bound each chunk handed to the
// embedder. Sized in runes (not bytes) so multi-byte text isn't split
// mid-character.
const (
	chunkRuneSize    = 800
	chunkRuneOverlap = 100
)

// ChunkEmbedProcessor implements system.MaterialProcessor: it extracts text
// from a material's stored file, splits it into overlapping chunks, and
// embeds each chunk.
//
// There is no PDF/DOCX/etc. text-extraction library wired in, so extraction
// only reads plain-text-shaped files (source code, markdown, csv, json,
// yaml, logs, ...). Binary formats and remote sources with no local file
// (wiki/git/url) fall back to a single chunk built from the material's name,
// which keeps them searchable by title even though their body isn't indexed.
type ChunkEmbedProcessor struct {
	Embedder Embedder
}

// NewMaterialProcessor returns a system.MaterialProcessor that extracts,
// chunks, and embeds a material's content using embedder.
func NewMaterialProcessor(embedder Embedder) *ChunkEmbedProcessor {
	return &ChunkEmbedProcessor{Embedder: embedder}
}

// Process implements system.MaterialProcessor.
func (p *ChunkEmbedProcessor) Process(ctx context.Context, m *system.Material) ([]system.ProcessedChunk, error) {
	text := extractText(m)
	pieces := chunkText(text)
	if len(pieces) == 0 {
		pieces = []string{fallbackText(m)}
	}

	embeddings, err := p.Embedder.Embed(ctx, pieces)
	if err != nil {
		return nil, fmt.Errorf("embed material %s: %w", m.ID, err)
	}
	if len(embeddings) != len(pieces) {
		return nil, fmt.Errorf("embedder returned %d vectors for %d chunks", len(embeddings), len(pieces))
	}

	out := make([]system.ProcessedChunk, len(pieces))
	for i, text := range pieces {
		out[i] = system.ProcessedChunk{Text: text, Embedding: embeddings[i]}
	}
	return out, nil
}

// extractText reads and returns a material's plain-text content, falling
// back to its name when there's no local file or the file isn't text.
func extractText(m *system.Material) string {
	if m.StoragePath == "" {
		return fallbackText(m)
	}
	data, err := os.ReadFile(m.StoragePath)
	if err != nil || !looksLikeText(data) {
		return fallbackText(m)
	}
	return string(data)
}

func fallbackText(m *system.Material) string {
	if m.Filename != "" {
		return m.Filename
	}
	return m.SourceLocation
}

func looksLikeText(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	const sniffLen = 8000
	sample := data
	if len(sample) > sniffLen {
		sample = sample[:sniffLen]
	}
	if !utf8.Valid(sample) {
		return false
	}
	for _, b := range sample {
		if b == 0 {
			return false
		}
	}
	return true
}

// chunkText splits text into overlapping, rune-bounded chunks.
func chunkText(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	runes := []rune(text)
	if len(runes) <= chunkRuneSize {
		return []string{string(runes)}
	}

	step := chunkRuneSize - chunkRuneOverlap
	var chunks []string
	for start := 0; start < len(runes); start += step {
		end := min(start+chunkRuneSize, len(runes))
		chunks = append(chunks, strings.TrimSpace(string(runes[start:end])))
		if end == len(runes) {
			break
		}
	}
	return chunks
}
