package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"sk-ant-api03-1234567890abcdef", "sk-a****cdef"},
		{"short", "********"},
		{"12345678", "********"},
		{"", "********"},
		{"123456789", "1234****6789"},
	}
	for _, tt := range tests {
		if got := maskAPIKey(tt.in); got != tt.want {
			t.Errorf("maskAPIKey(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// setHome points os.UserHomeDir at a temp dir for the duration of the test.
func setHome(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)        // unix
	t.Setenv("USERPROFILE", dir) // windows, harmless elsewhere
}

func TestSaveConfig_Roundtrip(t *testing.T) {
	setHome(t)

	if err := saveConfig("", "sk-ant-test-key", "claude-opus-4-8"); err != nil {
		t.Fatalf("saveConfig: %v", err)
	}
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.ApiKey != "sk-ant-test-key" || cfg.Model != "claude-opus-4-8" {
		t.Fatalf("roundtrip mismatch: %+v", cfg)
	}
	if cfg.Provider != providerAPI {
		t.Errorf("provider = %q, want %q", cfg.Provider, providerAPI)
	}
}

func TestSaveConfig_UpdateKeyKeepsModel(t *testing.T) {
	setHome(t)
	_ = saveConfig("", "first-key", "claude-opus-4-8")

	// Updating only the key must not reset the model.
	if err := saveConfig("", "second-key", ""); err != nil {
		t.Fatalf("saveConfig: %v", err)
	}
	cfg, _ := loadConfig()
	if cfg.Model != "claude-opus-4-8" {
		t.Errorf("model reset to %q, want claude-opus-4-8", cfg.Model)
	}
	if cfg.ApiKey != "second-key" {
		t.Errorf("key = %q, want second-key", cfg.ApiKey)
	}
}

func TestSaveConfig_RejectsInvalidModel(t *testing.T) {
	setHome(t)
	err := saveConfig("", "some-key", "not-a-real-model")
	if err == nil || !strings.Contains(err.Error(), "invalid model") {
		t.Fatalf("expected invalid model error, got %v", err)
	}
}

func TestSaveConfig_RequiresAPIKey(t *testing.T) {
	setHome(t)
	err := saveConfig("", "", "claude-sonnet-4-6")
	if err == nil || !strings.Contains(err.Error(), "API key is required") {
		t.Fatalf("expected API key error, got %v", err)
	}
}

func TestSaveConfig_RejectsInvalidProvider(t *testing.T) {
	setHome(t)
	err := saveConfig("carrier-pigeon", "some-key", "")
	if err == nil || !strings.Contains(err.Error(), "invalid provider") {
		t.Fatalf("expected invalid provider error, got %v", err)
	}
}

// The claude CLI holds its own credentials, so this path must save with no key.
func TestSaveConfig_ClaudeCodeNeedsNoAPIKey(t *testing.T) {
	setHome(t)
	if err := saveConfig(providerClaudeCode, "", ""); err != nil {
		t.Fatalf("saveConfig: %v", err)
	}
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Provider != providerClaudeCode {
		t.Errorf("provider = %q, want %q", cfg.Provider, providerClaudeCode)
	}
	if cfg.Model != defaultModel {
		t.Errorf("model = %q, want %q", cfg.Model, defaultModel)
	}
}

// Configs written before providers existed have no provider field.
func TestLoadConfig_DefaultsProvider(t *testing.T) {
	setHome(t)
	path, err := configPath()
	if err != nil {
		t.Fatalf("configPath: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	legacy := `{"api_key":"sk-ant-old","model":"claude-sonnet-5"}`
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfg.Provider != providerAPI {
		t.Errorf("provider = %q, want %q", cfg.Provider, providerAPI)
	}
}
