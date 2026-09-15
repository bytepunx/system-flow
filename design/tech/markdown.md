---
title: Markdown rendering
updated: 2026-09-15
status: active
---

# Markdown rendering in flaiover

| Package | Purpose |
|---------|---------|
| `markdown-it` | Rendering, with `markdown-it-anchor` for heading links and `markdown-it-task-lists` for acceptance criteria checkboxes |
| `yaml` | Front matter parsing, same library the search indexer uses |
| `shiki` | Code block highlighting, lazy loaded |
| `mermaid` | Diagrams from fenced `mermaid` blocks, rendered client side, lazy loaded |
| `@tailwindcss/typography` | Prose styling |

Rendering happens in the browser from raw markdown served by `/api/docs/file`, so the server stays a file reader and the SPA owns presentation. Relative links between documents are rewritten to dashboard routes.
