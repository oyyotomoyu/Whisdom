package system

import "github.com/whisdom/server/storage"

// TrainingMaterialPath returns the current default destination for newly
// uploaded training material.
func (s *Store) TrainingMaterialPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.trainingMaterialPath
}

// SetTrainingMaterialPath validates, normalizes, and stores a new default
// destination path. Callers are responsible for checking the acting user
// holds config.manage or materials.path.edit and for audit-logging the
// change.
func (s *Store) SetTrainingMaterialPath(path string) (string, error) {
	if err := storage.ValidateDestinationPath(path); err != nil {
		return "", err
	}
	normalized := storage.NormalizeDestinationPath(path)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.trainingMaterialPath = normalized
	return normalized, nil
}
