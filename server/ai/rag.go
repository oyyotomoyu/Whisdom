package ai

import (
	"context"
	"strings"

	"github.com/whisdom/server/system"
)

// RAGService retrieves the company knowledge most relevant to a query.
type RAGService interface {
	Retrieve(ctx context.Context, query string) ([]ContextChunk, error)
}

// MaterialFilenameRAG is a placeholder RAGService: it matches query words
// against ready materials' filenames. There is no text extraction, chunking,
// or embedding pipeline yet (see docs/server.md's "Material And Training
// Path Config" section), so this cannot do real semantic retrieval — it
// exists so the chat API's source-citation contract can be exercised end to
// end before the real vector-search RAG backend lands.
type MaterialFilenameRAG struct {
	store *system.Store
}

// NewMaterialFilenameRAG returns a RAGService placeholder backed by store.
func NewMaterialFilenameRAG(store *system.Store) *MaterialFilenameRAG {
	return &MaterialFilenameRAG{store: store}
}

// Retrieve implements RAGService.
func (r *MaterialFilenameRAG) Retrieve(ctx context.Context, query string) ([]ContextChunk, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	keywords := strings.Fields(strings.ToLower(query))
	if len(keywords) == 0 {
		return nil, nil
	}

	const maxResults = 3
	var chunks []ContextChunk
	for _, m := range r.store.ListMaterials() {
		if m.Status != system.MaterialReady || !m.RAGAvailable {
			continue
		}
		name := strings.ToLower(m.Filename)
		for _, kw := range keywords {
			if len(kw) >= 3 && strings.Contains(name, kw) {
				chunks = append(chunks, ContextChunk{MaterialID: m.ID, Name: m.Filename})
				break
			}
		}
		if len(chunks) >= maxResults {
			break
		}
	}
	return chunks, nil
}
