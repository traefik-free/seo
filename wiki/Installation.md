# Installation

The plugin is written in Go and integrates as a Traefik middleware.

## Enabling the Plugin

Add the plugin to Traefik's **static configuration** (traefik.yml, traefik.toml, or environment variables):

```yaml
experimental:
  plugins:
    seo:
      moduleName: "github.com/traefik-free/seo"
      version: "v0.1.2"   # specify the latest release version
```

## Versions

Use a release tag from [Releases](https://github.com/traefik-free/seo/releases). Traefik will download and compile the plugin on startup.

## Next Step

After enabling the plugin, add the middleware to your routers — see [Configuration](Configuration).
