package apis

import (
	"net/http"

	"github.com/whisdom/server/system"
)

// NewRouter builds the full /api/v1 route table. Each protected route names
// the exact permission(s) docs/server.md assigns it; requireAuth always runs
// first so handlers can assume an authenticated user is present.
func NewRouter(a *App) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/health", a.public(a.handleHealth))

	mux.Handle("POST /api/v1/auth/login", a.public(a.handleLogin))
	mux.Handle("POST /api/v1/auth/refresh", a.public(a.handleRefresh))
	mux.Handle("POST /api/v1/auth/logout", a.protected(a.handleLogout))
	mux.Handle("GET /api/v1/auth/me", a.protected(a.handleMe))

	mux.Handle("GET /api/v1/users", a.protected(a.handleListUsers, system.PermUsersRead))
	mux.Handle("GET /api/v1/users/{id}", a.protected(a.handleGetUser, system.PermUsersRead))
	mux.Handle("POST /api/v1/users", a.protected(a.handleCreateUser, system.PermUsersManage))
	mux.Handle("PATCH /api/v1/users/{id}", a.protected(a.handleUpdateUser, system.PermUsersManage))
	mux.Handle("DELETE /api/v1/users/{id}", a.protected(a.handleDeleteUser, system.PermUsersManage))

	mux.Handle("GET /api/v1/roles", a.protected(a.handleListRoles, system.PermRolesRead))
	mux.Handle("POST /api/v1/roles", a.protected(a.handleCreateRole, system.PermRolesManage))
	mux.Handle("PATCH /api/v1/roles/{id}", a.protected(a.handleUpdateRole, system.PermRolesManage))
	mux.Handle("GET /api/v1/permissions", a.protected(a.handleListPermissions, system.PermRolesRead))

	mux.Handle("POST /api/v1/conversations", a.protected(a.handleCreateConversation, system.PermChatUse))
	mux.Handle("GET /api/v1/conversations", a.protected(a.handleListConversations, system.PermConversationsRead))
	mux.Handle("GET /api/v1/conversations/{id}", a.protected(a.handleGetConversation, system.PermConversationsRead))
	mux.Handle("DELETE /api/v1/conversations/{id}", a.protected(a.handleDeleteConversation, system.PermConversationsDelete))
	mux.Handle("POST /api/v1/conversations/{id}/messages", a.protected(a.handleSendMessage, system.PermChatUse))

	mux.Handle("GET /api/v1/materials", a.protected(a.handleListMaterials, system.PermMaterialsRead))
	mux.Handle("GET /api/v1/materials/{id}", a.protected(a.handleGetMaterial, system.PermMaterialsRead))
	mux.Handle("GET /api/v1/materials/{id}/status", a.protected(a.handleMaterialStatus, system.PermMaterialsRead))
	mux.Handle("POST /api/v1/materials", a.protected(a.handleUploadMaterial, system.PermMaterialsUpload))
	mux.Handle("PATCH /api/v1/materials/{id}", a.protected(a.handleUpdateMaterial, system.PermMaterialsPathEdit))
	mux.Handle("DELETE /api/v1/materials/{id}", a.protected(a.handleDeleteMaterial, system.PermMaterialsDelete))
	// Both spellings are wired to the same handler: docs/server.md documents
	// "/reprocess" while the README example (and the React client) use
	// "/process".
	mux.Handle("POST /api/v1/materials/{id}/process", a.protected(a.handleReprocessMaterial, system.PermMaterialsUpload))
	mux.Handle("POST /api/v1/materials/{id}/reprocess", a.protected(a.handleReprocessMaterial, system.PermMaterialsUpload))

	mux.Handle("GET /api/v1/corrections", a.protected(a.handleListCorrections, system.PermCorrectionsRead))
	mux.Handle("GET /api/v1/corrections/{id}", a.protected(a.handleGetCorrection, system.PermCorrectionsRead))
	mux.Handle("POST /api/v1/corrections", a.protected(a.handleCreateCorrection, system.PermCorrectionsCreate))
	mux.Handle("POST /api/v1/messages/{id}/corrections", a.protected(a.handleCreateCorrection, system.PermCorrectionsCreate))
	mux.Handle("PATCH /api/v1/corrections/{id}", a.protected(a.handleUpdateCorrection, system.PermCorrectionsUpdate))
	mux.Handle("DELETE /api/v1/corrections/{id}", a.protected(a.handleDeleteCorrection, system.PermCorrectionsDelete))

	mux.Handle("GET /api/v1/config", a.protected(a.handleGetConfig, system.PermConfigRead))
	// PATCH /api/v1/config accepts config.manage OR materials.path.edit; that
	// OR can't be expressed by protected()'s AND-only permission list, so the
	// handler itself checks hasAnyPermission after plain authentication.
	mux.Handle("PATCH /api/v1/config", a.protected(a.handleUpdateConfig))

	mux.Handle("GET /api/v1/logs", a.protected(a.handleListLogs, system.PermLogsRead))

	return mux
}
