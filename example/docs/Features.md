# Features Showcase

Explore Klarity's supported features with live examples.

---

## Full Page Search

If you have node installed klarity will generate a search index of the generated site with [PageFind](https://pagefind.app/)

This allows the users to search for any term on the whole wiki and immidietly nvigate to it, search can be opened with `CTRL+K`
or with the search button on the top right of the page. 
 
---

## Wikilinks

Link to other docs using Obsidian-style syntax:

```
See [[main.md]] for the introduction.
```

Result: See [[main.md]] for the introduction.

---

## Autolinks

URLs and emails are automatically linked:

```
Visit https://klarity.example.com or email hello@example.com
```

Result: Visit https://klarity.example.com or email hello@example.com

---

## Syntax Highlighting

Code blocks are highlighted:

The use the ``` syntax and both inline `codeblocks` and blocks are supported:

```go
// hello.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Klarity!")
}
```

---

## GFM (GitHub Flavored Markdown)

```markdown
- [x] Task lists
- [ ] Unchecked item

| Name    | Value |
| ------- | ----- |
| Klarity | 🚀     |

~~Strikethrough~~
```

- [x] Task lists
- [ ] Unchecked item

| Name    | Value |
| ------- | ----- |
| Klarity | 🚀     |

~~Strikethrough~~

---

## LaTeX Math

```markdown
Inline math: $E=mc^2$

Block math:

$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$
```

Inline math: $E=mc^2$

Block math:

$$
\int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
$$

---

## Raw HTML

You can embed HTML directly:

```html
<div style="color: teal;">Custom HTML block!</div>
```

<div style="color: teal;">Custom HTML block!</div>

---

## Callouts

```markdown
> [!NOTE]
> This is a note callout.

> [!WARNING]
> This is a warning callout.

> [!TIP]
> Tips can be shown like this!

> [!CAUTION]
> This callout advises caution!

> [!IMPORTANT]
> This carries imprtant info.
```

> [!NOTE]
> This is a note callout.

> [!WARNING]
> This is a warning callout.

> [!TIP]
> Tips can be shown like this!

> [!CAUTION]
> This callout advises caution!

> [!IMPORTANT]
> This carries imprtant info.

---
