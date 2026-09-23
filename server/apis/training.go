package apis

import (
	"net/http"
	"time"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

type trainingJobResponse struct {
	ID            string     `json:"id"`
	MaterialIDs   []string   `json:"materialIds"`
	CorrectionIDs []string   `json:"correctionIds"`
	Status        string     `json:"status"`
	Error         string     `json:"error,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

func toTrainingJobResponse(j *system.TrainingJob) trainingJobResponse {
	resp := trainingJobResponse{
		ID:            j.ID,
		MaterialIDs:   j.MaterialIDs,
		CorrectionIDs: j.CorrectionIDs,
		Status:        string(j.Status),
		Error:         j.Error,
		CreatedAt:     j.CreatedAt,
	}
	if !j.StartedAt.IsZero() {
		resp.StartedAt = &j.StartedAt
	}
	if !j.CompletedAt.IsZero() {
		resp.CompletedAt = &j.CompletedAt
	}
	return resp
}

type createTrainingJobRequest struct {
	MaterialIDs   []string `json:"materialIds"`
	CorrectionIDs []string `json:"correctionIds"`
}

// handleCreateTrainingJob implements POST /api/v1/training/jobs: it prepares
// the selected (already-ready) materials and (already-approved) corrections
// and triggers a training/adapter job. Uploading or correcting alone never
// does this — training only starts from an explicit request here.
func (a *App) handleCreateTrainingJob(w http.ResponseWriter, r *http.Request) {
	var req createTrainingJobRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, _ := currentUser(r.Context())
	job, err := a.Store.CreateTrainingJob(req.MaterialIDs, req.CorrectionIDs, user.ID)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "started training job "+job.ID)
	writeJSON(w, http.StatusCreated, toTrainingJobResponse(job))
}

// handleListTrainingJobs implements GET /api/v1/training/jobs.
func (a *App) handleListTrainingJobs(w http.ResponseWriter, r *http.Request) {
	jobs := a.Store.ListTrainingJobs()
	out := make([]trainingJobResponse, len(jobs))
	for i, j := range jobs {
		out[i] = toTrainingJobResponse(j)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleGetTrainingJob implements GET /api/v1/training/jobs/{id}.
func (a *App) handleGetTrainingJob(w http.ResponseWriter, r *http.Request) {
	job, err := a.Store.GetTrainingJob(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toTrainingJobResponse(job))
}
