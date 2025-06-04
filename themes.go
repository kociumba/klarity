package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeVarsCSS(v VarsConfig, outputDir string) error {
	outPath := filepath.Join(outputDir, "vars.css")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	lines := []string{}

	if v.BGMain != "" {
		lines = append(lines, fmt.Sprintf("  --bg-main:    %s;", v.BGMain))
	}
	if v.BGPanel != "" {
		lines = append(lines, fmt.Sprintf("  --bg-panel:   %s;", v.BGPanel))
	}
	if v.BGHover != "" {
		lines = append(lines, fmt.Sprintf("  --bg-hover:   %s;", v.BGHover))
	}
	if v.BGActive != "" {
		lines = append(lines, fmt.Sprintf("  --bg-active:  %s;", v.BGActive))
	}

	if v.BorderSoft != "" {
		lines = append(lines, fmt.Sprintf("  --border-color-soft: %s;", v.BorderSoft))
	}
	if v.BorderHard != "" {
		lines = append(lines, fmt.Sprintf("  --border-color-hard: %s;", v.BorderHard))
	}

	if v.AccentPrimary != "" {
		lines = append(lines, fmt.Sprintf("  --accent-primary:   %s;", v.AccentPrimary))
	}
	if v.AccentSecondary != "" {
		lines = append(lines, fmt.Sprintf("  --accent-secondary: %s;", v.AccentSecondary))
	}

	if v.TextMain != "" {
		lines = append(lines, fmt.Sprintf("  --text-main: %s;", v.TextMain))
	}
	if v.TextDim != "" {
		lines = append(lines, fmt.Sprintf("  --text-dim:  %s;", v.TextDim))
	}

	if len(lines) == 0 {
		return nil
	}

	fmt.Fprintln(f, ":root {")
	for _, line := range lines {
		fmt.Fprintln(f, line)
	}
	fmt.Fprintln(f, "}")

	return nil
}
