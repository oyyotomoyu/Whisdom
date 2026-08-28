package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/whisdom/server/ai"
	"github.com/whisdom/server/config"
	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

func newTestApp(t *testing.T) (*App, http.Handler) {
	t.Helper()

	logSvc, err := logs.NewService(t.TempDir())
	if err != nil {
		t.Fatalf("logs.NewService: %v", err)
	}
	t.Cleanup(func() { logSvc.Close() })

	store := system.NewStore("/training/default/")
	store.SeedDefaultRoles()

	app := &App{
		Store: store,
		Logs:  logSvc,
		Config: &config.Config{
			JWTSecret:       "test-secret",
			AccessTokenTTL:  time.Minute,
			RefreshTokenTTL: time.Hour,
			DataDir:         t.TempDir(),
		},
		Model: ai.NewStubModel(),
		RAG:   ai.NewMaterialFilenameRAG(store),
	}
	return app, NewRouter(app)
}

func createTestUser(t *testing.T, app *App, roleName, email, password string) *system.User {
	t.Helper()
	role, ok := app.Store.RoleByName(roleName)
	if !ok {
		t.Fatalf("role %q not seeded", roleName)
	}
	hash, err := system.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	user, err := app.Store.CreateUser("Test User", email, hash, role.ID)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	return user
}

func doJSON(t *testing.T, handler http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func loginAndGetToken(t *testing.T, handler http.Handler, email, password string) string {
	t.Helper()
	rec := doJSON(t, handler, http.MethodPost, "/api/v1/auth/login", "", loginRequest{Username: email, Password: password})
	if rec.Code != http.StatusOK {
		t.Fatalf("login failed: status=%d body=%s", rec.Code, rec.Body.String())
	}
	var resp loginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	return resp.AccessToken
}

func TestHealthIsPublic(t *testing.T) {
	_, handler := newTestApp(t)
	rec := doJSON(t, handler, http.MethodGet, "/api/v1/health", "", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health = %d, want 200", rec.Code)
	}
}

func TestConversationsRequireAuth(t *testing.T) {
	_, handler := newTestApp(t)
	rec := doJSON(t, handler, http.MethodGet, "/api/v1/conversations", "", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /conversations without token = %d, want 401", rec.Code)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	app, handler := newTestApp(t)
	createTestUser(t, app, system.RoleUser, "user@example.com", "correct-password")

	rec := doJSON(t, handler, http.MethodPost, "/api/v1/auth/login", "", loginRequest{Username: "user@example.com", Password: "wrong"})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("login with wrong password = %d, want 401", rec.Code)
	}
}

func TestRegularUserCannotAccessMaterials(t *testing.T) {
	app, handler := newTestApp(t)
	createTestUser(t, app, system.RoleUser, "user@example.com", "correct-password")
	token := loginAndGetToken(t, handler, "user@example.com", "correct-password")

	rec := doJSON(t, handler, http.MethodGet, "/api/v1/materials", token, nil)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("regular user GET /materials = %d, want 403", rec.Code)
	}
}

func TestAdministratorCanAccessMaterials(t *testing.T) {
	app, handler := newTestApp(t)
	createTestUser(t, app, system.RoleAdministrator, "admin@example.com", "correct-password")
	token := loginAndGetToken(t, handler, "admin@example.com", "correct-password")

	rec := doJSON(t, handler, http.MethodGet, "/api/v1/materials", token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("admin GET /materials = %d, want 200", rec.Code)
	}
}

func TestChatFlowEndToEnd(t *testing.T) {
	app, handler := newTestApp(t)
	createTestUser(t, app, system.RoleUser, "user@example.com", "correct-password")
	token := loginAndGetToken(t, handler, "user@example.com", "correct-password")

	rec := doJSON(t, handler, http.MethodPost, "/api/v1/conversations", token, nil)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /conversations = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	var conv conversationSummaryResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &conv); err != nil {
		t.Fatalf("decode conversation: %v", err)
	}

	rec = doJSON(t, handler, http.MethodPost, "/api/v1/conversations/"+conv.ID+"/messages", token, sendMessageRequest{Message: "How should we handle a refund?"})
	if rec.Code != http.StatusOK {
		t.Fatalf("POST messages = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	var chat chatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &chat); err != nil {
		t.Fatalf("decode chat response: %v", err)
	}
	if chat.Answer == "" {
		t.Error("expected a non-empty answer")
	}

	rec = doJSON(t, handler, http.MethodGet, "/api/v1/conversations/"+conv.ID, token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET conversation = %d, want 200", rec.Code)
	}
	var detail conversationDetailResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatalf("decode conversation detail: %v", err)
	}
	if len(detail.Messages) != 2 {
		t.Fatalf("expected 2 messages (user + assistant), got %d", len(detail.Messages))
	}
	if detail.Messages[0].Role != system.MessageRoleUser || detail.Messages[1].Role != system.MessageRoleAssistant {
		t.Errorf("unexpected message roles: %v, %v", detail.Messages[0].Role, detail.Messages[1].Role)
	}
}

func TestConfigPathValidation(t *testing.T) {
	app, handler := newTestApp(t)
	createTestUser(t, app, system.RoleAdministrator, "admin@example.com", "correct-password")
	token := loginAndGetToken(t, handler, "admin@example.com", "correct-password")

	rec := doJSON(t, handler, http.MethodPatch, "/api/v1/config", token, updateConfigRequest{TrainingMaterialPath: "not-absolute"})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("PATCH /config with invalid path = %d, want 400", rec.Code)
	}

	rec = doJSON(t, handler, http.MethodPatch, "/api/v1/config", token, updateConfigRequest{TrainingMaterialPath: "/training/support/"})
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /config with valid path = %d, want 200: %s", rec.Code, rec.Body.String())
	}
}
