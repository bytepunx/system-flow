---
title: Markdown rendering
updated: 2026-11-09
status: active
topics: [dashboard]
---

# Markdown rendering in flaiover

| Package | Version | Purpose |
|---------|---------|---------|
| `markdown-it` | 15.0 | Rendering, with `markdown-it-anchor` 10 for heading links and `markdown-it-task-lists` 2.1 for acceptance criteria checkboxes |
| `yaml` | 2.9 | Front matter parsing on the server, same library the search indexer uses |
| `shiki` | 4.4 | Code block highlighting, lazy loaded after the HTML is in the DOM |
| `mermaid` | 12.0 | Diagrams from fenced `mermaid` blocks, rendered client side, lazy loaded |
| `@tailwindcss/typography` | 0.5 | Prose styling |

Rendering happens in the browser from raw markdown served by `/api/docs/file`, so the server stays a file reader and the SPA owns presentation. Relative `.md` links are rewritten to `/docs/<repo path>` explorer routes; anchors and absolute URLs are left alone. `src/lib/markdown.ts` (S-0012).

## Raw HTML is shown as text

markdown-it runs with `html: false`, so raw HTML in any markdown the dashboard renders, a `<script>`, an `<img onerror>`, or a `<div onclick>`, is shown as text and never becomes markup (S-0627, [I-0294](../issues/I-0294-the-dashboard-renders-repository-markdown-with-markdown-it-html-true-and-no.md)). Documents, item bodies, threads, and messages are written by agents, imported repositories, and forge reviews, and [safety.md](../conventions/safety.md) treats such content as data; with `html: true` any of them could run script in the operator's session, which can accept, move, and edit items. Every page and component that shows markdown takes its HTML from `render()`, so the one option covers them all.

- The task-list checkboxes and the heading permalinks are tokens the plugins make, which markdown-it prints whatever `html` says, so they still render.
- markdown-it's own link check refuses `javascript:`, `vbscript:`, `file:`, and non-image `data:` links, which then show as text; the link rewrite runs only on links that pass it.
- Mermaid runs with `securityLevel: 'strict'`.
- No document under `design/` or `docs/` relies on raw HTML outside code spans and fences. A sanitizer such as DOMPurify was weighed and not chosen: it adds a dependency and an allow-list to keep, for HTML no document needs.
