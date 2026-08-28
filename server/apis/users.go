package apis

import (
	"net/http"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

type userAccountResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

func (a *App) toUserAccount(u *system.User) (userAccountResponse, error) {
	role, err := a.Store.GetRole(u.RoleID)
	if err != nil {
		return userAccountResponse{}, err
	}
	return userAccountResponse{ID: u.ID, Name: u.Name, Email: u.Email, Role: role.Name, Active: u.Active}, nil
}

// handleListUsers implements GET /api/v1/users.
func (a *App) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users := a.Store.ListUsers()
	out := make([]userAccountResponse, 0, len(users))
	for _, u := range users {
		account, err := a.toUserAccount(u)
		if err != nil {
			continue
		}
		out = append(out, account)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleGetUser implements GET /api/v1/users/{id}.
func (a *App) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := a.Store.GetUser(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	account, err := a.toUserAccount(user)
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, account)
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

// handleCreateUser implements POST /api/v1/users.
func (a *App) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Email == "" || req.Role == "" || len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "name, email, role, and a password of at least 8 characters are required")
		return
	}

	role, ok := a.Store.RoleByName(req.Role)
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown role")
		return
	}

	hash, err := system.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	user, err := a.Store.CreateUser(req.Name, req.Email, hash, role.ID)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "created user "+user.ID)

	account, err := a.toUserAccount(user)
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, account)
}

type updateUserRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Role   *string `json:"role"`
	Active *bool   `json:"active"`
}

// handleUpdateUser implements PATCH /api/v1/users/{id}.
func (a *App) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	patch := system.UserPatch{Name: req.Name, Email: req.Email, Active: req.Active}
	if req.Role != nil {
		role, ok := a.Store.RoleByName(*req.Role)
		if !ok {
			writeError(w, http.StatusBadRequest, "unknown role")
			return
		}
		patch.RoleID = &role.ID
	}

	user, err := a.Store.UpdateUser(r.PathValue("id"), patch)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "updated user "+user.ID)

	account, err := a.toUserAccount(user)
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, account)
}

// handleDeleteUser implements DELETE /api/v1/users/{id}.
func (a *App) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.Store.DeleteUser(id); err != nil {
		writeSystemError(w, err)
		return
	}
	logs.FromContext(r.Context()).Log("info", "deleted user "+id)
	writeJSON(w, http.StatusOK, nil)
}
