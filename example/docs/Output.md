# Klarity output

Klarity outputs static html files in its configured `output_dir` path.  
All output is fully static and portable.

## Dependencies

These files do not use any framework like React or Svelte.  
Currently, only two dependencies are in use:

- [MathJax](https://www.mathjax.org/): for rendering LaTeX mathematical notation.
- [swup](https://swup.js.org/): for SPA (single page application) navigation.

SPA navigation can be disabled with `visual.use_spa = false` in the [[Config.md|config]] file.

> [!WARNING]
> Turning off SPA navigation is not recommended since Klarity is intended to be used with it and tested accordingly.

## Entry

In the `output_dir` you will always find at least:

- `index.html`: the entry point of the generated site
- `style.css`: the default built-in CSS styling *(take a look [[Theming.md|here]] for info on themeing)*

If you provide a `favicon` (in any browser supported format), it will be copied to the output.

If `ignore_out = true`, Klarity generates a `.gitignore` file in the output directory to ignore all output files.  
Set `ignore_out = false` if you want to commit the generated files (for example, when manually deploying).

## Link Handling

All links and asset paths in the output respect your configured `base_url`.  
If links are broken, double-check your `base_url` in `klarity.toml`.
