# Traefik SEO Plugin

![Traefik SEO Plugin](https://github.com/traefik-free/seo/raw/main/.assets/icon.png)

A [Traefik](https://traefik.io/) middleware that dynamically generates **sitemap.xml** and **robots.txt** based on the paths served by your application. The plugin automatically collects successful (HTTP 200) routes that are not ignored and adds them to the sitemap. This improves SEO by providing search engines with an up-to-date site map and robots instructions.

## Quick Start

1. Register the plugin in Traefik's [static configuration](Installation).
2. Add the middleware to your routers in the [dynamic configuration](Configuration).
3. Open `/sitemap.xml` and `/robots.txt` on your domain.

## Key Features

| Feature | Description |
|---------|-------------|
| [Sitemap](Features#sitemap) | Dynamic XML sitemap generation from served requests |
| [Robots.txt](Features#robotstxt) | Automatic robots.txt with sitemap reference |
| [GTM](Features#google-tag-manager) | Google Tag Manager injection into HTML |
| [Hreflang](Features#hreflang) | Canonical and alternate links for multilingual pages |
| [Mobile](Features#mobile-alternate) | Alternate link for mobile version |

## Links

- [Traefik Plugins](https://plugins.traefik.io/plugins/6950e4534cda2b265225fa58/seo-generator)
- [Traefik Plugin Documentation](https://doc.traefik.io/traefik/plugins/overview/)
