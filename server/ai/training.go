package ai

import (
	"context"
	"time"

	"github.com/whisdom/server/system"
)

// simulatedTrainingDuration stands in for a real fine-tuning/adapter run.
// There is no local training backend wired in yet, so this lets the job API,
// status storage, and logging (docs/server.md's "Training job start,
// success, and failure") be exercised end to end today.
const simulatedTrainingDuration = 2 * time.Second

// SimulatedTrainingRunner is a placeholder system.TrainingRunner. Swapping
// in a real fine-tuning/adapter backend means implementing the same
// interface against that backend's job API; callers never need to change.
type SimulatedTrainingRunner struct{}

// NewSimulatedTrainingRunner returns a placeholder TrainingRunner.
func NewSimulatedTrainingRunner() *SimulatedTrainingRunner { return &SimulatedTrainingRunner{} }

// Run implements system.TrainingRunner.
func (r *SimulatedTrainingRunner) Run(ctx context.Context, job system.TrainingJob, materials []*system.Material, corrections []*system.Correction) error {
	select {
	case <-time.After(simulatedTrainingDuration):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
