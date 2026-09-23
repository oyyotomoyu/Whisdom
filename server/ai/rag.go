package ai

import (
	"context"
	"sort"
	"strings"

	"github.com/whisdom/server/system"
)

// RAGService retrieves the library knowledge most relevant to a query.
type RAGService interface {
	Retrieve(ctx context.Context, query string) ([]ContextChunk, error)
}

// ragMaxResults and ragMinScore bound what VectorRAG hands to the model: at
// most a handful of chunks, and only ones that actually resemble the query,
// so unrelated material is never sent as context (docs/server.md's "Avoid
// sending unrelated material to the model").
const (
	ragMaxResults = 4
	ragMinScore   = 0.05
)

// VectorRAG retrieves the chunks whose embeddings are most similar to the
// query's embedding, searching every ready material's chunk store. It
// replaces the earlier filename-substring placeholder now that materials are
// actually chunked and embedded (see processor.go).
type VectorRAG struct {
	store    *system.Store
	embedder Embedder
}

// NewVectorRAG returns a RAGService backed by store's embedded material
// chunks.
func NewVectorRAG(store *system.Store, embedder Embedder) *VectorRAG {
	return &VectorRAG{store: store, embedder: embedder}
}

type scoredChunk struct {
	chunk        system.MaterialChunk
	materialName string
	score        float64
}

// Retrieve implements RAGService.
func (r *VectorRAG) Retrieve(ctx context.Context, query string) ([]ContextChunk, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	vectors, err := r.embedder.Embed(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	queryVector := vectors[0]

	var candidates []scoredChunk
	for _, m := range r.store.ListMaterials() {
		if m.Status != system.MaterialReady || !m.RAGAvailable {
			continue
		}
		for _, chunk := range r.store.MaterialChunks(m.ID) {
			score := cosineSimilarity(queryVector, chunk.Embedding)
			if score < ragMinScore {
				continue
			}
			candidates = append(candidates, scoredChunk{chunk: chunk, materialName: m.Filename, score: score})
		}
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].score > candidates[j].score })
	if len(candidates) > ragMaxResults {
		candidates = candidates[:ragMaxResults]
	}

	out := make([]ContextChunk, len(candidates))
	for i, c := range candidates {
		out[i] = ContextChunk{MaterialID: c.chunk.MaterialID, Name: c.materialName, Text: c.chunk.Text}
	}
	return out, nil
}
