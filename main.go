package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

//go:generate postcss --use autoprefixer postcss-pxtorem cssnano --no-map -o assets/style.min.css assets/style.css

//go:embed templates/*
var templates embed.FS

//go:embed assets/*
var assets embed.FS

const appVersion = "v0.0.0"

var pwd string

var CLI struct {
	Init   InitCmd   `cmd:"" help:"Initialize a new Klarity project for writing docs."`
	Build  BuildCmd  `cmd:"" help:"Build Klarity docs from a directory."`
	Dev    DevServer `cmd:"" help:"Opens a dev local dev server for developement."`
	Clean  CleanCmd  `cmd:"" help:"Cleans out all output files from a klarity project"`
	Doctor DoctorCmd `cmd:"" help:"Diagnoses potential issues in a klarity project"`
	VersionCmd
}

type VersionCmd struct {
	Version kong.VersionFlag `name:"version" help:"Display version."`
}

type InitCmd struct {
	Path string `arg:"" name:"path" help:"The directory where the Klarity project should be initialized (e.g., '.' or '/path/to/project')." type:"path"`
}

type BuildCmd struct {
	Path string `arg:"" name:"path" help:"The directory containing the Klarity project to build" type:"path"`
}

type DevServer struct {
	Path string `arg:"" name:"path" help:"The directory containing the Klarity project" type:"path"`
}

type CleanCmd struct {
	Path string `arg:"" name:"path" help:"The directory containing the Klarity project"`
}

type DoctorCmd struct {
	Path string `arg:"" name:"paht" help:"The directory containing the Klarity project"`
}

func (c *DoctorCmd) Run(ctx *kong.Context) error {
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return err
	}
	cfg := ReadConfig(path)

	if cfg.Base_URL == "/" || cfg.Base_URL == "" {
		slog.Warn("the base_url is not configured for distribution")
	}

	if cfg.Entry == "" {
		slog.Error("no entry file configured")
	}

	if len(cfg.Doc_dirs) <= 0 {
		slog.Error("no doc directories configured")
	}

	icons, err := validateFavicons(path)
	if err != nil {
		return err
	}

	if len(icons) == 0 {
		slog.Warn("no favicon detected in the root of the project")
	} else if len(icons) > 1 {
		slog.Warn("multiple favicons detected with different extensions", "favicons", icons)
	}

	return nil
}

func validateFavicons(path string) ([]string, error) {
	var faviconExtList = []string{".ico", ".png", ".svg", ".gif", ".apng", ".jpg"}
	foundFavicons := []string{}

	for _, ext := range faviconExtList {
		faviconPath := filepath.Join(path, fmt.Sprintf("favicon%s", ext))
		if _, err := os.Stat(faviconPath); err == nil {
			foundFavicons = append(foundFavicons, faviconPath)
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking for favicon %s: %w", faviconPath, err)
		}
	}

	return foundFavicons, nil
}

