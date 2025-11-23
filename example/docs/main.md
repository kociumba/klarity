# Klarity Documentation

Welcome to the documentation for **Klarity** - a simple, light-config markdown docs generator.

These docs are written using Klarity itself, to see the source for them look in the [github repo](https://github.com/kociumba/klarity/tree/main/example).

## 🚀 Quick Start

1. **Install Klarity:**
   ```shell
    go install github.com/kociumba/klarity@latest
   ```

> [!NOTE]
> If you want to use full page search you will also need the [pagefind binary](https://github.com/pagefind/pagefind/releases) or [node.js](https://nodejs.org/en) installed, if pagefind or node is not found full page search is simply disabled on the generated site.

2. **Initialize a new project:**
   ```shell
   klarity init my-docs
   ```

3. **Start the dev server:**
   ```shell
   klarity dev my-docs
   ```
   Open [http://localhost:5173](http://localhost:5173) in your browser.

4. **Build for production:**
    - make sure `base_url` is set correctly in `klarity.toml`
    - run the build command:
   ```shell
   klarity build my-docs
   ```

## 📁 Project Structure

A default Klarity project contains:
- `klarity.toml` - project configuration
- `docs/main.md` - your documentation entry point

## ✨ Features

For a list of features with examples look in [[Features.md]]

- **Wikilinks:** Use `[[other_file.md]]` to link between docs.
- **Autolinks:** URLs and emails are auto-linked.
- **Syntax Highlighting:** Code blocks are highlighted.
- **GFM:** Tables, task lists, and more.
- **LaTeX:** Inline and block math with MathJax.
- **Raw HTML:** Embed HTML directly.
- **Callouts:** Use GitHub-style alerts.

## 📝 Example

```
# My Project

Welcome to my docs! See [[getting_started.md]] for more info.

- [ ] Task list
- [x] Completed task

Inline math: $E=mc^2$
```

---

**Continue exploring the sidebar for more details on configuration and advanced usage.**
