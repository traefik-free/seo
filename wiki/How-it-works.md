# How It Works

## Request Handling

1. The middleware wraps the next handler.
2. For regular paths, successful responses (200) are recorded in the path map.
3. For `/sitemap.xml` and `/robots.txt`, generated content is returned.

## Sitemap

When sitemap is requested:
- Path map access is locked (mutex)
- Filtering by current host is applied
- URLs are sorted
- XML with loc, lastmod, priority is generated

## robots.txt

Simple text is generated with User-agent and sitemap reference. Sitemap URL is built from X-Forwarded-Proto and Host.

## SEO Links Injection

For HTML pages with locale paths:
- Canonical and alternate links are added before `</head>`
- Mobile alternate is added for desktop pages

## GTM Injection

When gtmID is set:
- Decompress gzip if needed
- Insert script before `</head>`
- Insert noscript after `<body>`
- Re-compress with gzip if needed
