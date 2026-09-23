package apis

import (
	"github.com/whisdom/server/ai"
	"github.com/whisdom/server/config"
	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

// App holds every dependency handlers need. It is constructed once in
// main.go and passed to NewRouter.
type App struct {
	Store  *system.Store
	Logs   *logs.Service
	Config *config.Config
	Model  ai.ModelService
	RAG    ai.RAGService
	// Models manages the active model target for the Models API
	// (GET /api/v1/models, .../current, POST .../{id}/activate). It owns the
	// same SwitchableModel that Model is set to.
	Models *ModelManager
}
