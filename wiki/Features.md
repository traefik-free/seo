# Features

## Sitemap

- Dynamic generation based on successful requests (HTTP 200)
- Host-aware filtering
- Priorities: 1.0 for root, 0.8 for others
- lastmod in UTC

## Robots.txt

- Simple format
- Sitemap reference
- `User-agent: *` allowed

## Google Tag Manager

When `gtmID` is set, the plugin:
- Adds GTM script before `</head>`
- Adds noscript iframe after `<body>`
- Properly handles gzip-compressed responses

## Hreflang

For pages with paths `/ru/`, `/en/`, `/mobile/ru/`, `/mobile/en/`, the plugin injects:
- `rel="canonical"` — always points to desktop URL
- `rel="alternate" hreflang="ru"` and `hreflang="en"`
- `rel="alternate" hreflang="x-default"` — points to homepage in defaultLang

## Mobile Alternate

On desktop pages, adds:
- `rel="alternate" media="only screen and (max-width: 640px)"` — link to mobile version

## Ignore Patterns

Built-in exclusions:

| Category | Examples |
|----------|----------|
| VCS | .git, .svn, .hg |
| Config | .env, .htaccess, .well-known |
| CMS | wp-admin, phpmyadmin, administrator |
| Dev | node_modules, vendor, .next, dist |
| API | graphql, swagger, health, metrics |
| Extensions | jpg, png, css, js, zip, pdf, etc. |

Additional patterns can be added via `ignore`.
