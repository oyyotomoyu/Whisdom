package system

import (
	"context"
	"errors"
	"time"
)

// trainingJobTimeout bounds how long a training job's (currently simulated)
// work may run before it's marked failed.
const trainingJobTimeout = 5 * time.Minute

// CreateTrainingJob validates the selected materials and corrections, marks
// the materials as training-available (upload alone never does this — see
// docs/server.md's "Uploading material should not automatically fine-tune
// the model"), and starts the job in the background.
func (s *Store) CreateTrainingJob(materialIDs, correctionIDs []string, createdByID string) (*TrainingJob, error) {
	s.mu.Lock()
	if len(materialIDs) == 0 && len(correctionIDs) == 0 {
		s.mu.Unlock()
		return nil, errors.New("at least one material or correction is required")
	}

	for _, id := range materialIDs {
		m, ok := s.materials[id]
		if !ok {
			s.mu.Unlock()
			return nil, ErrNotFound
		}
		if m.Status != MaterialReady {
			s.mu.Unlock()
			return nil, errors.New("material " + id + " is not ready")
		}
	}
	for _, id := range correctionIDs {
		c, ok := s.corrections[id]
		if !ok {
			s.mu.Unlock()
			return nil, ErrNotFound
		}
		if c.Status != CorrectionApproved {
			s.mu.Unlock()
			return nil, errors.New("correction " + id + " is not approved")
		}
		if c.Usage != CorrectionUsageTraining && c.Usage != CorrectionUsageAll {
			s.mu.Unlock()
			return nil, errors.New("correction " + id + " is not marked for training use")
		}
	}

	for _, id := range materialIDs {
		s.materials[id].TrainingAvailable = true
	}

	job := &TrainingJob{
		ID:            NewID("job"),
		MaterialIDs:   append([]string{}, materialIDs...),
		CorrectionIDs: append([]string{}, correctionIDs...),
		Status:        TrainingJobPending,
		CreatedByID:   createdByID,
		CreatedAt:     timeNow(),
	}
	s.trainingJobs[job.ID] = job
	s.mu.Unlock()

	s.runTrainingJob(job.ID)
	return job, nil
}

// ListTrainingJobs returns all training jobs, newest first.
func (s *Store) ListTrainingJobs() []*TrainingJob {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*TrainingJob, 0, len(s.trainingJobs))
	for _, j := range s.trainingJobs {
		out = append(out, j)
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].CreatedAt.After(out[j-1].CreatedAt); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// GetTrainingJob looks up a training job by ID.
func (s *Store) GetTrainingJob(id string) (*TrainingJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.trainingJobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return j, nil
}

func (s *Store) runTrainingJob(id string) {
	s.mu.Lock()
	job, ok := s.trainingJobs[id]
	if !ok {
		s.mu.Unlock()
		return
	}
	job.Status = TrainingJobRunning
	job.StartedAt = timeNow()
	runner := s.trainer

	materials := make([]*Material, 0, len(job.MaterialIDs))
	for _, mid := range job.MaterialIDs {
		if m, ok := s.materials[mid]; ok {
			materials = append(materials, m)
		}
	}
	corrections := make([]*Correction, 0, len(job.CorrectionIDs))
	for _, cid := range job.CorrectionIDs {
		if c, ok := s.corrections[cid]; ok {
			corrections = append(corrections, c)
		}
	}
	jobSnapshot := *job
	s.mu.Unlock()

	s.log("info", "training job "+id+" started")

	go func() {
		var err error
		if runner == nil {
			err = errors.New("no training runner configured")
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), trainingJobTimeout)
			defer cancel()
			err = runner.Run(ctx, jobSnapshot, materials, corrections)
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		job, ok := s.trainingJobs[id]
		if !ok {
			return
		}
		job.CompletedAt = timeNow()
		if err != nil {
			job.Status = TrainingJobFailed
			job.Error = err.Error()
			s.log("error", "training job "+id+" failed: "+err.Error())
			return
		}
		job.Status = TrainingJobSucceeded
		s.log("info", "training job "+id+" succeeded")
	}()
}
