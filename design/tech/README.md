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
| goccy/go-yaml | v1.19 | flai | [go-libraries.md](go-libraries.md) |
| charmbracelet/huh | v1.0 | flai | [go-libraries.md](go-libraries.md) |
| golangci-lint | v2.5 | flai | [go.md](go.md) |
| GoReleaser | v2 (pinned in scripts/install-tools.sh) | flai | [go.md](go.md) |
| Node.js | 24 LTS (`.nvmrc`; dev host runs 25) | flaiover | [node.md](node.md) |
| pnpm | 10.34.5 (`packageManager`) | flaiover | [node.md](node.md) |
| SvelteKit / Svelte | 2.63 / 5.56 | flaiover | [sveltekit.md](sveltekit.md) |
| Vite | 8 | flaiover | [sveltekit.md](sveltekit.md) |
| Tailwind CSS | 4.3 with typography | flaiover | [tailwind.md](tailwind.md) |
| TypeScript | 6 | flaiover | [sveltekit.md](sveltekit.md) |
| Apache ECharts | 6 | flaiover | [charts.md](charts.md) |
| yaml, chokidar | 2.9, 5.0 | flaiover | [sveltekit.md](sveltekit.md) |
| markdown-it, shiki, mermaid | 15.0, 4.4, 12.0 | flaiover | [markdown.md](markdown.md) |
| MiniSearch | 7.2 | flaiover | [search.md](search.md) |
| Vitest, Playwright | 4.1, 1.60 | flaiover | [sveltekit.md](sveltekit.md) |
| Docker | 29 (host), image on node:24-alpine | flaiover, flai | [docker.md](docker.md) |
| Git | 2.47 (host) | flai | [go-libraries.md](go-libraries.md) |
| GitHub Actions, GHCR | n/a | monorepo | [ci.md](ci.md) |
| markdownlint-cli2 | 0.20.0 locally via npx (`scripts/lint-md.sh`), action v24 in CI | monorepo, template | [ci.md](ci.md) |
| log/slog, pino, OpenTelemetry, @prometheus-io/client | slog; pino 10.3, otel sdk-node 0.222, client 0.16 | flai, flaiover | [observability.md](observability.md) |

Versions marked "latest" are pinned in lockfiles; this table is updated when the pin moves a major version.