func main() {
	log.SetFlags(log.Llongfile)
	ctx := kong.Parse(&CLI,
		kong.Name("klarity"),
		kong.Description("A very simple markdown docs generator."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Tree:    true,
		}),
		kong.Vars{"version": appVersion},
	)

	err := ctx.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func promptForConfirmation(prompt string) bool {
	var response string
	for {
		fmt.Printf("%s (y/n): ", prompt)
		fmt.Scanln(&response)
		response = strings.ToLower(string(response[0]))
		if response == "y" || response == "Y" {
			return true
		} else if response == "n" || response == "N" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}

func (c *BuildCmd) Run(ctx *kong.Context) error {
	pwd = c.Path
	InitMarkdown(pwd)
	return buildKlarity(c.Path)
}

type PageData struct {
	Title       string
	Content     template.HTML
	Base_URL    string
	FaviconPath string
	FavExt      string
	CustomCSS   string
	SPA         bool
	NavTree     []*NavFolder
	Current     string
}

type NavFolder struct {
	Label string
	Pages []*NavPage
	Open  bool
}

type NavPage struct {
	Title  string
	URL    string
	Active bool
}

var tpl = template.Must(template.ParseFS(templates, "templates/layout.html"))
var partial = template.Must(template.ParseFS(templates, "templates/partial.html"))

var config Config

func buildKlarity(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	c := ReadConfig(path)
	config = c
	docs, err := collectMarkdownFiles(c, path)
	if err != nil {
		return err
	}

	navTree := buildNavTree(path, docs, c.Doc_dirs, c.Entry, c.Title)

	var faviconPath string
	icons, err := validateFavicons(path)
	if err != nil {
		return err
	}

	if len(icons) > 1 {
		// return fmt.Errorf("more than one valid favicon found")
	} else {
		if len(icons) == 0 {
		} else {
			faviconPath = filepath.Join(path, c.Output_dir, filepath.Base(icons[0])) // probably needs better picking
		}
	}

	if c.Visual.CustomCSS != "" {
		custom, err := filepath.Abs(c.Visual.CustomCSS)
		if err != nil {
			return err
		}
		c.Visual.CustomCSS = custom
	}

	html_docs := make(map[string]string)
	for _, doc := range docs {
		b, err := os.ReadFile(doc)
		if err != nil {
			return err
		}
		currentlyRendering = doc
		html, err := renderMarkdown(b)
		if err != nil {
			return err
		}
		currentlyRendering = ""
		html_docs[doc] = html
	}

	c.Output_dir, err = cleanOutputDir(path, c.Output_dir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(c.Output_dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", c.Output_dir, err)
	}

	entry := filepath.Clean(filepath.Join(path, c.Entry))

	for f, page := range html_docs {
		relPath, err := filepath.Rel(path, f)
		if err != nil {
			return fmt.Errorf("unable to determine relative path for '%s': %w", f, err)
		}

		outPath := filepath.Join(c.Output_dir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".html")

		if err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory for '%s': %w", outPath, err)
		}

		var pageTitle string
		var isEntry bool

		if filepath.Clean(f) == entry {
			outPath = filepath.Join(c.Output_dir, "index.html")
			pageTitle = c.Title
			isEntry = true
		} else {
			pageTitle = strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
			isEntry = false
		}

		var relURL string
		if isEntry {
			relURL = "/"
		} else {
			relOut, err := filepath.Rel(c.Output_dir, outPath)
			if err != nil {
				return fmt.Errorf("unable to compute URL for '%s': %w", outPath, err)
			}
			relURL = "/" + filepath.ToSlash(relOut)
		}

		dot_to_blank := func(path string) string {
			if path == "." {
				return ""
			}
			return path
		}

		data := PageData{
			Title:       pageTitle,
			Content:     template.HTML(page),
			Base_URL:    normalizeURL(c.Base_URL),
			FaviconPath: dot_to_blank(filepath.Base(faviconPath)),
			FavExt:      strings.ToLower(filepath.Ext(faviconPath)),
			SPA:         c.Visual.SPA,
			CustomCSS:   dot_to_blank(filepath.Base(c.Visual.CustomCSS)),
			NavTree:     navTree,
			Current:     relURL,
		}

		for _, folder := range data.NavTree {
			folder.Open = false
			for _, pg := range folder.Pages {
				if pg.URL == data.Current {
					pg.Active = true
					folder.Open = true
				}
			}
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("error creating file '%s': %w", outPath, err)
		}

		// if isEntry {
		if err := tpl.Execute(outFile, data); err != nil {
			outFile.Close()
			return fmt.Errorf("error rendering template to '%s': %w", outPath, err)
		}
		// } else {
		// 	if err := partial.Execute(outFile, data); err != nil {
		// 		outFile.Close()
		// 		return fmt.Errorf("error rendering template to '%s': %w", outPath, err)
		// 	}
		// }
		outFile.Close()
	}

	f, err := os.Create(filepath.Join(c.Output_dir, "style.css"))
	if err != nil {
		return err
	}
	defer f.Close()
	built_in_css, err := assets.ReadFile("assets/style.min.css")
	if err != nil {
		return err
	}
	f.Write(built_in_css)

	if c.Ignore_out {
		ignoreTemplate := `# THIS FILE IS AUTOMATICALLY GENERATED, DO NOT MODIFY!

# This file has been automatically generated by Klarity to ingore it's build output
*
`

		ignore, err := os.Create(filepath.Join(c.Output_dir, ".gitignore"))
		if err != nil {
			return fmt.Errorf("could not create .gitignore in %s", c.Output_dir)
		}
		defer ignore.Close()

		ignore.WriteString(ignoreTemplate)
	} else {
		err := os.Remove(filepath.Join(c.Output_dir, ".gitignore"))
		if err != nil {
			return fmt.Errorf("could not remove .gitignore from %s", c.Output_dir)
		}
	}

	writeVarsCSS(c.Visual.Vars, c.Output_dir)

	if faviconPath != "" {
		if err := CopyFile(icons[0], faviconPath); err != nil {
			return err
		}
	}

	if c.Visual.CustomCSS != "" {
		dst := filepath.Join(c.Output_dir, filepath.Base(c.Visual.CustomCSS))
		if err := CopyFile(c.Visual.CustomCSS, dst); err != nil {
			return err
		}
	}

	return nil
}

func (c *CleanCmd) Run(ctx *kong.Context) error {
	path, err := filepath.Abs(c.Path)
	if err != nil {
		return err
	}
	cfg := ReadConfig(path)
	cfg.Output_dir, err = cleanOutputDir(path, cfg.Output_dir)
	if err != nil {
		return err
	}
	fmt.Println("cleaned all build artifacts from", cfg.Output_dir)
	return nil
}

func cleanOutputDir(basePath, outputDir string) (string, error) {
	absOutputDir, err := filepath.Abs(filepath.Join(basePath, outputDir))
	if err != nil {
		return "", err
	}

	if err := os.RemoveAll(absOutputDir); err != nil {
		return "", fmt.Errorf("failed to remove existing output directory '%s': %w", absOutputDir, err)
	}

	return absOutputDir, nil
}

func (c *InitCmd) Run(ctx *kong.Context) error {
	pwd = c.Path
	// InitMarkdown(pwd)
	return initKlarity(c.Path)
}

type initState int

const (
	no initState = iota
	not_empty
	klarity_project
)

func alreadyExists(path string) initState {
	exist, err := os.Stat(path)
	if os.IsNotExist(err) {
		return no
	}

	if !exist.IsDir() {
		return no
	}

	klarityTomlPath := filepath.Join(path, "klarity.toml")
	_, err = os.Stat(klarityTomlPath)
	klarityTomlExists := !os.IsNotExist(err)

	if !klarityTomlExists {
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return not_empty
		}
		if len(dirEntries) > 0 {
			return not_empty
		}
		return no
	}

	config := ReadConfig(path)
	mdFileExists := false

	for _, docDir := range config.Doc_dirs {
		docsPath := filepath.Join(path, docDir)
		if _, err := os.Stat(docsPath); os.IsNotExist(err) {
			continue
		}

		_ = filepath.Walk(docsPath, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error walking path %q: %v\n", p, err)
				return nil
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				mdFileExists = true
				return filepath.SkipDir
			}
			return nil
		})

		if mdFileExists {
			break
		}
	}

	if klarityTomlExists && mdFileExists {
		return klarity_project
	}

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return not_empty
	}

	if len(dirEntries) > 0 {
		return not_empty
	}

	return no
}

