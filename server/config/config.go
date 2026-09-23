// Package config loads server configuration from a JSON file with
// environment variable overrides for deployment-specific values.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds server-wide settings loaded at startup.
type Config struct {
	Addr                 string        `json:"addr"`
	LogDir               string        `json:"log_dir"`
	DataDir              string        `json:"data_dir"`
	TrainingMaterialPath string        `json:"training_material_path"`
	Model                ModelConfig   `json:"model"`
	AccessTokenTTL       time.Duration `json:"-"`
	RefreshTokenTTL      time.Duration `json:"-"`
	JWTSecret            string        `json:"-"`
}

// ModelConfig controls the local language model runtime used for inference.
type ModelConfig struct {
	Provider       string  `json:"provider"`
	Name           string  `json:"name"`
	EmbeddingModel string  `json:"embedding_model"`
	RuntimeURL     string  `json:"runtime_url"`
	Temperature    float64 `json:"temperature"`
}

const (
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour
)

// Load reads the JSON config at path, then applies environment overrides.
// Missing files fall back to built-in defaults so the server can start
// without any configuration present.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Addr:                 ":8080",
		LogDir:               "log",
		DataDir:              "data",
		TrainingMaterialPath: "/training/company-default/",
		Model: ModelConfig{
			Provider:       "stub",
			Name:           "stub",
			EmbeddingModel: "nomic-embed-text",
			RuntimeURL:     "",
			Temperature:    0.2,
		},
		AccessTokenTTL:  defaultAccessTokenTTL,
		RefreshTokenTTL: defaultRefreshTokenTTL,
	}

	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	if v := os.Getenv("PORT"); v != "" {
		cfg.Addr = ":" + v
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("LOG_DIR"); v != "" {
		cfg.LogDir = v
	}
	if v := os.Getenv("DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("MODEL_PROVIDER"); v != "" {
		cfg.Model.Provider = v
	}
	if v := os.Getenv("MODEL_NAME"); v != "" {
		cfg.Model.Name = v
	}
	if v := os.Getenv("MODEL_RUNTIME_URL"); v != "" {
		cfg.Model.RuntimeURL = v
	}
	if v := os.Getenv("MODEL_EMBEDDING_NAME"); v != "" {
		cfg.Model.EmbeddingModel = v
	}

	cfg.JWTSecret = os.Getenv("WHISDOM_JWT_SECRET")

	return cfg, nil
}

// EnsureJWTSecret returns the configured signing secret, generating a random
// ephemeral one when none was supplied. generated reports whether a new
// secret was created, so the caller can warn that restarts will invalidate
// existing sessions.
func (c *Config) EnsureJWTSecret() (generated bool, err error) {
	if c.JWTSecret != "" {
		return false, nil
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return false, fmt.Errorf("generate jwt secret: %w", err)
	}
	c.JWTSecret = hex.EncodeToString(buf)
	return true, nil
}
