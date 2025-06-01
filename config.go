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
	Doc_dirs   []string     `toml:"doc_dirs"`
	Entry      string       `toml:"entry"`
	Ignore_out bool         `toml:"ignore_out"`
	Visual     VisualConfig `toml:"visual"`
	Dev        DevConfig    `toml:"dev"`
}

type VisualConfig struct {
	Theme string `toml:"theme"`
}

type DevConfig struct {
	Port int `toml:"port"`
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
		Doc_dirs:   []string{"docs"},
		Entry:      "docs/main.md",
		Ignore_out: true,
		Visual: VisualConfig{
			Theme: "rose-pine-moon",
		},
		Dev: DevConfig{
			Port: 5173,
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
		log.Fatal(err)
	}

	var c Config
	err = toml.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
