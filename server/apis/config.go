package apis

import (
	"net/http"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

type configResponse struct {
	TrainingMaterialPath string `json:"training_material_path"`
}

// handleGetConfig implements GET /api/v1/config.
func (a *App) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, configResponse{TrainingMaterialPath: a.Store.TrainingMaterialPath()})
}

type updateConfigRequest struct {
	TrainingMaterialPath string `json:"training_material_path"`
}

// handleUpdateConfig implements PATCH /api/v1/config.
func (a *App) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var req updateConfigRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	logger := logs.FromContext(r.Context())
	user, ok := currentUser(r.Context())
	if !ok || !hasAnyPermission(user, system.PermConfigManage, system.PermMaterialsPathEdit) {
		writeError(w, http.StatusForbidden, "permission denied")
		return
	}
	previous := a.Store.TrainingMaterialPath()

	updated, err := a.Store.SetTrainingMaterialPath(req.TrainingMaterialPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Log("info", "training material path changed by "+user.ID+": "+previous+" -> "+updated)
	writeJSON(w, http.StatusOK, configResponse{TrainingMaterialPath: updated})
}