func initKlarity(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	switch alreadyExists(path) {
	case klarity_project:
		if !promptForConfirmation(fmt.Sprintf("A Klarity project already exists at '%s'. Do you want to reinitialize it? (not recommended, hard deletes files).", path)) {
			fmt.Println("Initialization cancelled.")
			return nil
		}
		c := ReadConfig(path)
		for _, dir := range c.Doc_dirs {
			os.RemoveAll(filepath.Join(path, dir))
		}
		os.Remove(filepath.Join(path, "klarity.toml"))
	case not_empty:
		if !promptForConfirmation(fmt.Sprintf("The directory '%s' is not empty. Do you want to proceed with initialization? (will mix files).", path)) {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	case no:
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	CreateConfig(path)
	os.Mkdir(filepath.Join(path, "docs"), os.ModePerm)
	f, err := os.Create(filepath.Join(path, "docs", "main.md"))
	if err != nil {
		return err
	}
	_, err = f.WriteString("# Welcome to Klarity!\n\nStart writing your docs here.\n\n```go\nfunc klarity() {\n\tfmt.Println(\"Hello Klarity!\")\n}\n```")
	if err != nil {
		return err
	}

	fmt.Printf("A Klarity project has been successfully created 🚀\n\nGet started by running:\n\nklarity dev %s\n\nor editing %s\n", filepath.Base(path),
		filepath.Join(path, "docs", "main.md"))

	return nil
}
