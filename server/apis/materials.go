package apis

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/storage"
	"github.com/whisdom/server/system"
)

// maxUploadSize caps a single material upload at 200MB (mainly audio/video).
const maxUploadSize = 200 << 20

type materialResponse struct {
	ID                string    `json:"id"`
	Filename          string    `json:"filename"`
	Type              string    `json:"type"`
	SourceType        string    `json:"sourceType"`
	SourceLocation    string    `json:"sourceLocation"`
	UploadedAt        time.Time `json:"uploadedAt"`
	UploadedBy        string    `json:"uploadedBy"`
	Status            string    `json:"status"`
	DestinationPath   string    `json:"destinationPath"`
	RAGAvailable      bool      `json:"ragAvailable"`
	TrainingAvailable bool      `json:"trainingAvailable"`
}

func toMaterialResponse(m *system.Material) materialResponse {
	return materialResponse{
		ID:                m.ID,
		Filename:          m.Filename,
		Type:              m.Type,
		SourceType:        string(m.SourceType),
		SourceLocation:    m.SourceLocation,
		UploadedAt:        m.UploadedAt,
		UploadedBy:        m.UploadedByName,
		Status:            string(m.Status),
		DestinationPath:   m.DestinationPath,
		RAGAvailable:      m.RAGAvailable,
		TrainingAvailable: m.TrainingAvailable,
	}
}

// handleListMaterials implements GET /api/v1/materials.
func (a *App) handleListMaterials(w http.ResponseWriter, r *http.Request) {
	materials := a.Store.ListMaterials()
	out := make([]materialResponse, len(materials))
	for i, m := range materials {
		out[i] = toMaterialResponse(m)
	}
	writeJSON(w, http.StatusOK, out)
}

// handleGetMaterial implements GET /api/v1/materials/{id}.
func (a *App) handleGetMaterial(w http.ResponseWriter, r *http.Request) {
	m, err := a.Store.GetMaterial(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toMaterialResponse(m))
}

// handleUploadMaterial implements POST /api/v1/materials. The client sends
// either multipart form data with a "file" part or JSON source metadata for
// already-existing organization knowledge locations.
func (a *App) handleUploadMaterial(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		a.handleCreateMaterialSource(w, r)
		return
	}

	logger := logs.FromContext(r.Context())

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "invalid or too large upload")
		return
	}

	destinationPath := r.FormValue("destination_path")
	if err := storage.ValidateDestinationPath(destinationPath); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	destinationPath = storage.NormalizeDestinationPath(destinationPath)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	user, _ := currentUser(r.Context())
	id := system.NewID("mat")

	storagePath, err := storage.SaveMaterialFile(a.Config.DataDir, id, header.Filename, file)
	if err != nil {
		logger.Log("error", "failed to store material file: "+err.Error())
		writeError(w, http.StatusInternalServerError, "failed to store file")
		return
	}

	materialType := strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), ".")
	if materialType == "" {
		materialType = "file"
	}

	sourceLocation := r.FormValue("source_location")
	if strings.TrimSpace(sourceLocation) == "" {
		sourceLocation = storage.SanitizeFilename(header.Filename)
	}

	m := a.Store.CreateMaterialFromSource(
		id,
		storage.SanitizeFilename(header.Filename),
		materialType,
		user.ID,
		user.Name,
		destinationPath,
		string(system.MaterialSourceUpload),
		sourceLocation,
		storagePath,
	)

	logger.Log("info", "uploaded material "+m.ID)
	writeJSON(w, http.StatusCreated, toMaterialResponse(m))
}

type createMaterialSourceRequest struct {
	SourceType      string `json:"source_type"`
	SourceLocation  string `json:"source_location"`
	DestinationPath string `json:"destination_path"`
}

func (a *App) handleCreateMaterialSource(w http.ResponseWriter, r *http.Request) {
	var req createMaterialSourceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	destinationPath := req.DestinationPath
	if err := storage.ValidateDestinationPath(destinationPath); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	destinationPath = storage.NormalizeDestinationPath(destinationPath)

	sourceType := strings.TrimSpace(req.SourceType)
	sourceLocation := strings.TrimSpace(req.SourceLocation)
	if err := validateMaterialSource(sourceType, sourceLocation); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, _ := currentUser(r.Context())
	id := system.NewID("mat")
	filename := materialNameFromSource(sourceType, sourceLocation)

	m := a.Store.CreateMaterialFromSource(
		id,
		filename,
		materialTypeFromSource(sourceType, sourceLocation),
		user.ID,
		user.Name,
		destinationPath,
		sourceType,
		sourceLocation,
		"",
	)

	logs.FromContext(r.Context()).Log("info", "added material source "+m.ID)
	writeJSON(w, http.StatusCreated, toMaterialResponse(m))
}

