package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"slices"

	mathjax "github.com/litao91/goldmark-mathjax"
	enclaveCallout "github.com/quailyquaily/goldmark-enclave/callout"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
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
	return slices.Contains(themes, name)
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
			enclaveCallout.New(),
			&anchor.Extender{},
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

// finally fixed the resolver, but still needs [[!]] resolution xd
func (KlarityResolver) ResolveWikilink(n *wikilink.Node) (destination []byte, err error) {
	if currentlyRendering == "" {
		return nil, nil
	}
	base := normalizeURL(ReadConfig(pwd).Base_URL)

	entryAbs := filepath.Clean(filepath.Join(pwd, config.Entry))
	// log.Print("entry: ", entryAbs)

	targetRaw := string(n.Target)
	ext := filepath.Ext(targetRaw)

	// log.Printf("raw: %s, ext: %s", targetRaw, ext)

	var candidateMD string
	baseDir := filepath.Dir(currentlyRendering)
	if ext == "" {
		candidateMD = filepath.Join(baseDir, targetRaw+".md")
	} else if ext == ".md" {
		candidateMD = filepath.Join(baseDir, targetRaw)
	} else {
		return nil, nil // not .md link | handle like default resolver
	}
	candidateMD = filepath.Clean(candidateMD)

	// log.Print("candidate: ", candidateMD)

	relCand, err := filepath.Rel(pwd, candidateMD)
	if err != nil {
		return nil, err
	}
	relCand = filepath.ToSlash(relCand)

	// log.Print("relCandidate: ", relCand)

	var dest string
	if candidateMD == entryAbs {
		dest = base + "/index.html"
	} else {
		dest = base + "/" + strings.TrimSuffix(relCand, ".md") + ".html"
	}

	if len(n.Fragment) > 0 {
		dest += "#" + string(n.Fragment)
	}

	// log.Printf("resolved %s\n", dest)
	return []byte(dest), nil
}

func normalizeURL(url string) string {
	if url == "/" {
		return ""
	}
	return strings.TrimRight(url, "/")
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
