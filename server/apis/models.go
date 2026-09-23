package apis

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/whisdom/server/ai"
	"github.com/whisdom/server/config"
	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

// ModelManager tracks which model target is currently active and can list
// and switch between the targets available for the server's configured
// provider. It owns the SwitchableModel that chat.go's Generate/Summarize
// calls actually run against, so activating a target takes effect
// immediately for every subsequent request.
type ModelManager struct {
	baseConfig config.ModelConfig
	switchable *ai.SwitchableModel

	mu       sync.RWMutex
	activeID string
}

// NewModelManager returns a ModelManager whose active target matches
// baseConfig (the target the server started with).
func NewModelManager(baseConfig config.ModelConfig, switchable *ai.SwitchableModel) *ModelManager {
	provider := normalizeProvider(baseConfig.Provider)
	return &ModelManager{
		baseConfig: baseConfig,
		switchable: switchable,
		activeID:   modelTargetID(provider, baseConfig.Name),
	}
}

func normalizeProvider(provider string) string {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" {
		provider = "stub"
	}
	return provider
}

func modelTargetID(provider, name string) string {
	return provider + ":" + name
}

// ActiveID returns the currently active target's ID.
func (mgr *ModelManager) ActiveID() string {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()
	return mgr.activeID
}

// Targets lists the model targets selectable for this server's configured
// provider: a fixed "stub" placeholder, plus, when the provider is ollama,
// every model currently pulled into the local Ollama runtime (discovered
// live via its /api/tags endpoint).
func (mgr *ModelManager) Targets(ctx context.Context) []system.ModelTarget {
	targets := []system.ModelTarget{{
		ID:          modelTargetID("stub", "stub"),
		Name:        "stub",
		Provider:    "stub",
		Description: "Deterministic placeholder; no model runtime required.",
	}}

	if normalizeProvider(mgr.baseConfig.Provider) != "ollama" {
		return targets
	}

	names, err := ai.ListOllamaModels(ctx, mgr.baseConfig.RuntimeURL)
	if err != nil || len(names) == 0 {
		targets = append(targets, system.ModelTarget{
			ID:          modelTargetID("ollama", mgr.baseConfig.Name),
			Name:        mgr.baseConfig.Name,
			Provider:    "ollama",
			Description: "Configured Ollama model (its /api/tags catalog was unavailable).",
		})
		return targets
	}
	for _, name := range names {
		targets = append(targets, system.ModelTarget{ID: modelTargetID("ollama", name), Name: name, Provider: "ollama"})
	}
	return targets
}

// Activate switches the live model to the target named by id, constructing
// a fresh ModelService for it and redirecting the shared SwitchableModel.
func (mgr *ModelManager) Activate(ctx context.Context, id string) (system.ModelTarget, error) {
	for _, target := range mgr.Targets(ctx) {
		if target.ID != id {
			continue
		}

		var svc ai.ModelService
		switch target.Provider {
		case "stub":
			svc = ai.NewStubModel()
		case "ollama":
			targetCfg := mgr.baseConfig
			targetCfg.Provider = "ollama"
			targetCfg.Name = target.Name
			svc = ai.NewOllamaModel(targetCfg)
		default:
			return system.ModelTarget{}, fmt.Errorf("unsupported provider %q", target.Provider)
		}

		mgr.switchable.Set(svc)
		mgr.mu.Lock()
		mgr.activeID = id
		mgr.mu.Unlock()
		return target, nil
	}
	return system.ModelTarget{}, system.ErrNotFound
}

type modelTargetResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Description string `json:"description,omitempty"`
	Active      bool   `json:"active"`
}

func toModelTargetResponse(t system.ModelTarget, activeID string) modelTargetResponse {
	return modelTargetResponse{ID: t.ID, Name: t.Name, Provider: t.Provider, Description: t.Description, Active: t.ID == activeID}
}

// handleListModels implements GET /api/v1/models.
func (a *App) handleListModels(w http.ResponseWriter, r *http.Request) {
	activeID := a.Models.ActiveID()
	targets := a.Models.Targets(r.Context())
	out := make([]modelTargetResponse, len(targets))
	for i, t := range targets {
		out[i] = toModelTargetResponse(t, activeID)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleCurrentModel implements GET /api/v1/models/current.
func (a *App) handleCurrentModel(w http.ResponseWriter, r *http.Request) {
	activeID := a.Models.ActiveID()
	for _, t := range a.Models.Targets(r.Context()) {
		if t.ID == activeID {
			writeJSON(w, http.StatusOK, toModelTargetResponse(t, activeID))
			return
		}
	}
	writeError(w, http.StatusNotFound, "active model target not found")
}

// handleActivateModel implements POST /api/v1/models/{id}/activate.
func (a *App) handleActivateModel(w http.ResponseWriter, r *http.Request) {
	target, err := a.Models.Activate(r.Context(), r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	logs.FromContext(r.Context()).Log("info", "activated model target "+target.ID)
	writeJSON(w, http.StatusOK, toModelTargetResponse(target, target.ID))
}
