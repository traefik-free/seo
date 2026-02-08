# Limitations

| Limitation | Description |
|------------|-------------|
| Dynamic only | Paths are collected at runtime; no static route scanning |
| Memory | All unique paths are stored in memory. Suitable for small to medium sites |
| Changefreq | No `<changefreq>` in sitemap |
| Lastmod | Always set to generation time, not per-path modification time |

## Recommendations for Large Sites

- Consider periodic flushing of the path map
- Or external storage (would require plugin modifications)
