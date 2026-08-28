package system

// Permission is one grantable capability. Handlers require a specific
// permission; roles are just named permission sets.
type Permission = string

const (
	PermChatUse Permission = "chat.use"

	PermConversationsRead   Permission = "conversations.read"
	PermConversationsDelete Permission = "conversations.delete"

	PermMaterialsRead     Permission = "materials.read"
	PermMaterialsUpload   Permission = "materials.upload"
	PermMaterialsDelete   Permission = "materials.delete"
	PermMaterialsPathEdit Permission = "materials.path.edit"

	PermCorrectionsCreate Permission = "corrections.create"
	PermCorrectionsRead   Permission = "corrections.read"
	PermCorrectionsUpdate Permission = "corrections.update"
	PermCorrectionsDelete Permission = "corrections.delete"

	PermUsersRead   Permission = "users.read"
	PermUsersManage Permission = "users.manage"
	PermRolesRead   Permission = "roles.read"
	PermRolesManage Permission = "roles.manage"

	PermConfigRead   Permission = "config.read"
	PermConfigManage Permission = "config.manage"

	PermLogsRead Permission = "logs.read"

	PermModelsRead   Permission = "models.read"
	PermModelsManage Permission = "models.manage"
	PermTrainingRun  Permission = "training.run"

	PermSystemManage Permission = "system.manage"
)

// AllPermissions lists every known permission, used to validate role
// updates and to serve GET /api/v1/permissions.
var AllPermissions = []Permission{
	PermChatUse,
	PermConversationsRead, PermConversationsDelete,
	PermMaterialsRead, PermMaterialsUpload, PermMaterialsDelete, PermMaterialsPathEdit,
	PermCorrectionsCreate, PermCorrectionsRead, PermCorrectionsUpdate, PermCorrectionsDelete,
	PermUsersRead, PermUsersManage, PermRolesRead, PermRolesManage,
	PermConfigRead, PermConfigManage,
	PermLogsRead,
	PermModelsRead, PermModelsManage, PermTrainingRun,
	PermSystemManage,
}

// Default role names seeded on first run, matching the example role table in
// the project README.
const (
	RoleUser          = "user"
	RoleTrainer       = "trainer"
	RoleManager       = "manager"
	RoleAdministrator = "administrator"
)

func defaultRoles() []*Role {
	user := []Permission{PermChatUse, PermConversationsRead, PermConversationsDelete}

	trainer := append(append([]Permission{}, user...),
		PermMaterialsRead, PermMaterialsUpload,
		PermCorrectionsCreate, PermCorrectionsRead,
	)

	manager := append(append([]Permission{}, trainer...),
		PermMaterialsDelete, PermMaterialsPathEdit,
		PermCorrectionsUpdate, PermCorrectionsDelete,
	)

	administrator := append(append([]Permission{}, manager...), AllPermissions...)

	return []*Role{
		{ID: NewID("role"), Name: RoleUser, Permissions: dedupe(user)},
		{ID: NewID("role"), Name: RoleTrainer, Permissions: dedupe(trainer)},
		{ID: NewID("role"), Name: RoleManager, Permissions: dedupe(manager)},
		{ID: NewID("role"), Name: RoleAdministrator, Permissions: dedupe(administrator)},
	}
}

func dedupe(perms []Permission) []Permission {
	seen := make(map[Permission]bool, len(perms))
	out := make([]Permission, 0, len(perms))
	for _, p := range perms {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}
