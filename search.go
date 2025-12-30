package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func injectSearchUI(outputDir, baseURL string) error {
	normalized := normalizeURL(baseURL)
	if !strings.HasSuffix(normalized, "/") {
		normalized += "/"
	}

	data := struct {
		BundlePath string
	}{
		BundlePath: normalized + "pagefind/",
	}

	var buf bytes.Buffer
	if err := searchTpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render search template: %w", err)
	}
	searchHTML := buf.String()

	return filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// TODO: also hardcode editor.html exclusion here
		if filepath.Base(path) == "editor.html" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if !bytes.Contains(content, []byte("</body>")) {
			return nil
		}

		newContent := bytes.ReplaceAll(content,
			[]byte("</body>"),
			[]byte(searchHTML+"\n</body>"),
		)

		return os.WriteFile(path, newContent, 0644)
	})
}
