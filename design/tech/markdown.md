---
title: Markdown rendering
updated: 2026-09-15
status: active
---

# Markdown rendering in flaiover

| Package | Version | Purpose |
|---------|---------|---------|
| `markdown-it` | 15.0 | Rendering, with `markdown-it-anchor` 10 for heading links and `markdown-it-task-lists` 2.1 for acceptance criteria checkboxes |
| `yaml` | 2.9 | Front matter parsing on the server, same library the search indexer uses |
| `shiki` | 4.4 | Code block highlighting, lazy loaded after the HTML is in the DOM |
| `mermaid` | 12.0 | Diagrams from fenced `mermaid` blocks, rendered client side, lazy loaded |
| `@tailwindcss/typography` | 0.5 | Prose styling |

Rendering happens in the browser from raw markdown served by `/api/docs/file`, so the server stays a file reader and the SPA owns presentation. Relative `.md` links are rewritten to `/docs/<repo path>` explorer routes; anchors and absolute URLs are left alone. `src/lib/markdown.ts` (S-012).
