# Configuration

## Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `sitemapPath` | string | `/sitemap.xml` | Path for sitemap |
| `robotsPath` | string | `/robots.txt` | Path for robots.txt |
| `ignore` | []string | — | Regex patterns to exclude paths |
| `gtmID` | string | — | Google Tag Manager container ID (GTM-XXXXXX) |
| `defaultLang` | string | `en` | Default language for x-default hreflang |
| `supportedLangs` | []string | `["ru", "en"]` | List of languages for hreflang alternates |

## Example (YAML)

```yaml
http:
  middlewares:
    seo-middleware:
      plugin:
        seo:
          sitemapPath: "/sitemap.xml"
          robotsPath: "/robots.txt"
          gtmID: "GTM-0000000"
          defaultLang: "en"
          supportedLangs: ["ru", "en"]
          ignore:
            - "^/admin/.*"
            - ".*\\.log$"
  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: "my-service"
      middlewares:
        - seo-middleware
```

## Ignored Paths

By default, the following are excluded:
- VCS, config files, admin panels (wp-admin, phpmyadmin, etc.)
- API and documentation (graphql, swagger, health)
- Build artifacts (node_modules, dist, target)
- Media and static assets (images, fonts, archives)

See [Features#ignore-patterns](Features#ignore-patterns) for details.
