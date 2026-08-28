package system

import "testing"

func TestDefaultRolesHierarchy(t *testing.T) {
	store := NewStore("/training/default/")
	roles := store.SeedDefaultRoles()

	byName := make(map[string]*Role, len(roles))
	for _, r := range roles {
		byName[r.Name] = r
	}

	for _, name := range []string{RoleUser, RoleTrainer, RoleManager, RoleAdministrator} {
		if _, ok := byName[name]; !ok {
			t.Fatalf("missing seeded role %q", name)
		}
	}

	has := func(perms []Permission, want Permission) bool {
		for _, p := range perms {
			if p == want {
				return true
			}
		}
		return false
	}

	if !has(byName[RoleUser].Permissions, PermChatUse) {
		t.Error("user role must have chat.use")
	}
	if has(byName[RoleUser].Permissions, PermMaterialsUpload) {
		t.Error("user role must not have materials.upload")
	}
	if !has(byName[RoleTrainer].Permissions, PermMaterialsUpload) {
		t.Error("trainer role must have materials.upload")
	}
	if has(byName[RoleTrainer].Permissions, PermMaterialsDelete) {
		t.Error("trainer role must not have materials.delete")
	}
	if !has(byName[RoleManager].Permissions, PermMaterialsDelete) {
		t.Error("manager role must have materials.delete")
	}
	if has(byName[RoleManager].Permissions, PermUsersManage) {
		t.Error("manager role must not have users.manage (administration)")
	}
	for _, p := range AllPermissions {
		if !has(byName[RoleAdministrator].Permissions, p) {
			t.Errorf("administrator role is missing permission %q", p)
		}
	}
}

func TestSeedDefaultRolesIsIdempotent(t *testing.T) {
	store := NewStore("/training/default/")
	first := store.SeedDefaultRoles()
	second := store.SeedDefaultRoles()
	if len(first) != len(second) {
		t.Fatalf("SeedDefaultRoles not idempotent: got %d roles then %d", len(first), len(second))
	}
	if len(store.ListRoles()) != len(first) {
		t.Fatalf("SeedDefaultRoles created duplicate roles: store has %d, expected %d", len(store.ListRoles()), len(first))
	}
}

func TestCreateRoleRejectsUnknownPermission(t *testing.T) {
	store := NewStore("/training/default/")
	if _, err := store.CreateRole("custom", []Permission{"not.a.real.permission"}); err == nil {
		t.Error("CreateRole should reject an unknown permission")
	}
}

func TestCreateRoleRejectsDuplicateName(t *testing.T) {
	store := NewStore("/training/default/")
	if _, err := store.CreateRole("custom", []Permission{PermChatUse}); err != nil {
		t.Fatalf("CreateRole: %v", err)
	}
	if _, err := store.CreateRole("custom", []Permission{PermChatUse}); err == nil {
		t.Error("CreateRole should reject a duplicate role name")
	}
}
