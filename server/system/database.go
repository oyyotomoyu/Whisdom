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
	chunks        map[string][]MaterialChunk // material ID -> its chunks
	trainingJobs  map[string]*TrainingJob

	trainingMaterialPath string

	// processor turns a stored material's file into embedded chunks in the
	// background. Set by the caller that constructs Store (main.go, or a
	// test) so system stays free of a direct dependency on the ai package.
	processor MaterialProcessor
	// trainer runs a training job's (currently simulated) work in the
	// background, same wiring rationale as processor.
	trainer TrainingRunner
	// logHook records background job outcomes (material processing,
	// training jobs) as system log lines with empty IP/user ID, matching
	// docs/server.md's "Background jobs use system logs" rule. It is nil in
	// tests that don't care about log output.
	logHook func(status, content string)
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
		chunks:               make(map[string][]MaterialChunk),
		trainingJobs:         make(map[string]*TrainingJob),
		trainingMaterialPath: defaultTrainingPath,
	}
}

// SetMaterialProcessor installs the pipeline used to turn an uploaded
// material into searchable chunks. Must be called before any material is
// created; materials created before this is set fail processing immediately.
func (s *Store) SetMaterialProcessor(p MaterialProcessor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processor = p
}

// SetTrainingRunner installs the backend that executes training jobs.
func (s *Store) SetTrainingRunner(t TrainingRunner) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.trainer = t
}

// SetLogHook installs the callback used to record background job outcomes.
func (s *Store) SetLogHook(fn func(status, content string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logHook = fn
}

func (s *Store) log(status, content string) {
	s.mu.RLock()
	fn := s.logHook
	s.mu.RUnlock()
	if fn != nil {
		fn(status, content)
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
