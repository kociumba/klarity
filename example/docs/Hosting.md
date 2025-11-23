# Hosting

This page covers how to host a Klarity page, with a GitHub Pages example.

---

## 📦 Building Your Site

After you set up your Klarity project (more info [[main.md#-quick-start]]), make sure your `klarity.toml` is configured.

**Important:** Set the `base_url` in `klarity.toml` to the path where your site will be hosted.  
For GitHub Pages, this is usually:

```
base_url = "https://<username>.github.io/<repo>/"
```

If you use a custom domain or root, adjust accordingly.

Build your site:

```shell
klarity build [path]
```

This will generate static HTML in the output directory (default: `public`).

---

## 🚀 Deploying to GitHub Pages

1. **Push your code to GitHub.**

2. **Set up GitHub Pages:**
   - Go to your repository's settings.
   - Under "Pages", set the source to the `/public` if manually building and set `ignore_out = false` in `klarity.toml` directory in your branch or use GitHub Actions as below.
> [!CAUTION]
> Hosting by manually building and commiting changes is not recommended, and only really sustainable if the site never really changes 

3. **(Recommended) Use GitHub Actions for automatic deployment:**

Klarity includes a sample workflow at [`.github/workflows/pages.yaml`](https://github.com/kociumba/klarity/blob/main/.github/workflows/pages.yaml)

This is a bare bones example deploy script you can use to deploy your Klarity site to github pages, using github actions (this workflow assumes your repo is the root of your Klarity project):
```yaml
name: Build and deploy Klarity page

on:
  push:
    branches: ["main"]

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  build-and-deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}

    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4


      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "1.x"
        
      - name: Install klarity
        run: | 
          go install github.com/kociumba/klarity@latest

      - name: Build example site
        run: |
            klarity build .

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
            path: 'public'

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

This workflow will build and deploy your docs automatically on every push to `main`.

This workflow is also fully compatible with generating full page search for your build, since node is always pre installed on github actions runners.

---

## 🌐 Other Hosting Options

Since Klarity outputs static HTML, you can host your docs anywhere that serves static files, such as:
- [Fly.io](https://fly.io/)
- [Vercel](https://vercel.com/)
- [Netlify](https://www.netlify.com/)
- [Cloudflare Pages](https://www.cloudflare.com/)
- Your own server

Just upload the contents of your output directory.