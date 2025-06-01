package main

import (
	"path/filepath"
	"sort"
	"strings"
)

func buildNavTree(root string, docs []string, docDirs []string, entry string, siteTitle string) []*NavFolder {
	absDocDirs := make([]string, 0, len(docDirs))
	for _, dd := range docDirs {
		abs := filepath.Clean(filepath.Join(root, dd))
		absDocDirs = append(absDocDirs, abs)
	}

	entryAbs := ""
	if entry != "" {
		entryAbs = filepath.Clean(filepath.Join(root, entry))
	}

	folderMap := make(map[string][]*NavPage)

	for _, absPath := range docs {
		cleanPath := filepath.Clean(absPath)

		if entryAbs != "" && cleanPath == entryAbs {
			folderMap[""] = append(folderMap[""], &NavPage{
				Title: siteTitle,
				URL:   "/",
			})
			continue
		}

		var matchedDocDir string
		var relInDocDir string
		for _, dd := range absDocDirs {
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
			continue
		}

		relToRoot, err := filepath.Rel(root, cleanPath)
		if err != nil {
			continue
		}
		url := "/" + strings.TrimSuffix(filepath.ToSlash(relToRoot), ".md") + ".html"

		parts := strings.Split(filepath.ToSlash(relInDocDir), "/")
		var folderKey string
		var pageTitle string
		if len(parts) == 1 {
			folderKey = ""
			pageTitle = strings.TrimSuffix(parts[0], ".md")
		} else {
			folderKey = parts[0]
			pageTitle = strings.TrimSuffix(parts[len(parts)-1], ".md")
		}

		folderMap[folderKey] = append(folderMap[folderKey], &NavPage{
			Title: pageTitle,
			URL:   url,
		})
	}

	var out []*NavFolder

	if pages, ok := folderMap[""]; ok {
		var entryPage *NavPage
		var otherPages []*NavPage
		for _, p := range pages {
			if p.URL == "/" {
				entryPage = p
			} else {
				otherPages = append(otherPages, p)
			}
		}

		sort.Slice(otherPages, func(i, j int) bool {
			return otherPages[i].Title < otherPages[j].Title
		})

		var sortedPages []*NavPage
		if entryPage != nil {
			sortedPages = append(sortedPages, entryPage)
		}
		sortedPages = append(sortedPages, otherPages...)

		out = append(out, &NavFolder{
			Label: "",
			Pages: sortedPages,
			Open:  true,
		})
	}

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
			Open:  false,
		}
		out = append(out, nf)
	}

	return out
}
