package system

import "fmt"

// SeedDefaultRoles populates the four built-in roles if no roles exist yet.
// It is idempotent so it is safe to call on every startup.
func (s *Store) SeedDefaultRoles() []*Role {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.roles) > 0 {
		out := make([]*Role, 0, len(s.roles))
		for _, r := range s.roles {
			out = append(out, r)
		}
		return out
	}

	roles := defaultRoles()
	for _, r := range roles {
		s.roles[r.ID] = r
	}
	return roles
}

// ListRoles returns all roles ordered by name.
func (s *Store) ListRoles() []*Role {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return sortedByField(s.roles, func(a, b *Role) bool { return a.Name < b.Name })
}

// GetRole looks up a role by ID.
func (s *Store) GetRole(id string) (*Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.roles[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

// RoleByName looks up a role by its unique name.
func (s *Store) RoleByName(name string) (*Role, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.roles {
		if r.Name == name {
			return r, true
		}
	}
	return nil, false
}

// CreateRole adds a new named permission set.
func (s *Store) CreateRole(name string, permissions []Permission) (*Role, error) {
	if err := validatePermissions(permissions); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.roles {
		if r.Name == name {
			return nil, fmt.Errorf("%w: role name already exists", ErrConflict)
		}
	}

	role := &Role{ID: NewID("role"), Name: name, Permissions: dedupe(permissions)}
	s.roles[role.ID] = role
	return role, nil
}

// UpdateRole patches a role's name and/or permission set.
func (s *Store) UpdateRole(id string, name *string, permissions *[]Permission) (*Role, error) {
	if permissions != nil {
		if err := validatePermissions(*permissions); err != nil {
			return nil, err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	role, ok := s.roles[id]
	if !ok {
		return nil, ErrNotFound
	}
	if name != nil {
		role.Name = *name
	}
	if permissions != nil {
		role.Permissions = dedupe(*permissions)
	}
	return role, nil
}

func validatePermissions(perms []Permission) error {
	valid := make(map[Permission]bool, len(AllPermissions))
	for _, p := range AllPermissions {
		valid[p] = true
	}
	for _, p := range perms {
		if !valid[p] {
			return fmt.Errorf("unknown permission %q", p)
		}
	}
	return nil
}
