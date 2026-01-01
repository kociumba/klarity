package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Title      string       `toml:"title"`
	Output_dir string       `toml:"output_dir"`
	Base_URL   string       `toml:"base_url"`
	Doc_dirs   []string     `toml:"doc_dirs"`
	Entry      string       `toml:"entry"`
	Ignore_out bool         `toml:"ignore_out"`
	Visual     VisualConfig `toml:"visual"`
	Dev        DevConfig    `toml:"dev"`
	Editor     EditorConfig `toml:"editor"`
}

type VisualConfig struct {
	Theme     string     `toml:"theme"`
	SPA       bool       `toml:"use_spa"`
	CustomCSS string     `toml:"custom_css"`
	Vars      VarsConfig `toml:"vars"`
}

type VarsConfig struct {
	BGMain          string `toml:"bg_main"`
	BGPanel         string `toml:"bg_panel"`
	BGHover         string `toml:"bg_hover"`
	BGActive        string `toml:"bg_active"`
	BorderSoft      string `toml:"border_soft"`
	BorderHard      string `toml:"border_hard"`
	AccentPrimary   string `toml:"accent_primary"`
	AccentSecondary string `toml:"accent_secondary"`
	TextMain        string `toml:"text_main"`
	TextDim         string `toml:"text_dim"`
}

type DevConfig struct {
	Port int `toml:"port"`
}

type EditorConfig struct {
	Enable bool `toml:"enable_editor"`
}

// path is the directory klarity was called with
func CreateConfig(path string) {
	f, err := os.Create(filepath.Join(path, "klarity.toml"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := toml.Marshal(Config{
		Title:      "Hello klarity!",
		Output_dir: "public",
		Base_URL:   "/",
		Doc_dirs:   []string{"docs"},
		Entry:      "docs/main.md",
		Ignore_out: true,
		Visual: VisualConfig{
			Theme: "rose-pine-moon",
			SPA:   true,
		},
		Dev: DevConfig{
			Port: 5173,
		},
		Editor: EditorConfig{
			Enable: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	f.Write(b)
}

// path is the directory klarity was called with
func ReadConfig(path string) Config {
	b, err := os.ReadFile(filepath.Join(path, "klarity.toml"))
	if err != nil {
		log.Fatal(err, " ", filepath.Join(path, "klarity.toml"))
	}

	var c Config
	err = toml.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}

	if dev_server {
		c.Base_URL = "/"
	}

	return c
}
