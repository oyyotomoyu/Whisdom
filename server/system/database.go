// Package system owns business rules and data access for accounts,
// conversations, materials, and corrections. Store is an in-memory
// implementation; swapping in a real database later means reimplementing
// the same methods on a new type backed by SQL.
package system

import (
	"sort"
	"sync"
	"time"
)

// Store is a thread-safe in-memory data layer.
type Store struct {
	mu sync.RWMutex

	users         map[string]*User
	usersByEmail  map[string]string // lowercased email -> user ID
	roles         map[string]*Role
	conversations map[string]*Conversation
	materials     map[string]*Material
	corrections   map[string]*Correction

	trainingMaterialPath string
}

// NewStore returns an empty Store. Call SeedDefaults to populate default
// roles and an initial administrator account.
func NewStore(defaultTrainingPath string) *Store {
	return &Store{
		users:                make(map[string]*User),
		usersByEmail:         make(map[string]string),
		roles:                make(map[string]*Role),
		conversations:        make(map[string]*Conversation),
		materials:            make(map[string]*Material),
		corrections:          make(map[string]*Correction),
		trainingMaterialPath: defaultTrainingPath,
	}
}

func sortedByField[T any](items map[string]T, less func(a, b T) bool) []T {
	out := make([]T, 0, len(items))
	for _, v := range items {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return less(out[i], out[j]) })
	return out
}

func timeNow() time.Time { return time.Now() }
