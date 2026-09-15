---
title: Technology index
updated: 2026-09-15
status: active
---

# Technology choices

One file per technology. Each records the version in use, where it is used, why it was chosen, and what was considered instead. Update the file when a version changes.

| Technology | Version | Used in | File |
|------------|---------|---------|------|
| Go | 1.26 | flai | [go.md](go.md) |
| Cobra | v1.10.2 | flai | [go-libraries.md](go-libraries.md) |
| goccy/go-yaml | v1.18 | flai | [go-libraries.md](go-libraries.md) |
| charmbracelet/huh, lipgloss | latest v0.x | flai | [go-libraries.md](go-libraries.md) |
| golangci-lint | v2.5 | flai | [go.md](go.md) |
| GoReleaser | v2 | flai | [go.md](go.md) |
| Node.js | 24 LTS | flaiover | [node.md](node.md) |
| pnpm | 10 | flaiover | [node.md](node.md) |
| SvelteKit / Svelte | 2 / 5 | flaiover | [sveltekit.md](sveltekit.md) |
| Vite | 7 | flaiover | [sveltekit.md](sveltekit.md) |
| Tailwind CSS | 4 | flaiover | [tailwind.md](tailwind.md) |
| TypeScript | 5 | flaiover | [sveltekit.md](sveltekit.md) |
| Apache ECharts | 6 | flaiover | [charts.md](charts.md) |
| markdown-it, shiki, mermaid | latest | flaiover | [markdown.md](markdown.md) |
| MiniSearch | 7 | flaiover | [search.md](search.md) |
| Vitest, Playwright | latest | flaiover | [sveltekit.md](sveltekit.md) |
| Docker | 29 (host), image on node:24-alpine | flaiover, flai | [docker.md](docker.md) |
| Git | 2.47 (host) | flai | [go-libraries.md](go-libraries.md) |
| GitHub Actions, GHCR | n/a | monorepo | [ci.md](ci.md) |
| markdownlint-cli2 | latest | monorepo, template | [ci.md](ci.md) |

Versions marked "latest" are pinned in lockfiles; this table is updated when the pin moves a major version.
