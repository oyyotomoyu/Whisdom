package apis

import (
	"net/http"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

// handleListRoles implements GET /api/v1/roles.
func (a *App) handleListRoles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.Store.ListRoles())
}

type createRoleRequest struct {
	Name        string              `json:"name"`
	Permissions []system.Permission `json:"permissions"`
}

// handleCreateRole implements POST /api/v1/roles.
func (a *App) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	var req createRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	role, err := a.Store.CreateRole(req.Name, req.Permissions)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "created role "+role.ID)
	writeJSON(w, http.StatusCreated, role)
}

type updateRoleRequest struct {
	Name        *string              `json:"name"`
	Permissions *[]system.Permission `json:"permissions"`
}

// handleUpdateRole implements PATCH /api/v1/roles/{id}.
func (a *App) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	var req updateRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	role, err := a.Store.UpdateRole(r.PathValue("id"), req.Name, req.Permissions)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "updated role "+role.ID)
	writeJSON(w, http.StatusOK, role)
}

// handleListPermissions implements GET /api/v1/permissions.
func (a *App) handleListPermissions(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, system.AllPermissions)
}
