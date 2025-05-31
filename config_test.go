package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "klarity_test_")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return tempDir
}

func TestCreateConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "create config successfully"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := createTempDir(t)
			defer os.RemoveAll(tempDir)

			CreateConfig(tempDir)

			configPath := filepath.Join(tempDir, "klarity.toml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Errorf("CreateConfig() did not create klarity.toml at %s", configPath)
			}

			b, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("failed to read created config file: %v", err)
			}

			var c Config
			err = toml.Unmarshal(b, &c)
			if err != nil {
				t.Errorf("CreateConfig() created a klarity.toml that could not be unmarshalled: %v", err)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name         string
		setupContent Config // The config to write for the test
	}{
		{
			name:         "read basic config",
			setupContent: Config{Title: "My Project", Output_dir: "/build"},
		},
		{
			name:         "read config with different values",
			setupContent: Config{Title: "Another Site", Output_dir: "/public_html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := createTempDir(t)
			defer os.RemoveAll(tempDir)

			configPath := filepath.Join(tempDir, "klarity.toml")
			b, err := toml.Marshal(tt.setupContent)
			if err != nil {
				t.Fatalf("failed to marshal setup content for test: %v", err)
			}
			err = os.WriteFile(configPath, b, 0644)
			if err != nil {
				t.Fatalf("failed to write klarity.toml for test setup: %v", err)
			}

			got := ReadConfig(tempDir)

			gotBytes, err := toml.Marshal(got)
			if err != nil {
				t.Fatalf("failed to marshal read config for comparison: %v", err)
			}
			setupBytes, err := toml.Marshal(tt.setupContent)
			if err != nil {
				t.Fatalf("failed to marshal setup config for comparison: %v", err)
			}

			if string(gotBytes) != string(setupBytes) {
				t.Errorf("ReadConfig() returned unexpected content.\nGot:\n%s\nWant:\n%s", string(gotBytes), string(setupBytes))
			}
		})
	}
}
