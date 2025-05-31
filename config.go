package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Title      string   `toml:"title"`
	Output_dir string   `toml:"output_dir"`
	Doc_dirs   []string `toml:"doc_dirs"`
	Entry      string   `toml:"entry"`
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
