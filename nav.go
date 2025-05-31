package main

import (
	"path/filepath"
	"sort"
	"strings"
)

func buildNavTree(root string, docs []string, docDirs []string, entry string, siteTitle string) []*NavFolder {
	// 1) Normalize absolute paths of each doc_dir
	absDocDirs := make([]string, 0, len(docDirs))
	for _, dd := range docDirs {
		abs := filepath.Clean(filepath.Join(root, dd))
		absDocDirs = append(absDocDirs, abs)
	}

	// 2) Determine the absolute path of the entry file (if any)
	entryAbs := ""
	if entry != "" {
		entryAbs = filepath.Clean(filepath.Join(root, entry))
	}

	// 3) folderKey → []*NavPage
	folderMap := make(map[string][]*NavPage)

	for _, absPath := range docs {
		cleanPath := filepath.Clean(absPath)

		// 3.a) If this is the “entry” .md (to be rendered as /index.html):
		if entryAbs != "" && cleanPath == entryAbs {
			// We place it into folderMap[""] (top‐level, no label),
			// give it Title = siteTitle (from config.Title), URL = "/"
			folderMap[""] = append(folderMap[""], &NavPage{
				Title: siteTitle,
				URL:   "/", // entry always becomes "/"
			})
			continue // skip the rest
		}

		// 3.b) Find which doc_dir this file belongs to
		var matchedDocDir string
		var relInDocDir string
		for _, dd := range absDocDirs {
			// We want exact prefix + separator, so that "/home/.../docs" does not match "/home/.../docs-extra"
			if strings.HasPrefix(cleanPath, dd+string(filepath.Separator)) {
				matchedDocDir = dd
				rel, err := filepath.Rel(dd, cleanPath)
				if err != nil {
					relInDocDir = ""
				} else {
					relInDocDir = rel
				}
				break
			}
		}
		if matchedDocDir == "" {
			// If it’s not in any doc_dir, skip it
			continue
		}

		// 3.c) Compute the URL of this page.
		//     We take its path relative to `root`, swap ".md" → ".html", and prefix "/" for URL.
		relToRoot, err := filepath.Rel(root, cleanPath)
		if err != nil {
			continue
		}
		url := "/" + strings.TrimSuffix(filepath.ToSlash(relToRoot), ".md") + ".html"

		// 3.d) Compute the “folder key” based on relInDocDir:
		//     - If relInDocDir has no slash (e.g. "link.md"), folderKey = ""   (top‐level)
		//     - Otherwise, folderKey = first segment (e.g. "api" for "api/foo.md").
		parts := strings.Split(filepath.ToSlash(relInDocDir), "/")
		var folderKey string
		var pageTitle string
		if len(parts) == 1 {
			// directly under doc_dir
			folderKey = ""
			pageTitle = strings.TrimSuffix(parts[0], ".md")
		} else {
			// in a subfolder
			folderKey = parts[0]
			pageTitle = strings.TrimSuffix(parts[len(parts)-1], ".md")
		}

		// 3.e) Append NavPage to the folderMap
		folderMap[folderKey] = append(folderMap[folderKey], &NavPage{
			Title: pageTitle,
			URL:   url,
		})
	}

	// 4) Build the []*NavFolder in order:
	//    - If folderMap[""] exists, that becomes the first NavFolder with Label = "".
	//    - Then collect all non‐empty keys, sort them, and build NavFolders in that order.
	var out []*NavFolder

	// 4.a) Top‐level pages (folderKey = "")
	if pages, ok := folderMap[""]; ok {
		// Sort pages by Title
		sort.Slice(pages, func(i, j int) bool {
			return pages[i].Title < pages[j].Title
		})
		out = append(out, &NavFolder{
			Label: "", // signals “no folder header” in template
			Pages: pages,
			Open:  true, // top‐level group is always “open”
		})
	}

	// 4.b) Non‐empty keys, sorted lexicographically
	var keys []string
	for k := range folderMap {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		pages := folderMap[key]
		// Sort pages alphabetically
		sort.Slice(pages, func(i, j int) bool {
			return pages[i].Title < pages[j].Title
		})

		nf := &NavFolder{
			Label: key,
			Pages: pages,
			Open:  false, // we’ll flip this to true if it contains the “Current” page
		}
		out = append(out, nf)
	}

	return out
}
