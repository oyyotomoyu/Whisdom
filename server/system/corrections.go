package system

// CreateCorrection records a human-supplied fix to an AI answer, pending
// review.
func (s *Store) CreateCorrection(question, originalAnswer, correctedAnswer, note, relatedMaterialID string, usage CorrectionUsage, createdByID string) *Correction {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := &Correction{
		ID:                NewID("cor"),
		Question:          question,
		OriginalAnswer:    originalAnswer,
		CorrectedAnswer:   correctedAnswer,
		Note:              note,
		RelatedMaterialID: relatedMaterialID,
		Usage:             usage,
		Status:            CorrectionPending,
		CreatedByID:       createdByID,
		CreatedAt:         timeNow(),
	}
	s.corrections[c.ID] = c
	return c
}

// ListCorrections returns all corrections, newest first.
func (s *Store) ListCorrections() []*Correction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*Correction, 0, len(s.corrections))
	for _, c := range s.corrections {
		out = append(out, c)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].CreatedAt.After(out[j-1].CreatedAt); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// GetCorrection looks up a correction by ID.
func (s *Store) GetCorrection(id string) (*Correction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.corrections[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// CorrectionPatch describes a partial update; nil fields are left unchanged.
type CorrectionPatch struct {
	CorrectedAnswer *string
	Note            *string
	Usage           *CorrectionUsage
	Status          *CorrectionStatus
}

// UpdateCorrection applies a partial update, typically a review decision.
func (s *Store) UpdateCorrection(id string, patch CorrectionPatch) (*Correction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, ok := s.corrections[id]
	if !ok {
		return nil, ErrNotFound
	}
	if patch.CorrectedAnswer != nil {
		c.CorrectedAnswer = *patch.CorrectedAnswer
	}
	if patch.Note != nil {
		c.Note = *patch.Note
	}
	if patch.Usage != nil {
		c.Usage = *patch.Usage
	}
	if patch.Status != nil {
		c.Status = *patch.Status
	}
	return c, nil
}

// DeleteCorrection removes a correction.
func (s *Store) DeleteCorrection(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.corrections[id]; !ok {
		return ErrNotFound
	}
	delete(s.corrections, id)
	return nil
}
