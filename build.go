package main

import (
	"bytes"
	"log"
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

var md goldmark.Markdown

// list of current chroma themes
var themes = []string{
	"abap",
	"algol",
	"algol_nu",
	"arduino",
	"autumn",
	"average",
	"base16-snazzy",
	"borland",
	"bw",
	"catppuccin-frappe",
	"catppuccin-latte",
	"catppuccin-macchiato",
	"catppuccin-mocha",
	"colorful",
	"doom-one",
	"doom-one2",
	"dracula",
	"emacs",
	"evergarden",
	"friendly",
	"fruity",
	"github-dark",
	"github",
	"gruvbox-light",
	"gruvbox",
	"hr_high_contrast",
	"hrdark",
	"igor",
	"lovelace",
	"manni",
	"modus-operandi",
	"modus-vivendi",
	"monokai",
	"monokailight",
	"murphy",
	"native",
	"nord",
	"nordic",
	"onedark",
	"onesenterprise",
	"paraiso-dark",
	"paraiso-light",
	"pastie",
	"perldoc",
	"pygments",
	"rainbow_dash",
	"rose-pine-dawn",
	"rose-pine-moon",
	"rose-pine",
	"rpgle",
	"rrt",
	"solarized-dark",
	"solarized-dark256",
	"solarized-light",
	"swapoff",
	"tango",
	"tokyonight-day",
	"tokyonight-moon",
	"tokyonight-night",
	"tokyonight-storm",
	"trac",
	"vim",
	"vs",
	"vulcan",
	"witchhazel",
	"xcode-dark",
	"xcode",
}

func isValidTheme(name string) bool {
	for _, t := range themes {
		if t == name {
			return true
		}
	}
	return false
}

func InitMarkdown(path string) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	c := ReadConfig(path)
	theme := "rose-pine-moon" // default
	if c.Visual.Theme != "" && isValidTheme(c.Visual.Theme) {
		theme = c.Visual.Theme
	}

	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Linkify,
			extension.Table,
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
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
}

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
