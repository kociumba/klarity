# Config (`klarity.toml`)

Klarity uses a simple `klarity.toml` file for configuration.  
This file is created automatically when you run `klarity init`.

All paths in `klarity.toml` are always relative to it, this means `klarity.toml` always marks the root of the project.

---

## Required Fields

- **title**: The site title (shown on the main page).
- **output_dir**: Directory where the built HTML will be placed.
- **doc_dirs**: List of directories containing your markdown files.
- **entry**: The markdown file that becomes the main entry (`index.html`).

Example:
```toml
title = "My Docs"
output_dir = "output"
doc_dirs = ["docs"]
entry = "docs/main.md"
```

---

## Optional Fields

- **base_url**: The base URL your site will be hosted at.  
  - Set this to the full URL if hosting on GitHub Pages or a subdirectory.
  - Default: `/` *(only works on the local dev server or root hosting)*
    > [!IMPORTANT]
    > This is also required if you want to host on something like giuthub pages
- **ignore_out**: If `true`, Klarity will create a `.gitignore` in the output directory to ignore built files.  
  - Set to `false` if you want to commit the output.
- **[visual]**
  - **theme**: Code highlighting theme.  
    - Default: `"rose-pine-moon"`.  
    - See [theme gallery](https://xyproto.github.io/splash/docs/all.html) for options.
    > [!NOTE]
    > Background colors of those themes are not used for the sake consistency
  - **use_spa**: turn on or off single page navigation, it is highly recommended to keep this `true` since most of the testing it done with it, and [swup](https://swup.js.org/) which enables this behaviour isn't a big dependency.
- **[dev] port**: Port for the dev server.  
  - Default: `5173`.  
  - Must be between 1024-49151.

---

## Example `klarity.toml`

```toml
title = "Example Klarity site"
output_dir = "output"
base_url = "https://username.github.io/repo/"
doc_dirs = ["docs", "other"]
entry = "docs/entry.md"
ignore_out = true

[visual]
theme = "rose-pine-moon"
use_spa = true

[dev]
port = 42069
```

---

## Notes

- If you change `output_dir`, update your deployment scripts accordingly.
- Always set `base_url` to match your hosting path for correct link resolution.
- It is recommended to keep `visial.use_spa` at `true` since Klarity is mostly tested with it enabled.
