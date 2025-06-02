# CLI Reference

Klarity provides a simple command-line interface for managing your project.

---

## Commands

To get more info on these you can always run the general `klarity -h` or specific context sensitive `klarity <command> -h`

### `klarity init [path]`

Initializes a new Klarity project in the specified directory.  
Creates a `klarity.toml` and a starter `docs/main.md`.

> [!TIP]
> `klarity init` checks if the directory is empty and if it already contains a klarity project

---

### `klarity dev [path]`

Starts a local development server with live reload.  
Default address: [http://localhost:5173](http://localhost:5173)

---

### `klarity build [path]`

Builds your documentation into static HTML files in the output directory. This build is ready for hosting, meaning it resolves paths using the configured `base_url` (default output directory: `public`).

---

### `klarity clean [path]`

Removes all generated output files from the output directory.

---

### `klarity doctor [path]`

Diagnoses potential issues in your Klarity project, such as missing config fields or favicons.

---

### `klarity --version`

Displays the current version of Klarity.