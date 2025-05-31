package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/wikilink"
)

var currentlyRendering string

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Linkify,
		extension.Table,
		highlighting.NewHighlighting(
			highlighting.WithStyle("fruity"),
		),
		mathjax.MathJax,
		&wikilink.Extender{
			Resolver: KlarityResolver{},
		},
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithUnsafe(),
	),
)

type KlarityResolver struct{}

func (KlarityResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	if currentlyRendering == "" {
		return nil, nil
	}

	entryAbs := filepath.Clean(filepath.Join(pwd, config.Entry))
	currentAbs := filepath.Clean(currentlyRendering)
	target := string(n.Target)

	ext := filepath.Ext(target)
	if ext == "" {
		target += ".md"
		ext = ".md"
	}
	if ext == ".md" {
		target = strings.TrimSuffix(target, ".md") + ".html"
	}

	entryBase := strings.TrimSuffix(filepath.Base(config.Entry), ".md")
	targetBase := strings.TrimSuffix(filepath.Base(string(n.Target)), ".md")
	if entryBase != "" && targetBase == entryBase {
		target = "index.html"
	}

	var dest string
	if currentAbs == entryAbs {
		dest = target
	} else {
		dest = target
	}

	if len(n.Fragment) > 0 {
		dest += "#" + string(n.Fragment)
	}

	dest = filepath.ToSlash(dest)
	return []byte(dest), nil
}

func collectMarkdownFiles(config Config, root string) ([]string, error) {
	var files []string
	for _, dir := range config.Doc_dirs {
		full := filepath.Join(root, dir)
		err := filepath.Walk(full, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".md" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func renderMarkdown(src []byte) (string, error) {
	var buf bytes.Buffer
	err := md.Convert(src, &buf)
	return buf.String(), err
}
