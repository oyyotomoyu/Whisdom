package apis

import "net/http"

type healthResponse struct {
	Status         string `json:"status"`
	Database       string `json:"database"`
	Storage        string `json:"storage"`
	Model          string `json:"model"`
	VectorDatabase string `json:"vector_database"`
}

// handleHealth reports basic liveness. The in-memory store and local disk
// storage can't meaningfully fail independently of the process itself yet,
// so they report "ok" whenever the handler runs; model and vector_database
// report "not_configured" until those backends exist.
func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Status:         "ok",
		Database:       "ok",
		Storage:        "ok",
		Model:          "not_configured",
		VectorDatabase: "not_configured",
	})
}