func validateMaterialSource(sourceType, sourceLocation string) error {
	if sourceType == "" {
		return errors.New("source_type is required")
	}
	if strings.TrimSpace(sourceLocation) == "" {
		return errors.New("source_location is required")
	}
	if strings.ContainsRune(sourceLocation, 0) || strings.Contains(sourceLocation, "..") {
		return errors.New("source_location is invalid")
	}

	switch system.MaterialSourceType(sourceType) {
	case system.MaterialSourceLocalPath:
		if !strings.HasPrefix(sourceLocation, "/") {
			return errors.New("local source path must start with /")
		}
	case system.MaterialSourceWikiURL, system.MaterialSourceGitURL, system.MaterialSourceURL:
		parsed, err := url.ParseRequestURI(sourceLocation)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return errors.New("source_location must be a valid URL")
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return errors.New("source URL must use http or https")
		}
	default:
		return errors.New("source_type is unsupported")
	}
	return nil
}

func materialNameFromSource(sourceType, sourceLocation string) string {
	if system.MaterialSourceType(sourceType) == system.MaterialSourceLocalPath {
		return storage.SanitizeFilename(path.Base(sourceLocation))
	}
	parsed, err := url.Parse(sourceLocation)
	if err == nil {
		if base := path.Base(parsed.Path); base != "." && base != "/" {
			return storage.SanitizeFilename(base)
		}
		if parsed.Host != "" {
			return storage.SanitizeFilename(parsed.Host)
		}
	}
	return storage.SanitizeFilename(sourceLocation)
}

func materialTypeFromSource(sourceType, sourceLocation string) string {
	if ext := strings.TrimPrefix(strings.ToLower(path.Ext(sourceLocation)), "."); ext != "" {
		return ext
	}
	return sourceType
}

type updateMaterialRequest struct {
	DestinationPath string `json:"destination_path"`
}

// handleUpdateMaterial implements PATCH /api/v1/materials/{id}.
func (a *App) handleUpdateMaterial(w http.ResponseWriter, r *http.Request) {
	var req updateMaterialRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := storage.ValidateDestinationPath(req.DestinationPath); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	normalized := storage.NormalizeDestinationPath(req.DestinationPath)

	m, err := a.Store.UpdateMaterialDestination(r.PathValue("id"), normalized)
	if err != nil {
		writeSystemError(w, err)
		return
	}

	logs.FromContext(r.Context()).Log("info", "updated destination path for material "+m.ID)
	writeJSON(w, http.StatusOK, toMaterialResponse(m))
}

// handleDeleteMaterial implements DELETE /api/v1/materials/{id}.
func (a *App) handleDeleteMaterial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	logger := logs.FromContext(r.Context())

	m, err := a.Store.DeleteMaterial(id)
	if err != nil {
		writeSystemError(w, err)
		return
	}
	if err := storage.DeleteMaterialDir(a.Config.DataDir, m.ID); err != nil {
		logger.Log("warning", "failed to delete material file for "+m.ID+": "+err.Error())
	}

	logger.Log("info", "deleted material "+id)
	writeJSON(w, http.StatusOK, nil)
}

// handleReprocessMaterial implements POST /api/v1/materials/{id}/process.
func (a *App) handleReprocessMaterial(w http.ResponseWriter, r *http.Request) {
	m, err := a.Store.ReprocessMaterial(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}
	logs.FromContext(r.Context()).Log("info", "reprocessing material "+m.ID)
	writeJSON(w, http.StatusOK, toMaterialResponse(m))
}

type materialProcessingStatusResponse struct {
	ID         string `json:"id"`
	Filename   string `json:"filename"`
	Status     string `json:"status"`
	Processing struct {
		Extraction string `json:"extraction"`
		Chunking   string `json:"chunking"`
		Embedding  string `json:"embedding"`
	} `json:"processing"`
}

// handleMaterialStatus implements GET /api/v1/materials/{id}/status. Since
// there is no real extraction/chunking/embedding pipeline yet, every stage
// is reported as a single step matching the material's overall status.
func (a *App) handleMaterialStatus(w http.ResponseWriter, r *http.Request) {
	m, err := a.Store.GetMaterial(r.PathValue("id"))
	if err != nil {
		writeSystemError(w, err)
		return
	}

	stage := "pending"
	switch m.Status {
	case system.MaterialProcessing:
		stage = "in_progress"
	case system.MaterialReady:
		stage = "completed"
	case system.MaterialFailed:
		stage = "failed"
	}

	resp := materialProcessingStatusResponse{ID: m.ID, Filename: m.Filename, Status: string(m.Status)}
	resp.Processing.Extraction = stage
	resp.Processing.Chunking = stage
	resp.Processing.Embedding = stage
	writeJSON(w, http.StatusOK, resp)
}
