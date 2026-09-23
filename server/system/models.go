package system

import (
	"context"
	"time"
)

// Role is a named set of permissions.
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// User is an authenticated account.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	RoleID       string    `json:"-"`
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuthUser is the public, permission-resolved shape returned by the auth
// APIs and embedded in JWT claims lookups.
type AuthUser struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Role        string       `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// MessageRole distinguishes who authored a conversation message.
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
)

// MessageSource cites a material a RAG answer drew from.
type MessageSource struct {
	MaterialID string `json:"material_id"`
	Name       string `json:"name"`
}

// Message is one turn in a conversation.
type Message struct {
	ID             string          `json:"id"`
	ConversationID string          `json:"-"`
	Role           MessageRole     `json:"role"`
	Content        string          `json:"content"`
	Sources        []MessageSource `json:"sources,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

// Conversation is a chat thread owned by one user.
type Conversation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	Title     string    `json:"title"`
	Messages  []Message `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MaterialStatus tracks a training material through processing.
type MaterialStatus string

const (
	MaterialPending    MaterialStatus = "pending"
	MaterialProcessing MaterialStatus = "processing"
	MaterialReady      MaterialStatus = "ready"
	MaterialFailed     MaterialStatus = "failed"
)

// MaterialSourceType identifies how a material entered the library.
type MaterialSourceType string

const (
	MaterialSourceUpload    MaterialSourceType = "upload"
	MaterialSourceLocalPath MaterialSourceType = "local_path"
	MaterialSourceWikiURL   MaterialSourceType = "wiki_url"
	MaterialSourceGitURL    MaterialSourceType = "git_url"
	MaterialSourceURL       MaterialSourceType = "url"
)

// Material is one piece of company knowledge.
type Material struct {
	ID                string             `json:"id"`
	Filename          string             `json:"filename"`
	Type              string             `json:"type"`
	SourceType        MaterialSourceType `json:"source_type"`
	SourceLocation    string             `json:"source_location"`
	UploadedByID      string             `json:"-"`
	UploadedByName    string             `json:"uploaded_by"`
	UploadedAt        time.Time          `json:"uploaded_at"`
	DestinationPath   string             `json:"destination_path"`
	Status            MaterialStatus     `json:"status"`
	RAGAvailable      bool               `json:"rag_available"`
	TrainingAvailable bool               `json:"training_available"`
	StoragePath       string             `json:"-"`
}

// CorrectionUsage is where a correction may be applied.
type CorrectionUsage string

const (
	CorrectionUsageRAG        CorrectionUsage = "rag"
	CorrectionUsageEvaluation CorrectionUsage = "evaluation"
	CorrectionUsageTraining   CorrectionUsage = "training"
	CorrectionUsageAll        CorrectionUsage = "all"
)

// CorrectionStatus tracks review of a submitted correction.
type CorrectionStatus string

const (
	CorrectionPending  CorrectionStatus = "pending"
	CorrectionApproved CorrectionStatus = "approved"
	CorrectionRejected CorrectionStatus = "rejected"
)

// Correction records a human-supplied fix to an AI answer.
type Correction struct {
	ID                string           `json:"id"`
	Question          string           `json:"question"`
	OriginalAnswer    string           `json:"original_answer"`
	CorrectedAnswer   string           `json:"corrected_answer"`
	Note              string           `json:"note,omitempty"`
	RelatedMaterialID string           `json:"related_material_id,omitempty"`
	Usage             CorrectionUsage  `json:"usage"`
	Status            CorrectionStatus `json:"status"`
	CreatedByID       string           `json:"-"`
	CreatedAt         time.Time        `json:"created_at"`
}

// MaterialChunk is one embedded slice of a material's extracted text, used
// for RAG similarity search.
type MaterialChunk struct {
	ID         string
	MaterialID string
	Text       string
	Embedding  []float64
}

// ProcessedChunk is one chunk produced by a MaterialProcessor, before Store
// assigns it an ID and a MaterialID.
type ProcessedChunk struct {
	Text      string
	Embedding []float64
}

// MaterialProcessor turns a stored material's file into embedded chunks.
// Implementations live in the ai package (extraction, chunking, embedding);
// system only depends on this interface to avoid an import cycle.
type MaterialProcessor interface {
	Process(ctx context.Context, m *Material) ([]ProcessedChunk, error)
}

// TrainingJobStatus tracks a training/adapter job's lifecycle.
type TrainingJobStatus string

const (
	TrainingJobPending   TrainingJobStatus = "pending"
	TrainingJobRunning   TrainingJobStatus = "running"
	TrainingJobSucceeded TrainingJobStatus = "succeeded"
	TrainingJobFailed    TrainingJobStatus = "failed"
)

// TrainingJob is one request to train or adapt the model from approved
// corrections and selected materials.
type TrainingJob struct {
	ID            string
	MaterialIDs   []string
	CorrectionIDs []string
	Status        TrainingJobStatus
	Error         string
	CreatedByID   string
	CreatedAt     time.Time
	StartedAt     time.Time
	CompletedAt   time.Time
}

// TrainingRunner executes a training job's actual work. Implementations live
// in the ai package, same rationale as MaterialProcessor.
type TrainingRunner interface {
	Run(ctx context.Context, job TrainingJob, materials []*Material, corrections []*Correction) error
}

// ModelTarget is one selectable model the chat/summarize APIs can run
// against.
type ModelTarget struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Description string `json:"description,omitempty"`
}
