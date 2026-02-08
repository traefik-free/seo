# Contributing to Traefik SEO Plugin

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/traefik-free/seo/issues).
2. If not, open a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Traefik version and plugin version

### Suggesting Features

1. Open an issue with the `enhancement` label.
2. Describe the use case and proposed solution.
3. Discuss before implementing to align on the approach.

### Pull Requests

1. Fork the repository and create a branch from `main`.
2. Make your changes and add tests if applicable.
3. Ensure tests pass: `go test -v ./...`
4. Update documentation if needed.
5. Submit a pull request with a clear description of changes.

### Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Keep functions focused and readable.
- Add comments for exported types and functions.

### Commit Messages

- Use clear, descriptive commit messages.
- Prefer present tense: "Add changefreq to sitemap" instead of "Added changefreq".

## Development Setup

```bash
git clone https://github.com/traefik-free/seo.git
cd seo
go test ./...
```

## Questions?

Feel free to open an issue for any questions about contributing.
