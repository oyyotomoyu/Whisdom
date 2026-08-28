package apis

import (
	"net/http"
	"time"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

type correctionResponse struct {
	ID                string    `json:"id"`
	Question          string    `json:"question"`
	OriginalAnswer    string    `json:"originalAnswer"`
	CorrectedAnswer   string    `json:"correctedAnswer"`
	Note              string    `json:"note,omitempty"`
	RelatedMaterialID string    `json:"relatedMaterialId,omitempty"`
	Usage             string    `json:"usage"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
}

func toCorrectionResponse(c *system.Correction) correctionResponse {
	return correctionResponse{
		ID:                c.ID,
		Question:          c.Question,
		OriginalAnswer:    c.OriginalAnswer,
		CorrectedAnswer:   c.CorrectedAnswer,
		Note:              c.Note,
		RelatedMaterialID: c.RelatedMaterialID,
		Usage:             string(c.Usage),
		Status:            string(c.Status),
		CreatedAt:         c.CreatedAt,
	}
}

// handleListCorrections implements GET /api/v1/corrections.
func (a *App) handleListCorrections(w http.ResponseWriter, r *http.Request) {
	corrections := a.Store.ListCorrections()
	out := make([]correctionResponse, len(corrections))
	for i, c := range corrections {
		out[i] = toCorrectionResponse(c)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleGetCorrection implements GET /api/v1/corrections/{id}.
func (a *App) handleGetCorrection(w http.ResponseWriter, r *http.Request) {
	c, err := a.Store.GetCorrection(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCorrectionResponse(c))
}

type createCorrectionRequest struct {
	Question          string `json:"question"`
	OriginalAnswer    string `json:"originalAnswer"`
	CorrectedAnswer   string `json:"correctedAnswer"`
	Note              string `json:"note"`
	RelatedMaterialID string `json:"relatedMaterialId"`
	Usage             string `json:"usage"`
}

var validUsages = map[string]bool{
	string(system.CorrectionUsageRAG): true, string(system.CorrectionUsageEvaluation): true,
	string(system.CorrectionUsageTraining): true, string(system.CorrectionUsageAll): true,
}

// handleCreateCorrection implements POST /api/v1/corrections and
// POST /api/v1/messages/{id}/corrections.
func (a *App) handleCreateCorrection(w http.ResponseWriter, r *http.Request) {
	var req createCorrectionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Question == "" || req.CorrectedAnswer == "" {
		writeError(w, http.StatusBadRequest, "question and correctedAnswer are required")
		return
	}
	if req.Usage == "" {
		req.Usage = string(system.CorrectionUsageAll)
	}
	if !validUsages[req.Usage] {
		writeError(w, http.StatusBadRequest, "usage must be one of rag, evaluation, training, all")
		return
	}

	user, _ := currentUser(r.Context())
	c := a.Store.CreateCorrection(req.Question, req.OriginalAnswer, req.CorrectedAnswer, req.Note, req.RelatedMaterialID, system.CorrectionUsage(req.Usage), user.ID)

	logs.FromContext(r.Context()).Log("info", "created correction "+c.ID)
	writeJSON(w, http.StatusCreated, toCorrectionResponse(c))
}

type updateCorrectionRequest struct {
	CorrectedAnswer *string `json:"correctedAnswer"`
	Note            *string `json:"note"`
	Usage           *string `json:"usage"`
	Status          *string `json:"status"`
}

// handleUpdateCorrection implements PATCH /api/v1/corrections/{id}.
func (a *App) handleUpdateCorrection(w http.ResponseWriter, r *http.Request) {
	var req updateCorrectionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	patch := system.CorrectionPatch{CorrectedAnswer: req.CorrectedAnswer, Note: req.Note}
	if req.Usage != nil {
		if !validUsages[*req.Usage] {
			writeError(w, http.StatusBadRequest, "usage must be one of rag, evaluation, training, all")
			return
		}
		usage := system.CorrectionUsage(*req.Usage)
		patch.Usage = &usage
	}
	if req.Status != nil {
		status := system.CorrectionStatus(*req.Status)
		switch status {
		case system.CorrectionPending, system.CorrectionApproved, system.CorrectionRejected:
			patch.Status = &status
		default:
			writeError(w, http.StatusBadRequest, "status must be one of pending, approved, rejected")
			return
		}
	}

	c, err := a.Store.UpdateCorrection(r.PathValue("id"), patch)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "updated correction "+c.ID)
	writeJSON(w, http.StatusOK, toCorrectionResponse(c))
}

// handleDeleteCorrection implements DELETE /api/v1/corrections/{id}.
func (a *App) handleDeleteCorrection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.Store.DeleteCorrection(id); err != nil {
		writeSystemError(w, err)
		return
	}
	logs.FromContext(r.Context()).Log("info", "deleted correction "+id)
	writeJSON(w, http.StatusOK, nil)
}
