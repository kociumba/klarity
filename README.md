# Klarity

<img align="right" src="https://raw.githubusercontent.com/kociumba/klarity/512e21d916601b453ba78a17c6c05103f9f16932/example/favicon.svg" alt="klarity icon" width="150" height="150"/>

[![Klarity Wiki](https://img.shields.io/badge/Wiki-Built_with_Klarity-c94e51)](https://kociumba.github.io/klarity/)

Klarity is a simple markdown docs generator 📖

Motivation for clarity came from using things like quartz 4 and gitbook, that while they have nice workflows, are unintuative and complex.

Klarity tries to solve this problem by not requiring any special setup, or cloning anything, the whole generator is a single go binary that works out of the box.

## Usage

To use Klarity install it with (pre built releases will be provided at a later date):

```shell
go install github.com/kociumba/klarity@latest
```
Then initialize Klarity project using:

```shell
klarity init [path]
```
This creates the specified directory or uses an existing one, the created project contains only 2 files:

- `klarity.toml`
- `docs/main.md`

this is all you need to get started with Klarity, `klarity.toml` contains very simple configurations, this is how it looks by default:

```toml
title = "Hello klarity!" # title used on the main page
output_dir = "public" # where to output the built html
base_url = "/" # the url your site will be hosted on, something like https://<username>.github.io/<repo>/ if using github pages
doc_dirs = ["docs"] # all directories that are used during build
entry = "docs/main.md" # the file that will become index.html(the entry of the site)
ignore_out = true # creates a .gitignore that automatically ignores the output_dir

[visual]
theme = "rose-pine-moon" # the theme used for code, see more in the section below
use_spa = true # uses swup to provide spa navigation on the generated site

[dev]
port = 5173 # accepts ports between 1024-49151
```

Now that you have your project created, run:

```shell
klarity dev [path]
```

This will open a very simple dev server hosting your docs. This dev server opens on http://localhost:5173 and live reloads while you make changes to in your markdown files. 

> [!NOTE]
> Keep in mind this dev server is pretty bare bones and only partially supports file additions and deletions.

When you want to build the docs for hosting, you have to do 2 things

- set `base_url` in `klarity.toml` to the base url your site will be hosted on. If you are hosting klarity on github pages your base url will be `https://<username>.github.io/<repo name>/`

- run the build command to build a production version of your site:
    ```shell
    klarity build [path]
    ```

This builds klarity into the output dir from `klarity.toml`, this is ready for hosting using something like github pages, there is a simple entry, `index.html` in the root that should be picked up by any hosting provider like that.

Klarity will by default place a `.gitignore` in the output directory that ignores everything in it, this can also be turned off in the config, if you want to commit the generated files.

## Documentation

The entire Klarity wiki is written _with Klarity_. You can browse it live here:  
👉 [https://kociumba.github.io/klarity/](https://kociumba.github.io/klarity/)

The source code for it is availible in the [example](./example/) directory.


## Features

Klarity supports a wide range of markdown features:

- wikilinks - obsidian style `[[other_file.md]]` wiki links are supported (only `!` image resolution is missing for now)
- autolinks - things like emails or full links automatically become `<a>` in html
- syntax highliting - by default code is highlited using the `rose-pine-moon`, this can be changed to any of these themes: [theme gallery](https://xyproto.github.io/splash/docs/all.html)
- GFM(github flavoured markdown) - most gfm features are supproted like taks lists, tables or strikethrough
- Latex - math notation with inline and block latex is supported and rendered through [mathjax](https://www.mathjax.org/)
- raw html - inserting raw html into markdown is also fully supported
- github callouts - more info about them here: [github gfm docs](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax#alerts)

The generated pages are fully static and use spa like navigation, for a smooth experience.

This allows Klarity docs to be easely served in minutes on platforms like github pages.
