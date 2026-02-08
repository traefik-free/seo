# Usage

## Attaching the Middleware

1. Add the middleware to your Traefik routers (see [Configuration](Configuration)).
2. After successful requests (200 OK), paths are automatically added to the sitemap.

## Accessing Sitemap and Robots.txt

- **Sitemap:** `https://yourdomain.com/sitemap.xml`
- **Robots:** `https://yourdomain.com/robots.txt`

Example robots.txt:

```
User-agent: *
Sitemap: https://yourdomain.com/sitemap.xml
```

## Path Collection

- Only paths returning HTTP 200 are added
- Paths matching ignore patterns are excluded
- Root path `/` is always included

## GTM

If `gtmID` is set, GTM scripts are injected into HTML responses (text/html, 200). Gzip responses are handled correctly.

## Multilingual Pages

For paths `/ru/`, `/en/`, `/mobile/ru/`, `/mobile/en/`, canonical and hreflang links are added automatically.
