# Traefik SEO Plugin
![Traefik SEO Plugin](https://github.com/traefik-free/seo/raw/main/.assets/icon.png)
This plugin is a middleware for [Traefik](https://traefik.io/), designed to dynamically generate sitemap.xml and robots.txt files based on the paths served by your application. It automatically collects successful (HTTP 200) routes that are not ignored and includes them in the sitemap. This helps improve SEO by providing search engines with an up-to-date site map and robots instructions.

## Features
- Dynamic Sitemap Generation: Collects URLs from successful requests (HTTP 200) and generates an XML sitemap on-the-fly when requested.
- Robots.txt Support: Generates a simple robots.txt file that references the sitemap URL and allows all user agents.
- Configurable Paths: Customize the paths for sitemap (/sitemap.xml by default) and robots (/robots.txt by default).
- Ignore Patterns: Exclude specific paths using regular expressions. Includes comprehensive default ignores for VCS, config files, CMS/admin paths, API endpoints, build artifacts, and common non-SEO file extensions.
- Host-Aware Filtering: Only includes URLs matching the current host in the sitemap.
- Priorities and Lastmod: Assigns priorities (1.0 for root, 0.8 for others) and uses the current UTC time for lastmod.
- Thread-Safe: Uses mutex locking for concurrent access to the path map.
- Gzip Handling: Properly handles gzipped responses when injecting scripts and SEO links, ensuring content integrity.
- Google Tag Manager Integration: Optionally injects GTM script and noscript tags into HTML responses for analytics tracking.
- **Multilingual SEO (hreflang)**: Automatically injects canonical and alternate link tags for locale pages (`/ru/`, `/en/`, `/mobile/ru/`, `/mobile/en/`). Supports configurable default language and supported languages for x-default and hreflang alternates.
- **Mobile Alternate**: Injects `rel="alternate" media="only screen and (max-width: 640px)"` for desktop pages, pointing to the corresponding mobile version.

## Installation
This plugin is written in Go and can be integrated as a Traefik middleware plugin. To use it:

1. Build the Plugin: Clone the repository and build the plugin binary if needed, or use it directly in your Traefik configuration.
  
2. Traefik Configuration: Enable experimental plugins in your Traefik static configuration (e.g., traefik.toml or YAML):
   ```yaml
   experimental:
     plugins:
       seo:
         moduleName: "github.com/traefik-free/seo" # Replace with the actual Go module path
         version: "v0.1.2" # Specify the version
   ```
  
   For more details on Traefik plugins, refer to the [Traefik Plugin Documentation](https://doc.traefik.io/traefik/plugins/overview/).

## Configuration
The plugin accepts a JSON configuration with the following options:
- sitemapPath (string, optional): Path where the sitemap is served. Default: /sitemap.xml.
- robotsPath (string, optional): Path where robots.txt is served. Default: /robots.txt.
- ignore (array of strings, optional): List of regex patterns to ignore when collecting paths for the sitemap. These are compiled as Go regular expressions.
- gtmID (string, optional): Google Tag Manager container ID (e.g., "GTM-XXXXXX"). If provided, the plugin will automatically inject the GTM script into the <head> and noscript iframe into the <body> of HTML responses. This enables easy analytics tracking without modifying your application code.
- defaultLang (string, optional): Default language for x-default hreflang (e.g., "en"). Used for pages with locale paths (/ru/, /en/, /mobile/ru/, /mobile/en/). Default: "en".
- supportedLangs (array of strings, optional): List of supported language codes for hreflang alternates (e.g., ["ru", "en"]). Default: ["ru", "en"].

> **Note:** SEO links (canonical, hreflang, mobile alternate) are injected only for HTML pages with locale paths: `/ru/`, `/en/`, `/mobile/ru/`, or `/mobile/en/`. Other path formats are not modified.

### Default Ignore Patterns
The plugin includes comprehensive built-in ignore patterns to exclude non-SEO-relevant files and paths:
- **Config & backup**: .env, .bak, .old, .example, .sample, .tmpl, .tpl, .dist, config, wp-, sitemap, undefined.
- **VCS & sensitive**: .git, .svn, .hg, .bzr, .htaccess, .htpasswd, .well-known, .aws, .ssh, .mysql_history, .pgpass, .netrc.
- **CMS & admin**: wp-admin, wp-login, wp-content, wp-includes, administrator, phpmyadmin, phpinfo, server-status, actuator, jmx, jolokia, heapdump.
- **Dev & build**: node_modules, vendor, .idea, .vscode, .next, .nuxt, target, dist, build, __pycache__, .pytest_cache.
- **API & docs**: api/v[0-9]+/, graphql, swagger, openapi, api-docs, health, metrics, status.
- **Config files**: Dockerfile, docker-compose, package.json, composer.json, Gemfile, pom.xml, web.config, etc.
- **File extensions**: images (jpg, png, gif, webp, svg, etc.), documents (pdf, doc, docx), media (mp3, mp4, avi), fonts (ttf, woff), code (js, css, json, xml), archives (zip, tar, gz), executables, and more.

You can add custom ignores via the configuration.

### Example Traefik Configuration (YAML)
Attach the middleware to a router in your dynamic configuration:
```yaml
http:
  middlewares:
    seo-middleware:
      plugin:
        seo:
          sitemapPath: "/sitemap.xml"
          robotsPath: "/robots.txt"
          gtmID: "GTM-0000000"  # Your Google Tag Manager ID
          defaultLang: "en"     # Default language for x-default hreflang
          supportedLangs: ["ru", "en"]  # Languages for hreflang alternates
          ignore:
            - "^/admin/.*" # Custom ignore for admin paths
            - ".*\\.log$" # Ignore log files
  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: "my-service"
      middlewares:
        - seo-middleware
```

## Usage
1. Attach to Routers: Add the middleware to your Traefik routers. As requests are handled successfully (200 OK), non-ignored paths are collected.
  
2. Access Sitemap: Visit https://yourdomain.com/sitemap.xml (or your configured path). The sitemap will include all collected URLs, sorted alphabetically, with priorities and lastmod timestamps.
3. Access Robots.txt: Visit https://yourdomain.com/robots.txt. It will contain:
   ```
   User-agent: *
   Sitemap: https://yourdomain.com/sitemap.xml
   ```

4. Path Collection: Only paths that return HTTP 200 and do not match ignore patterns are added. The root path (/) is always included if not present.
5. GTM Integration: If gtmID is set, HTML responses (text/html, status 200) will have the GTM script added before </head> and noscript after <body>. Gzipped responses are decompressed, modified, and re-gzipped automatically.
6. **Multilingual SEO links**: For pages with locale paths (`/ru/`, `/en/`, `/mobile/ru/`, `/mobile/en/`), the plugin automatically injects:
   - `rel="canonical"` — always points to the desktop URL (mobile pages canonicalize to `/ru/...` or `/en/...`).
   - `rel="alternate" hreflang="ru"` and `hreflang="en"` — language alternates for each supported language.
   - `rel="alternate" hreflang="x-default"` — points to the default language homepage (desktop only).
   - `rel="alternate" media="only screen and (max-width: 640px)"` — mobile version URL (desktop pages only).

## How It Works
- Request Handling: The middleware wraps the next handler. For non-special paths, it records successful requests.
- Sitemap Build: When /sitemap.xml is requested, it locks the path map, filters by host, sorts, and generates XML.
- Robots Build: When /robots.txt is requested, it generates a simple text file with the sitemap reference.
- Scheme Detection: Uses X-Forwarded-Proto or falls back to the request scheme for full URLs.
- **SEO Links Injection**: For locale pages, injects canonical and hreflang alternate links before `</head>`, plus mobile alternate for desktop pages.
- GTM Injection: If gtmID is configured, injects GTM scripts into HTML responses, handling gzipped content by decompressing, modifying, and re-compressing if necessary.

## Limitations
- Dynamic Only: Paths are collected at runtime; no static scanning of routes.
- Memory Usage: Stores all unique paths in memory. Suitable for small to medium sites; for large sites, consider periodic flushing or external storage.
- No Change Frequency: Currently, no <changefreq> in sitemap (can be added if needed).
- Lastmod: Always set to the current time on generation, not per-path modification time.

## Contributing
Contributions are welcome! Feel free to open issues or pull requests for improvements, such as adding more features or optimizing performance.

## License
This plugin is licensed under the MIT License. See [LICENSE](LICENSE) for details.

Plugin link: https://plugins.traefik.io/plugins/6950e4534cda2b265225fa58/seo-generator