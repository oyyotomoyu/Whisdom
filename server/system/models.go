package system

import "time"

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

// Material is one uploaded piece of company knowledge.
type Material struct {
	ID                string         `json:"id"`
	Filename          string         `json:"filename"`
	Type              string         `json:"type"`
	UploadedByID      string         `json:"-"`
	UploadedByName    string         `json:"uploaded_by"`
	UploadedAt        time.Time      `json:"uploaded_at"`
	DestinationPath   string         `json:"destination_path"`
	Status            MaterialStatus `json:"status"`
	RAGAvailable      bool           `json:"rag_available"`
	TrainingAvailable bool           `json:"training_available"`
	StoragePath       string         `json:"-"`
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
