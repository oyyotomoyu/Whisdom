// Command server is Whisdom's Go backend: it owns authentication,
// permissions, conversation storage, material handling, and (eventually)
// model/RAG orchestration behind a REST API consumed by the React frontend.
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/whisdom/server/ai"
	"github.com/whisdom/server/apis"
	"github.com/whisdom/server/config"
	"github.com/whisdom/server/logs"
	"github.com/whisdom/server/system"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load("config/default.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if generated, err := cfg.EnsureJWTSecret(); err != nil {
		return err
	} else if generated {
		log.Println("warning: WHISDOM_JWT_SECRET not set; generated an ephemeral signing key. " +
			"All sessions will be invalidated on restart. Set WHISDOM_JWT_SECRET for production.")
	}

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	logSvc, err := logs.NewService(cfg.LogDir)
	if err != nil {
		return fmt.Errorf("start log service: %w", err)
	}
	defer logSvc.Close()

	stopRetention := make(chan struct{})
	defer close(stopRetention)
	logSvc.StartRetentionLoop(stopRetention)

	logSvc.Log(logs.StatusInfo, "", "", "server starting")

	store := system.NewStore(cfg.TrainingMaterialPath)
	store.SeedDefaultRoles()
	if err := seedAdministrator(store, logSvc); err != nil {
		return fmt.Errorf("seed administrator: %w", err)
	}

	embedder := newEmbedder(cfg)
	store.SetMaterialProcessor(ai.NewMaterialProcessor(embedder))
	store.SetTrainingRunner(ai.NewSimulatedTrainingRunner())
	store.SetLogHook(func(status, content string) { _ = logSvc.Log(status, "", "", content) })

	switchableModel := ai.NewSwitchableModel(newModelService(cfg))

	app := &apis.App{
		Store:  store,
		Logs:   logSvc,
		Config: cfg,
		Model:  switchableModel,
		Models: apis.NewModelManager(cfg.Model, switchableModel),
		RAG:    ai.NewVectorRAG(store, embedder),
	}

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           apis.NewRouter(app),
		ReadHeaderTimeout: 10 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("whisdom server listening on %s", cfg.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
		close(serveErr)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		log.Println("shutting down...")
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logSvc.Log(logs.StatusInfo, "", "", "server shutting down")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	return nil
}

// seedAdministrator creates the initial administrator account. There is no
// persistence yet (Store is in-memory), so this runs on every startup;
// WHISDOM_ADMIN_EMAIL/WHISDOM_ADMIN_PASSWORD pin deterministic credentials
// for scripted/dev environments, otherwise a random password is generated
// and printed once to the console — never to the log file.
func seedAdministrator(store *system.Store, logSvc *logs.Service) error {
	role, ok := store.RoleByName(system.RoleAdministrator)
	if !ok {
		return errors.New("administrator role missing after seeding default roles")
	}

	email := os.Getenv("WHISDOM_ADMIN_EMAIL")
	if email == "" {
		email = "admin@whisdom.local"
	}
	password := os.Getenv("WHISDOM_ADMIN_PASSWORD")
	generated := password == ""
	if generated {
		var err error
		password, err = randomPassword()
		if err != nil {
			return err
		}
	}

	hash, err := system.HashPassword(password)
	if err != nil {
		return err
	}

	if _, err := store.CreateUser("Administrator", email, hash, role.ID); err != nil {
		return err
	}

	logSvc.Log(logs.StatusInfo, "", "", "seeded initial administrator account")

	fmt.Println("----------------------------------------------------------------")
	fmt.Println(" Initial administrator account (in-memory store: recreated on restart)")
	fmt.Printf(" Email:    %s\n", email)
	if generated {
		fmt.Printf(" Password: %s\n", password)
		fmt.Println(" (random; set WHISDOM_ADMIN_PASSWORD to pin it)")
	} else {
		fmt.Println(" Password: <from WHISDOM_ADMIN_PASSWORD>")
	}
	fmt.Println("----------------------------------------------------------------")

	return nil
}

func randomPassword() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate password: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func newModelService(cfg *config.Config) ai.ModelService {
	switch strings.ToLower(strings.TrimSpace(cfg.Model.Provider)) {
	case "", "stub":
		return ai.NewStubModel()
	case "ollama":
		return ai.NewOllamaModel(cfg.Model)
	default:
		log.Printf("warning: unsupported MODEL_PROVIDER %q; using stub model", cfg.Model.Provider)
		return ai.NewStubModel()
	}
}

func newEmbedder(cfg *config.Config) ai.Embedder {
	switch strings.ToLower(strings.TrimSpace(cfg.Model.Provider)) {
	case "", "stub":
		return ai.NewHashEmbedder()
	case "ollama":
		return ai.NewOllamaEmbedder(cfg.Model)
	default:
		log.Printf("warning: unsupported MODEL_PROVIDER %q; using hash embedder", cfg.Model.Provider)
		return ai.NewHashEmbedder()
	}
}
