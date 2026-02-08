# GitHub Wiki

These files should be pushed to your project's **wiki** repository.

## How to add the wiki

1. Enable Wiki in repository settings: **Settings → Features → Wiki → Allow**.

2. Clone the wiki repository:
   ```bash
   git clone https://github.com/traefik-free/seo.wiki.git
   cd seo.wiki
   ```

3. Copy all `.md` files from this folder to the wiki repository root (except this README.md).

4. Commit and push:
   ```bash
   git add .
   git commit -m "Add wiki pages"
   git push origin main
   ```

The wiki will be available at: `https://github.com/traefik-free/seo/wiki`
