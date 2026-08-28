package system

import (
	"fmt"
	"strings"
)

// CreateUser adds a new account. passwordHash must already be hashed by the
// caller (see auth.go); the store never sees plaintext passwords.
func (s *Store) CreateUser(name, email, passwordHash, roleID string) (*User, error) {
	emailKey := strings.ToLower(strings.TrimSpace(email))

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByEmail[emailKey]; exists {
		return nil, fmt.Errorf("%w: email already registered", ErrConflict)
	}
	if _, ok := s.roles[roleID]; !ok {
		return nil, fmt.Errorf("%w: role", ErrNotFound)
	}

	user := &User{
		ID:           NewID("usr"),
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		RoleID:       roleID,
		Active:       true,
		CreatedAt:    timeNow(),
	}
	s.users[user.ID] = user
	s.usersByEmail[emailKey] = user.ID
	return user, nil
}

// ListUsers returns all accounts ordered by name.
func (s *Store) ListUsers() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return sortedByField(s.users, func(a, b *User) bool { return a.Name < b.Name })
}

// GetUser looks up an account by ID.
func (s *Store) GetUser(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

// GetUserByEmail looks up an account by email, case-insensitively.
func (s *Store) GetUserByEmail(email string) (*User, error) {
	emailKey := strings.ToLower(strings.TrimSpace(email))

	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.usersByEmail[emailKey]
	if !ok {
		return nil, ErrNotFound
	}
	return s.users[id], nil
}

// UserPatch describes a partial update; nil fields are left unchanged.
type UserPatch struct {
	Name         *string
	Email        *string
	RoleID       *string
	Active       *bool
	PasswordHash *string
}

// UpdateUser applies a partial update to an existing account.
func (s *Store) UpdateUser(id string, patch UserPatch) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	if patch.RoleID != nil {
		if _, ok := s.roles[*patch.RoleID]; !ok {
			return nil, fmt.Errorf("%w: role", ErrNotFound)
		}
		user.RoleID = *patch.RoleID
	}
	if patch.Email != nil && !strings.EqualFold(*patch.Email, user.Email) {
		emailKey := strings.ToLower(strings.TrimSpace(*patch.Email))
		if _, exists := s.usersByEmail[emailKey]; exists {
			return nil, fmt.Errorf("%w: email already registered", ErrConflict)
		}
		delete(s.usersByEmail, strings.ToLower(user.Email))
		s.usersByEmail[emailKey] = user.ID
		user.Email = *patch.Email
	}
	if patch.Name != nil {
		user.Name = *patch.Name
	}
	if patch.Active != nil {
		user.Active = *patch.Active
	}
	if patch.PasswordHash != nil {
		user.PasswordHash = *patch.PasswordHash
	}

	return user, nil
}

// DeleteUser removes an account.
func (s *Store) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	delete(s.usersByEmail, strings.ToLower(user.Email))
	return nil
}

// ResolveAuthUser expands a User into the public shape carrying its role
// name and resolved permission list.
func (s *Store) ResolveAuthUser(u *User) (*AuthUser, error) {
	role, err := s.GetRole(u.RoleID)
	if err != nil {
		return nil, err
	}
	return &AuthUser{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		Role:        role.Name,
		Permissions: role.Permissions,
	}, nil
}
